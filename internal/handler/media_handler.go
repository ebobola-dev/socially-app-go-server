package handler

import (
	"errors"
	"fmt"
	"strings"

	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	minio_service "github.com/ebobola-dev/socially-app-go-server/internal/service/minio"
	image_util "github.com/ebobola-dev/socially-app-go-server/internal/util/image"

	"github.com/gofiber/fiber/v2"
)

type mediaHandler struct{}

func NewMediaHandler() IMediaHandler {
	return &mediaHandler{}
}

func (h *mediaHandler) Get(c *fiber.Ctx) error {
	s := middleware.GetAppScope(c)
	strBucket := c.Params("bucket")
	path := c.Params("*")
	bucket, bErr := minio_service.BucketFromString(strBucket)
	if bErr != nil {
		return common_error.NewBadRequestErr("Invalid bucket name, allowed: avavars, posts, messages, apks")
	}
	obj, stat, err := s.MinioService.Get(c.Context(), bucket, path)
	if errors.Is(err, minio_service.ErrObjectNotFound) {
		var respErr = common_error.NewMinioNotFoundErr(bucket.Name, path)
		//% If not found && bucket is image -> try to return original(size) image
		if bucket.IsImage {
			requestedSize, sizeErr := image_util.SizeFromPath(path)
			if sizeErr != nil {
				return respErr
			}
			newPath := strings.Replace(
				path,
				fmt.Sprintf("%s.jpg", requestedSize),
				fmt.Sprintf("%s.jpg", image_util.SizeOriginal),
				1,
			)
			obj, stat, err = s.MinioService.Get(c.Context(), bucket, newPath)
			if err != nil {
				return respErr
			}
			s.Log.Debug("Size[%s] not found, returing original...", requestedSize.String())
		} else {
			return respErr
		}
	} else if err != nil {
		return err
	}
	c.Set("Content-Type", stat.ContentType)
	c.Response().Header.Set("Content-Length", fmt.Sprintf("%d", stat.Size))
	return c.SendStream(obj, int(stat.Size))
}
