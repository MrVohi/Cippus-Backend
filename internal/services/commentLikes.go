package services

import (
	"cippus-backend/internal/models"

	"gorm.io/gorm"
)

type CommentLikeService struct {
	db *gorm.DB
}

func NewCommentLikeService(db *gorm.DB) *CommentLikeService {
	return &CommentLikeService{db: db}
}

func (s *CommentLikeService) ToggleLike(CommentID, userID uint) error {
	var existing models.CommentLike
	result := s.db.Where("comment_id = ? AND user_id = ?", CommentID, userID).First(&existing)

	if result.Error == gorm.ErrRecordNotFound {
		newLike := models.CommentLike{CommentID: CommentID, UserID: userID, Liked: true}
		return s.db.Create(&newLike).Error

	}

	if result.Error != nil {
		return result.Error
	}

	existing.Liked = !existing.Liked
	return s.db.Save(&existing).Error
}

func (s *CommentLikeService) GetLikesForComment(CommentID, callerID uint) (int64, bool, error) {
	var count int64
	userLiked := false
	countResult := s.db.Model(&models.CommentLike{}).Where("comment_id = ? AND liked = true", CommentID).Count(&count)
	if countResult.Error != nil {
		return 0, false, countResult.Error
	}

	if callerID != 0 {
		var row models.CommentLike
		err := s.db.Where("comment_id = ? AND user_id = ? AND liked = true", CommentID, callerID).First(&row).Error
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
