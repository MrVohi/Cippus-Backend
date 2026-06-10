package services

import (
    "cippus-backend/internal/models"
    "gorm.io/gorm"
)

type CommentService struct {
    db *gorm.DB
}

func NewCommentService(db *gorm.DB) *CommentService {
    return &CommentService{db: db}
}

func (s *CommentService) GetCommentsByPostID(postID uint) ([]models.Comment, error) {
    var comments []models.Comment
    result := s.db.Preload("Author").Where("post_id = ?", postID).Find(&comments)
    return comments, result.Error
}

func (s *CommentService) CreateComment(postID uint, userID uint, content string) (models.Comment, error) {
    comment := models.Comment{
        PostID:  postID,
        UserID:  userID,
        Content: content,
    }
    result := s.db.Create(&comment)
    return comment, result.Error
}

func (s *CommentService) DeleteComment(commentID uint, userID uint) error {
    result := s.db.Where("id = ? AND user_id = ?", commentID, userID).Delete(&models.Comment{})
    return result.Error
}