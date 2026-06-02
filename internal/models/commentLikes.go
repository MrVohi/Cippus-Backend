package models

import "gorm.io/gorm"

type CommentLike struct {
	gorm.Model
	UserID    uint `gorm:"uniqueIndex:idx_comment_like;not null"`
	User      User
	CommentID uint `gorm:"uniqueIndex:idx_comment_like;not null"`
	Comment   Comment
	Liked     bool
}
