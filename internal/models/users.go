package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username          string `gorm:"uniqueIndex;not null"`
	Email             string `gorm:"uniqueIndex;not null" json:"-"`
	Bio               string
	PasswordHash      string `json:"-"`
	AvatarURL         string
	Role              UserRole
	NotifReply        bool
	NotifFollow       bool
	NotifLike         bool
	NotifMention      bool
	NotifDM           bool
	NotifPostFlagged  bool
	NotifPostApproved bool
}
