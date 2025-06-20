package handler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strings"
	"time"

	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	auth_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/auth"
	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	user_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/user"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	"github.com/ebobola-dev/socially-app-go-server/internal/model"
	"github.com/ebobola-dev/socially-app-go-server/internal/repository"
	minio_service "github.com/ebobola-dev/socially-app-go-server/internal/service/minio"
	image_util "github.com/ebobola-dev/socially-app-go-server/internal/util/image"
	"github.com/ebobola-dev/socially-app-go-server/internal/util/nullable"
	pagination "github.com/ebobola-dev/socially-app-go-server/internal/util/pagination"
	"github.com/ebobola-dev/socially-app-go-server/internal/util/short_flag"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
)

type userHandler struct{}

func NewUserHandler() IUserHandler {
	return &userHandler{}
}

func (h *userHandler) CheckUsername(c *fiber.Ctx) error {
	scope := middleware.GetAppScope(c)
	payload := struct {
		Username string `validate:"required,username_length,username_charset,username_start_digit,username_start_dot"`
	}{
		Username: c.Query("username"),
	}
	if err := scope.Validate.Struct(payload); err != nil {
		return err
	}

	tx := middleware.GetTX(c)
	exists, ex_err := scope.UserRepository.ExistsByUsername(tx, payload.Username)
	if ex_err != nil {
		return ex_err
	}
	return c.JSON(fiber.Map{
		"username": payload.Username,
		"exists":   exists,
	})
}

func (h *userHandler) GetById(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	payload := struct {
		UserId string `validate:"required,uuid4"`
	}{
		UserId: c.Params("user_id"),
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}
	userId := uuid.MustParse(payload.UserId)
	tx := middleware.GetTX(c)
	short := short_flag.FromFiberCtx(c)
	user, get_err := s.UserRepository.GetByID(tx, userId, repository.GetUserOptions{
		IncludeDeleted:     true,
		CountSubscriptions: !short,
	})
	if get_err != nil && !errors.Is(get_err, gorm.ErrRecordNotFound) {
		return get_err
	} else if errors.Is(get_err, gorm.ErrRecordNotFound) {
		return common_error.NewRecordNotFoundErr("User")
	}

	return c.JSON(user.ToDto(model.SerializeUserOptions{
		Safe:  user.ID == middleware.GetUserId(c),
		Short: short,
	}))

}

func (h *userHandler) DeleteMyAccount(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	userId := middleware.GetUserId(c)
	tx := middleware.GetTX(c)

	user, _ := s.UserRepository.GetByID(tx, userId, repository.GetUserOptions{})

	//% Soft delete user
	if err := s.UserRepository.SoftDelete(tx, userId); errors.Is(err, gorm.ErrRecordNotFound) {
		return common_error.NewRecordNotFoundErr("User")
	} else if err != nil {
		return err
	}

	//% Delete refresh tokens
	if _, err := s.RefreshTokenRepository.DeleteByUserId(tx, userId); errors.Is(err, gorm.ErrRecordNotFound) {
		return auth_error.ErrInvalidToken
	} else if err != nil {
		return err
	}

	//% Delete avatar is exists
	if user.AvatarID != nil {
		if err := s.MinioService.DeleteAvatar(c.Context(), user.AvatarID.String()); err != nil {
			s.Log.Exception(err)
			return api_error.NewUnexceptedErr(err)
		}
	}

	tx.Commit()
	return auth_error.ErrAccountDeleted
}

func (h *userHandler) Search(c *fiber.Ctx) error {
	tx := middleware.GetTX(c)
	s := middleware.GetAppScope(c)
	userId := middleware.GetUserId(c)
	pag := middleware.GetPagination(c)

	pattern := c.Query("pattern")
	users, err := s.UserRepository.Search(tx, repository.SearchUsersOptions{
		Pattern:    pattern,
		Pagination: pag,
		IgnoreId:   userId,
	})
	if err != nil {
		return err
	}
	return c.JSON(struct {
		Pattern    string                `json:"pattern"`
		Count      int                   `json:"count"`
		Pagination pagination.Pagination `json:"pagination"`
		Users      []model.ShortUserDto  `json:"users"`
	}{
		Pattern:    pattern,
		Count:      len(users),
		Pagination: pag,
		Users: lo.Map(users, func(user model.User, _ int) model.ShortUserDto {
			return user.ToShortDto()
		}),
	})
}

