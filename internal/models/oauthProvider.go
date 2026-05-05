package models

import "gorm.io/gorm"

type OAuthProvider struct {
	gorm.Model
	UserID         uint
	User           User
	Provider       OAuthProviderName
	ProviderUserId string
}
