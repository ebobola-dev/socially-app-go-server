package minio_service

import (
	"fmt"
)

type Bucket struct {
	Name    string
	IsImage bool
}

var (
	AvatarsBucket = &Bucket{
		Name:    "avatars",
		IsImage: true,
	}
	PostsBucket = &Bucket{
		Name:    "posts",
		IsImage: true,
	}
	MessagesBucket = &Bucket{
		Name:    "messages",
		IsImage: true,
	}
	Apks = &Bucket{
		Name:    "apks",
		IsImage: false,
	}
)

var BucketList = []*Bucket{AvatarsBucket, PostsBucket, MessagesBucket, Apks}

func IsValidBucket(bucket string) bool {
	return bucket == AvatarsBucket.Name ||
		bucket == PostsBucket.Name ||
		bucket == MessagesBucket.Name ||
		bucket == Apks.Name
}

func BucketFromString(strBucket string) (*Bucket, error) {
	switch strBucket {
	case AvatarsBucket.Name:
		return AvatarsBucket, nil
	case PostsBucket.Name:
		return PostsBucket, nil
	case MessagesBucket.Name:
		return MessagesBucket, nil
	case Apks.Name:
		return Apks, nil
	default:
		return nil, fmt.Errorf("invalid bucket name: %s", strBucket)
	}
}
