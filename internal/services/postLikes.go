package services

import (
	"cippus-backend/internal/models"

	"gorm.io/gorm"
)

type PostLikeService struct {
	db *gorm.DB
}

func NewPostLikeService(db *gorm.DB) *PostLikeService {
	return &PostLikeService{db: db}
}

func (s *PostLikeService) LikePost(postID uint, userID uint) error {
	like := models.PostLike{}
	result := s.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&like)

	if result.Error == gorm.ErrRecordNotFound {
		like = models.PostLike{
			PostID: postID,
			UserID: userID,
			Liked:  true,
		}
		return s.db.Create(&like).Error
	}

	like.Liked = !like.Liked
	return s.db.Save(&like).Error
}

func (s *PostLikeService) GetLikeCount(postID uint) (int64, error) {
	var count int64
	result := s.db.Model(&models.PostLike{}).Where("post_id = ? AND liked = ?", postID, true).Count(&count)
	return count, result.Error
}

func (s *PostLikeService) HasUserLiked(postID uint, userID uint) (bool, error) {
	like := models.PostLike{}
	result := s.db.Where("post_id = ? AND user_id = ? AND liked = ?", postID, userID, true).First(&like)
	if result.Error == gorm.ErrRecordNotFound {
		return false, nil
	}
	return like.Liked, result.Error
}
