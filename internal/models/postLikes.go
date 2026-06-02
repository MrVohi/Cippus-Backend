package models

import "gorm.io/gorm"

type PostLike struct {
	gorm.Model
	UserID uint `gorm:"uniqueIndex:idx_post_like;not null"`
	User   User
	PostID uint `gorm:"uniqueIndex:idx_post_like;not null"`
	Post   Post
	Liked  bool
}
