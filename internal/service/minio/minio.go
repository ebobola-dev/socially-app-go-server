package minio_service

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/ebobola-dev/socially-app-go-server/internal/config"
	api_error "github.com/ebobola-dev/socially-app-go-server/internal/errors"
	common_error "github.com/ebobola-dev/socially-app-go-server/internal/errors/common"
	image_util "github.com/ebobola-dev/socially-app-go-server/internal/util/image"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type IMinioService interface {
	Save(ctx context.Context, bucket *Bucket, objectName string, data []byte, contentType string) error
	Delete(ctx context.Context, bucket *Bucket, objectName string) error
	Get(ctx context.Context, bucket *Bucket, objectName string) (*minio.Object, minio.ObjectInfo, error)
	DeleteAvatar(ctx context.Context, avatarID string) error
}

type minioService struct {
	Client *minio.Client
}

func NewMinioService(ctx context.Context, cfg *config.MinioConfig) IMinioService {
	client, err := minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.USER, cfg.PASSWORD, ""),
		Secure: false,
	})
	if err != nil {
		panic(fmt.Errorf("minio connect failed: %w", err))
	}

	svc := &minioService{
		Client: client,
	}

	for _, bucket := range BucketList {
		if exists, err := client.BucketExists(ctx, bucket.Name); err != nil {
			panic(fmt.Errorf("bucketExists(%s): %w", bucket.Name, err))
		} else if !exists {
			err = client.MakeBucket(ctx, bucket.Name, minio.MakeBucketOptions{})
			if err != nil {
				panic(fmt.Errorf("makeBucket(%s): %w", bucket.Name, err))
			}
		}
	}

	return svc
}

func (m *minioService) Save(ctx context.Context, bucket *Bucket, objectName string, data []byte, contentType string) error {
	_, err := m.Client.PutObject(ctx, bucket.Name, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

func (m *minioService) Delete(ctx context.Context, bucket *Bucket, objectName string) error {
	err := m.Client.RemoveObject(ctx, bucket.Name, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		if minio.ToErrorResponse(err).StatusCode == 404 {
			return common_error.NewMinioNotFoundErr(objectName)
		}
		return err
	}
	return nil

}

func (m *minioService) Get(ctx context.Context, bucket *Bucket, objectName string) (*minio.Object, minio.ObjectInfo, error) {
	obj, err := m.Client.GetObject(ctx, bucket.Name, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, minio.ObjectInfo{}, err
	}
	stat, statErr := obj.Stat()
	if statErr != nil {
		if minio.ToErrorResponse(statErr).StatusCode == 404 {
			return nil, minio.ObjectInfo{}, common_error.NewMinioNotFoundErr(objectName)
		}
		return nil, minio.ObjectInfo{}, statErr
	}
	return obj, stat, nil
}

func (m *minioService) DeleteAvatar(ctx context.Context, avatarID string) error {
	for _, size := range image_util.ImageSizesList {
		objectName := fmt.Sprintf("%s/%s.jpg", avatarID, size.String())
		err := m.Delete(ctx, AvatarsBucket, objectName)
		var apiErr api_error.IApiError
		if err != nil {
			if errors.As(err, &apiErr) {
				continue
			}
			return err
		}
	}
	return nil
}
