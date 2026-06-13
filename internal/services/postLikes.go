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

func (s *PostLikeService) ToggleLike(postID, userID uint) error {
	var existing models.PostLike
	result := s.db.Where("post_id = ? AND user_id = ?", postID, userID).First(&existing)

	if result.Error == gorm.ErrRecordNotFound {
		newLike := models.PostLike{PostID: postID, UserID: userID, Liked: true}
		return s.db.Create(&newLike).Error

	}

	if result.Error != nil {
		return result.Error
	}

	existing.Liked = !existing.Liked
	return s.db.Save(&existing).Error
}

func (s *PostLikeService) GetLikesForPost(postID, callerID uint) (int64, bool, error) {
	var count int64
	userLiked := false
	countResult := s.db.Model(&models.PostLike{}).Where("post_id = ? AND liked = true", postID).Count(&count)
	if countResult.Error != nil {
		return 0, false, countResult.Error
	}

	if callerID != 0 {
		var row models.PostLike
		err := s.db.Where("post_id = ? AND user_id = ? AND liked = true", postID, callerID).First(&row).Error
		if err == nil {
			userLiked = true
		} else if err == gorm.ErrRecordNotFound {
			userLiked = false
		} else if err != nil {
			return 0, false, err
		}
	}

	return count, userLiked, nil
}