func (h *userHandler) UpdateProfile(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	payload := struct {
		Fullname    *string                         `json:"fullname" validate:"omitempty,max=32"`
		Username    *string                         `json:"username" validate:"omitempty,username_length,username_charset,username_start_digit,username_start_dot"`
		Gender      nullable.Nullable[model.Gender] `json:"gender" validate:"omitempty,gender"`
		DateOfBirth *string                         `json:"date_of_birth" validate:"omitempty,datebt"`
		AboutMe     *string                         `json:"about_me" validate:"omitempty,max=256"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return common_error.ErrInvalidJSON
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}
	tx := middleware.GetTX(c)
	userId := middleware.GetUserId(c)
	user, _ := s.UserRepository.GetByID(tx, userId, repository.GetUserOptions{})

	hasUpdates := false
	if payload.Fullname != nil && !nullable.StringEqual(payload.Fullname, user.Fullname) {
		user.Fullname = payload.Fullname
		hasUpdates = true
	}
	if payload.Username != nil && *payload.Username != user.Username {
		user.Username = *payload.Username
		hasUpdates = true
	}
	if payload.Gender.Present && !nullable.StringEqual((*string)(payload.Gender.Value), (*string)(user.Gender)) {
		user.Gender = payload.Gender.Value
		hasUpdates = true
	}
	if payload.DateOfBirth != nil {
		dob, _ := time.Parse("02.01.2006", *payload.DateOfBirth)
		if dob != user.DateOfBirth {
			user.DateOfBirth = dob
			hasUpdates = true
		}
	}
	if payload.AboutMe != nil && !nullable.StringEqual(payload.AboutMe, user.AboutMe) {
		user.AboutMe = payload.AboutMe
		hasUpdates = true
	}
	if !hasUpdates {
		return user_error.ErrNothingToUpdateProfile
	}
	if err := s.UserRepository.Update(tx, user); err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"updated_user": user.ToFullDto(true),
	})
}

func (h *userHandler) UpdatePassword(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	payload := struct {
		NewPassword string `json:"new_password" validate:"required,password"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return common_error.NewInvalidJsonErr("Need json body { 'new_password': string } (at least one letter, at least one digit, between 8 and 32 characters)")
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}
	tx := middleware.GetTX(c)
	userId := middleware.GetUserId(c)
	user, _ := s.UserRepository.GetByID(tx, userId, repository.GetUserOptions{})

	hashedPassword, err := s.HashService.HashPassword(payload.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hashedPassword
	if err := s.UserRepository.Update(tx, user); err != nil {
		return err
	}
	return c.SendStatus(200)
}

func (h *userHandler) UpdateAvatar(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	payload := struct {
		AvatarType string `validate:"required,avatar_type"`
	}{
		AvatarType: c.FormValue("avatar_type"),
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}
	avatarType := *model.AvatarTypeFromString(&payload.AvatarType)
	tx := middleware.GetTX(c)
	userId := middleware.GetUserId(c)
	user, _ := s.UserRepository.GetByID(tx, userId, repository.GetUserOptions{})
	user.AvatarType = &avatarType

	//% if new avatar type is internal
	if avatarType != model.ExternalAvatar {
		//% if previous avatar exists - delete it
		if user.AvatarID != nil {
			s.MinioService.DeleteAvatar(c.Context(), user.AvatarID.String())
			user.AvatarID = nil
		}
		if err := s.UserRepository.Update(tx, user); err != nil {
			return err
		}
		return c.JSON(fiber.Map{
			"updated_user": user.ToFullDto(true),
		})
	}

	//% if new avatar type is external
	fileHeader, err := c.FormFile("avatar")
	if err != nil {
		return common_error.NewBadRequestErr("avatar field is required (must be a file)")
	}
	if fileHeader.Filename == "" || fileHeader.Size == 0 {
		return common_error.NewBadRequestErr("avatar field must be a file")
	}
	if fileHeader.Size > s.Cfg.MaxImageSize {
		return user_error.NewAvatarTooLargeErr(fileHeader.Size, s.Cfg.MaxImageSize)
	}
	extension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !slices.Contains(s.Cfg.AllowedImageExtensions, extension) {
		return user_error.NewBadAvatarExtensionErr(extension, s.Cfg.AllowedImageExtensions)
	}
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()
	var buffer bytes.Buffer
	limitedReader := io.LimitReader(file, s.Cfg.MaxImageSize+1)
	size, err := io.Copy(&buffer, limitedReader)
	if err != nil {
		return err
	}
	if size > s.Cfg.MaxImageSize {
		return user_error.NewAvatarTooLargeErr(fileHeader.Size, s.Cfg.MaxImageSize)
	}
	data := buffer.Bytes()
	if err := image_util.ValidateMime(data); err != nil {
		return user_error.NewInvalidImageErr(err.Error())
	}

	if err := image_util.ValidateImageDecode(data); err != nil {
		data, err = image_util.ConvertWithMagick(data)
		if err != nil {
			return user_error.NewInvalidImageErr(err.Error())
		}
		if err := image_util.ValidateImageDecode(data); err != nil {
			return user_error.NewInvalidImageErr(err.Error())
		}
	} else {
		data, err = image_util.ConvertToJPEG(data)
		if err != nil {
			return err
		}
	}
	spilttedImages, err := image_util.SplitImageBytes(data)
	if err != nil {
		return err
	}
	newAvatarId := uuid.New()
	for _, image := range spilttedImages {
		s.MinioService.Save(
			c.Context(),
			minio_service.AvatarsBucket,
			fmt.Sprintf("%s/%s.jpg", newAvatarId, image.Size.String()),
			image.Data, "image/jpeg",
		)
	}
	//% if previous avatar exists - delete it
	if user.AvatarID != nil {
		s.MinioService.DeleteAvatar(c.Context(), user.AvatarID.String())
	}
	user.AvatarID = &newAvatarId
	if err := s.UserRepository.Update(tx, user); err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"updated_user": user.ToFullDto(true),
	})
}

