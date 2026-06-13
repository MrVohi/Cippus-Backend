package handlers

import (
	"cippus-backend/internal/models"
	"cippus-backend/internal/services"
	"cippus-backend/internal/ws"
	"time"
)

type RegisterRequest struct {
	Email        string `json:"email"        binding:"required,email,max=255"`
	Username     string `json:"username"     binding:"required,min=2,max=30"`
	Password     string `json:"password"     binding:"required,min=8,max=128"`
	CaptchaToken string `json:"captcha"      binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=1,max=128"`
}

type PasswordResetRequestRequest struct {
	Email string `json:"email" binding:"required,email,max=255"`
}

type PasswordResetConfirmRequest struct {
	Token       string `json:"token"       binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=8,max=128"`
}

type AuthHandler struct {
	service *services.AuthService
}

type PatchMeRequest struct {
	Username          *string `json:"username" binding:"omitempty,min=2,max=30"`
	Bio               *string `json:"bio"      binding:"omitempty,max=500"`
	NotifReply        *bool   `json:"notif_replies"`
	NotifFollows      *bool   `json:"notif_follows"`
	NotifLike         *bool   `json:"notif_likes"`
	NotifMention      *bool   `json:"notif_mentions"`
	NotifDM           *bool   `json:"notif_dms"`
	NotifPostApproved *bool   `json:"notif_post_approved"`
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

type CommentResponse struct {
	ID        uint            `json:"id"`
	PostID    uint            `json:"post_id"`
	UserID    uint            `json:"user_id"`
	Author    *AuthorResponse `json:"author,omitempty"`
	Content   string          `json:"content"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

func toCommentResponse(c models.Comment) CommentResponse {
	resp := CommentResponse{
		ID:        c.ID,
		PostID:    c.PostID,
		UserID:    c.UserID,
		Content:   c.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
	if c.Author.ID != 0 {
		resp.Author = &AuthorResponse{
			UserID:    c.Author.ID,
			Username:  c.Author.Username,
			AvatarURL: c.Author.AvatarURL,
		}
	}
	return resp
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

type MessageHandler struct {
	service *services.MessageService
	Hub     *ws.Hub
	notif   *services.NotificationService
}

type PushHandler struct {
	service     *services.PushService
	vapidPublic string
}

type NotificationHandler struct {
	service *services.NotificationService
}

type SearchHandler struct {
	service *services.SearchService
}

type SearchPostResponse struct {
	PostResponse
	Distance float64 `json:"distance"`
}

func toSearchPostResponse(p services.PostWithDistance) SearchPostResponse {
	return SearchPostResponse{
		PostResponse: toPostResponse(p.Post),
		Distance:     p.Distance,
	}
}