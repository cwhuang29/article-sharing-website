package models

import (
	"time"
)

type Password struct {
	UserID    string
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Token     string    `gorm:"unique;not null;size:64"`
	MaxAge    int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoUpdateTime"`
}
