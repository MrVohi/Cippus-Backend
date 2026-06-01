package handlers

import (
	"cippus-backend/internal/models"
	"cippus-backend/internal/services"
	"time"
)

type RegisterRequest struct {
	Email        string `json:"email"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	CaptchaToken string `json:"captcha"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type PasswordResetRequestRequest struct {
	Email string `json:"email"`
}

type PasswordResetConfirmRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type AuthHandler struct {
	service *services.AuthService
}

type PatchMeRequest struct {
	Username *string `json:"username"`
	Bio      *string `json:"bio"`
}

type PostHandler struct {
	service *services.PostService
}

type MinioHandler struct {
	service      *services.MinioService
	postServices *services.PostService
}

type CategoryHandler struct {
	service *services.CategoryService
}

type CategoryResponse struct {
	CategoryID  uint   `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UserResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Bio       string `json:"bio"`
	Role      string `json:"role"`
	AvatarURL string `json:"avatarUrl"`
}

func toUserResponse(u models.User) UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Bio:       u.Bio,
		Role:      string(u.Role),
		AvatarURL: u.AvatarURL,
	}
}

type AuthorResponse struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type PostResponse struct {
	PostID           uint               `json:"post_id"`
	AuthorID         uint               `json:"author_id"`
	Author           *AuthorResponse    `json:"author,omitempty"`
	Title            string             `json:"title"`
	Content          string             `json:"content"`
	ImageURL         string             `json:"image_url"`
	Stuck            bool               `json:"stuck"`
	ModerationStatus string             `json:"moderation_status"`
	Categories       []CategoryResponse `json:"categories"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

func toPostResponse(p models.Post) PostResponse {
	cats := make([]CategoryResponse, len(p.Categories))
	for i, c := range p.Categories {
		cats[i] = CategoryResponse{CategoryID: c.ID, Name: c.Name, Description: c.Description}
	}
	resp := PostResponse{
		PostID:           p.ID,
		AuthorID:         p.UserId,
		Title:            p.Title,
		Content:          p.Content,
		ImageURL:         p.ImageURL,
		Stuck:            p.Stuck,
		ModerationStatus: string(p.ModerationStatus),
		Categories:       cats,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
	}
	if p.Author.ID != 0 {
		resp.Author = &AuthorResponse{
			UserID:    p.Author.ID,
			Username:  p.Author.Username,
			AvatarURL: p.Author.AvatarURL,
		}
	}
	return resp
}

type ProjectHandler struct {
	service *services.ProjectService
}
