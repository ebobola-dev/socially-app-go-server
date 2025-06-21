package model

import (
	"time"

	"github.com/google/uuid"
)

type UserSubscription struct {
	FollowerID uuid.UUID `gorm:"type:char(36);primaryKey"`
	TargetID   uuid.UUID `gorm:"type:char(36);primaryKey"`
	CreatedAt  time.Time `gorm:"not null;autoCreateTime(3)"`

	Follower User `gorm:"foreignKey:FollowerID;constraint:OnDelete:CASCADE"`
	Target   User `gorm:"foreignKey:TargetID;constraint:OnDelete:CASCADE"`
}

func (us *UserSubscription) ToFollowerDto() FollowerDto {
	return FollowerDto{
		FollowedAt: us.CreatedAt,
		Follower:   us.Follower.ToShortDto(),
	}
}

func (us *UserSubscription) ToFollowingDto() FollowingDto {
	return FollowingDto{
		FollowedAt: us.CreatedAt,
		Target:     us.Target.ToShortDto(),
	}
}

type FollowerDto struct {
	FollowedAt time.Time    `json:"followed_at"`
	Follower   ShortUserDto `json:"follower"`
}

type FollowingDto struct {
	FollowedAt time.Time    `json:"followed_at"`
	Target     ShortUserDto `json:"target"`
}
