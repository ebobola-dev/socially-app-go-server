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

func (us *UserSubscription) ToJson(options SerializeUserSubscriptionOptions) map[string]interface{} {
	result := make(map[string]interface{})
	result["followed_at"] = us.CreatedAt
	if options.IncludeFollower {
		result["follower"] = us.Follower.ToJson(SerializeUserOptions{Short: true})
	}
	if options.IncludeTarget {
		result["target"] = us.Target.ToJson(SerializeUserOptions{Short: true})
	}
	return result
}

type SerializeUserSubscriptionOptions struct {
	IncludeFollower bool
	IncludeTarget   bool
}
