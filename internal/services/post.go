package services

import (
	"cippus-backend/internal/models"
	"fmt"

	"gorm.io/gorm"
)

func (s *PostService) GetPosts(filters GetPostsFilter) ([]models.Post, error) {
	row := []models.Post{}
	query := s.db.Preload("Author").Preload("Categories")
	if filters.CategoryID != nil {
		query = query.Joins("JOIN post_categories ON post_categories.post_id = posts.id AND post_categories.category_id = ?", *filters.CategoryID)
	}

	if filters.AuthorID != nil {
		query = query.Where("user_id = ?", *filters.AuthorID)
	}

	if filters.Stuck != nil {
		query = query.Where("stuck = ?", *filters.Stuck)
	}

	if filters.Q != nil && *filters.Q != "" {
		query = query.Where("title ILIKE ? OR content ILIKE ?", "%"+*filters.Q+"%", "%"+*filters.Q+"%")
	}

	if filters.Sort == "new" {
		query = query.Order("created_at DESC")
	} else {
		query = query.Order("updated_at DESC")
	}

	result := query.Find(&row)
	if result.Error != nil {
		return nil, result.Error
	}

	return row, nil
}

func (s *PostService) GetPostsById(id uint) (models.Post, error) {
	row := models.Post{}
	result := s.db.Preload("Author").Preload("Categories").Where("id = ?", id).First(&row)
	if result.Error == gorm.ErrRecordNotFound {
		return models.Post{}, fmt.Errorf("post not found")
	} else if result.Error != nil {
		return models.Post{}, fmt.Errorf("error with db")
	}

	return row, nil
}

func (s *PostService) CreatePost(authorID uint, input PostInput) (models.Post, error) {
	post := models.Post{
		UserId:           authorID,
		Title:            input.Title,
		Content:          input.Content,
		ImageURL:         input.ImageURL,
		ModerationStatus: models.ModerationPending,
	}

	if input.Stuck != nil {
		post.Stuck = *input.Stuck
	}

	result := s.db.Create(&post)
	if result.Error != nil {
		return models.Post{}, fmt.Errorf("cannot register new post")
	}

	if input.CategoryIDs != nil {
		for _, category := range input.CategoryIDs {
			result := s.db.Create(&models.PostCategory{PostID: post.ID, CategoryID: category})
			if result.Error != nil {
				return post, result.Error
			}
		}

	}

	return post, nil
}

func (s *PostService) UpdatePost(id uint, authorID uint, userRole string, input PostInput) (models.Post, error) {
	post := models.Post{}
	if userRole == "" {
		return models.Post{}, fmt.Errorf("Cannot find user's role")
	}

	result := s.db.Model(&models.Post{}).Where("id = ?", id).First(&post)
	if result.Error == gorm.ErrRecordNotFound {
		return models.Post{}, fmt.Errorf("post not found")
	} else if result.Error != nil {
		return models.Post{}, fmt.Errorf("error with db")
	}

	if post.UserId != authorID && userRole == "user" {
		return models.Post{}, fmt.Errorf("forbidden")
	}

	if input.Title != "" {
		post.Title = input.Title
	}
	if input.Content != "" {
		post.Content = input.Content
	}
	if input.ImageURL != "" {
		post.ImageURL = input.ImageURL
	}
	if input.Stuck != nil {
		post.Stuck = *input.Stuck
	}

	post.ModerationStatus = models.ModerationPending

	result = s.db.Save(&post)
	if result.Error != nil {
		return models.Post{}, fmt.Errorf("cannot update post")
	}

	result = s.db.Where("post_id = ?", post.ID).Delete(&models.PostCategory{})
	if result.Error != nil {
		return models.Post{}, fmt.Errorf("cannot update post categories")
	}
	if input.CategoryIDs != nil {
		for _, category := range input.CategoryIDs {
			result := s.db.Save(&models.PostCategory{PostID: post.ID, CategoryID: category})
			if result.Error != nil {
				return post, result.Error
			}
		}

	}

	return post, nil
}

func (s *PostService) DeletePost(id uint, authorID uint, userRole string) error {
	post := models.Post{}
	if userRole == "" {
		return fmt.Errorf("Cannot find user's role")
	}

	result := s.db.Model(&models.Post{}).Where("id = ?", id).First(&post)
	if result.Error == gorm.ErrRecordNotFound {
		return fmt.Errorf("post not found")
	} else if result.Error != nil {
		return fmt.Errorf("error with db")
	}

	if post.UserId != authorID && userRole == "user" {
		return fmt.Errorf("forbidden")
	}

	result = s.db.Delete(&models.Post{}, id)
	return result.Error
}

func (s *PostService) UpdatePostImage(id uint, imageURL string) error {
	res := s.db.Model(&models.Post{}).Where("id = ?", id).Update("image_url", imageURL)
	return res.Error
}
func NewPostService(db *gorm.DB) *PostService {
	return &PostService{db: db}
}
