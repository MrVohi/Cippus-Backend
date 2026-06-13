package services

import (
	"cippus-backend/internal/models"
	"fmt"

	"gorm.io/gorm"
)

func NewCommentService(db *gorm.DB) *CommentService {
	return &CommentService{db: db}
}

func (s *CommentService) GetCommentsByPostID(postID uint) ([]models.Comment, error) {
	var comments []models.Comment
	
	result := s.db.Preload("Author").Where("post_id = ?", postID).Find(&comments)
	if result.Error != nil {
		return nil, result.Error
	}

	return comments, nil
}

func (s *CommentService) CreateComment(postID, userID uint, content string) (models.Comment, error) {
	comment := models.Comment{PostID: postID, UserID: userID, Content: content}
	result := s.db.Create(&comment)
	if result.Error != nil {
		return models.Comment{}, result.Error
	}

	return comment, nil
}

func (s *CommentService) UpdateComment(commentID, callerID uint, content string) (models.Comment, error) {
	var comment models.Comment
	result := s.db.Where("id = ? AND user_id = ?", commentID, callerID).First(&comment)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound  {
			return models.Comment{}, fmt.Errorf("not found or forbidden")
		}
		return models.Comment{}, result.Error
	}
	comment.Content = content
	result = s.db.Save(&comment)
	if result.Error != nil {
		return models.Comment{}, result.Error
	}
	return comment, nil
}

func (s *CommentService) DeleteComment(commentID, callerID uint, callerRole string) error {
	if callerRole == "moderator" || callerRole == "admin" {
		result := s.db.Delete(&models.Comment{}, commentID)
		if result.RowsAffected == 0 {
			return fmt.Errorf("not found or forbidden")
		}
	} else {
		result := s.db.Where("id = ? AND user_id = ?", commentID, callerID).Delete(&models.Comment{})
		if result.RowsAffected == 0 {
			return fmt.Errorf("not found or forbidden")
		}
	}
	return nil
}