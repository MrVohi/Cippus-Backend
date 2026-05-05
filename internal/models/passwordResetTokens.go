package models

import (
	"time"

	"gorm.io/gorm"
)

type PasswordResetToken struct {
	gorm.Model
	UserID    uint
	User      User
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
}
