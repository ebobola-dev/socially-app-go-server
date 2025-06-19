package handler

import (
	"fmt"
	"strings"

	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	"github.com/ebobola-dev/socially-app-go-server/internal/middleware"
	minio_service "github.com/ebobola-dev/socially-app-go-server/internal/service/minio"
	image_util "github.com/ebobola-dev/socially-app-go-server/internal/util/image"
	"github.com/minio/minio-go/v7"

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
	var obj *minio.Object
	var stat minio.ObjectInfo
	var err error
	if bucket.IsImage {
		requestedSize, sizeErr := image_util.SizeFromPath(path)
		if sizeErr != nil {
			return common_error.NewMinioNotFoundErr(path)
		}
		orderedSizes := image_util.GetOrderedSizeFrom(requestedSize)
		for _, size := range orderedSizes {
			newPath := strings.Replace(
				path,
				fmt.Sprintf("%s.jpg", requestedSize),
				fmt.Sprintf("%s.jpg", size),
				1,
			)
			obj, stat, err = s.MinioService.Get(c.Context(), bucket, newPath)
			if err != nil && size == image_util.SizeOriginal {
				return err
			}
			if err != nil {
				s.Log.Debug("size %s skipped", size)
			}
			if err == nil {
				break
			}
		}
	} else {
		obj, stat, err = s.MinioService.Get(c.Context(), bucket, path)
		if err != nil {
			return err
		}
	}
	c.Set("Content-Type", stat.ContentType)
	c.Response().Header.Set("Content-Length", fmt.Sprintf("%d", stat.Size))
	return c.SendStream(obj, int(stat.Size))
}
