package models

import "gorm.io/gorm"

type PushSubscription struct {
	gorm.Model
	UserID    uint
	User      User
	Endpoint  string
	P256dhKey string
	AuthKey   string
}
