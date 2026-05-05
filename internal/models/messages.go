package models

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	SenderID   uint
	Sender     User
	ReceiverID uint
	Receiver   User
	Content    string
	ReadAt     *time.Time
}