func (h *userHandler) DeleteAvatar(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	tx := middleware.GetTX(c)
	userId := middleware.GetUserId(c)
	user, _ := s.UserRepository.GetByID(tx, userId, repository.GetUserOptions{})
	user.AvatarType = nil
	if user.AvatarID != nil {
		if err := s.MinioService.DeleteAvatar(c.Context(), user.AvatarID.String()); err != nil {
			return err
		}
	}
	user.AvatarID = nil
	if err := s.UserRepository.Update(tx, user); err != nil {
		return err
	}
	return c.JSON(fiber.Map{
		"updated_user": user.ToFullDto(true),
	})
}

func (h *userHandler) GetPrivileges(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	payload := struct {
		UserId string `validate:"required,uuid4"`
	}{
		UserId: c.Query("user_id"),
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}
	targetUid := uuid.MustParse(payload.UserId)
	tx := middleware.GetTX(c)
	pag := middleware.GetPagination(c)
	privileges, err := s.UserRepository.GetPrivileges(tx, repository.GetUserPrivilegesOptions{
		Pagination: pag,
		UserID:     targetUid,
		CountUsers: true,
	})
	if err != nil {
		return err
	}
	return c.JSON(struct {
		UserId      uuid.UUID                `json:"user_id"`
		Count       int                      `json:"count"`
		Pagintation pagination.Pagination    `json:"pagination"`
		Privileges  []model.UserPrivilegeDto `json:"privileges"`
	}{
		UserId:      targetUid,
		Count:       len(privileges),
		Pagintation: pag,
		Privileges: lo.Map(privileges, func(userPrivilege model.UserPrivilege, _ int) model.UserPrivilegeDto {
			return userPrivilege.ToDto()
		}),
	})
}

