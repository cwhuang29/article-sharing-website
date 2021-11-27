package models

import (
	"time"
)

// This table is used for the password reset feature
type Password struct {
	UserID    string
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Token     string    `gorm:"unique;not null;size:64"`
	MaxAge    int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoUpdateTime"`
}
