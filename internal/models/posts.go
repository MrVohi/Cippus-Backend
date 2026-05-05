package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	UserId           uint
	Author           User `gorm:"foreignKey:UserID"`
	Title            string
	Content          string `gorm:"type:text"`
	ImageURL         string
	ModerationStatus ModerationStatus
}
