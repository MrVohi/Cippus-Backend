package handlers

import "cippus-backend/internal/services"

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
	Username *string `json:"username" binding:"omitempty,min=2,max=30"`
	Bio      *string `json:"bio"      binding:"omitempty,max=500"`
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
