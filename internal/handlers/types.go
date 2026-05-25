package handlers

import "cippus-backend/internal/services"

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

type ProjectHandler struct {
	service *services.ProjectService
}