func (h *userHandler) Follow(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	payload := struct {
		TargetUID string `validate:"required,uuid4"`
	}{
		TargetUID: c.Params("user_id"),
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}
	targetUID := uuid.MustParse(payload.TargetUID)
	tx := middleware.GetTX(c)
	subscruberId := middleware.GetUserId(c)
	newFollowersCount, err := s.UserRepository.Follow(tx, subscruberId, targetUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common_error.NewRecordNotFoundErr("Target user")
		}
		if errors.Is(err, repository.ErrAlreadyFollowing) {
			return common_error.ErrAlreadyFolowingConflict
		}
		if errors.Is(err, repository.ErrSubActionYourself) {
			return common_error.ErrFollowYourselfConflict
		}
		return err
	}
	return c.JSON(fiber.Map{
		"new_target_followers_count": newFollowersCount,
	})
}

func (h *userHandler) Unfollow(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	payload := struct {
		TargetUID string `validate:"required,uuid4"`
	}{
		TargetUID: c.Params("user_id"),
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}
	targetUID := uuid.MustParse(payload.TargetUID)
	tx := middleware.GetTX(c)
	subscruberId := middleware.GetUserId(c)
	newFollowersCount, err := s.UserRepository.Unfollow(tx, subscruberId, targetUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common_error.NewRecordNotFoundErr("Target user")
		}
		if errors.Is(err, repository.ErrNotFollowingAnyway) {
			return common_error.ErrNotFollowingConflict
		}
		if errors.Is(err, repository.ErrSubActionYourself) {
			return common_error.ErrFollowYourselfConflict
		}
		return err
	}
	return c.JSON(fiber.Map{
		"new_target_followers_count": newFollowersCount,
	})
}

func (h *userHandler) GetFollowers(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	payload := struct {
		TargetUID string `validate:"required,uuid4"`
	}{
		TargetUID: c.Params("user_id"),
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}
	targetUID := uuid.MustParse(payload.TargetUID)
	tx := middleware.GetTX(c)
	pag := middleware.GetPagination(c)
	followers, err := s.UserRepository.GetFollowers(tx, repository.GetSubscriptionsOptions{
		TargetUID:  targetUID,
		Pagination: pag,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common_error.NewRecordNotFoundErr("Target user")
		}
		return err
	}
	return c.JSON(struct {
		UserId      uuid.UUID             `json:"user_id"`
		Count       int                   `json:"count"`
		Pagintation pagination.Pagination `json:"pagination"`
		Followers   []model.FollowerDto   `json:"followers"`
	}{
		UserId:      targetUID,
		Count:       len(followers),
		Pagintation: pag,
		Followers: lo.Map(followers, func(us model.UserSubscription, _ int) model.FollowerDto {
			return us.ToFollowerDto()
		}),
	})
}

func (h *userHandler) GetFollowing(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	payload := struct {
		TargetUID string `validate:"required,uuid4"`
	}{
		TargetUID: c.Params("user_id"),
	}
	if err := s.Validate.Struct(payload); err != nil {
		return err
	}
	targetUID := uuid.MustParse(payload.TargetUID)
	tx := middleware.GetTX(c)
	pag := middleware.GetPagination(c)
	following, err := s.UserRepository.GetFollowing(tx, repository.GetSubscriptionsOptions{
		TargetUID:  targetUID,
		Pagination: pag,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return common_error.NewRecordNotFoundErr("Target user")
		}
		return err
	}
	return c.JSON(struct {
		UserId      uuid.UUID             `json:"user_id"`
		Count       int                   `json:"count"`
		Pagintation pagination.Pagination `json:"pagination"`
		Following   []model.FollowingDto  `json:"following"`
	}{
		UserId:      targetUID,
		Count:       len(following),
		Pagintation: pag,
		Following: lo.Map(following, func(us model.UserSubscription, _ int) model.FollowingDto {
			return us.ToFollowingDto()
		}),
	})
}
