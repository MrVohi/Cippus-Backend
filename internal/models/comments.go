package models

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PostID           uint
	Post             Post
	UserID           uint
	Author           User `gorm:"foreignKey:UserID"`
	Content          string
	ModerationStatus ModerationStatus
}
