package models

type PostCategory struct {
	PostID     uint `gorm:"primaryKey"`
	Post       Post
	CategoryID uint `gorm:"primaryKey"`
	Category   Category
}
