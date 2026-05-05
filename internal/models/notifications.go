package models

import (
	"time"

	"gorm.io/gorm"
)

type Notification struct {
	gorm.Model
	RecipientID uint
	Recipient   User
	ActorID     uint
	Actor       User
	Type        NotificationType
	EntityType  ContentType
	EntityID    uint
	ReadAt      *time.Time
}
