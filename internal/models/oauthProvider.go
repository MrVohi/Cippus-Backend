package models

import "gorm.io/gorm"

type OAuthProvider struct {
	gorm.Model
	UserID         uint
	User           User
	Provider       OAuthProviderName `gorm:"uniqueIndex:idx_oauth_provider;not null"`
	ProviderUserId string            `gorm:"uniqueIndex:idx_oauth_provider;not null"`
}
