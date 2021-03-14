package models

import (
	"time"
)

type Login struct {
	UserID    string
	User      User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Token     string    `gorm:"unique;not null;size:64"`
	MaxAge    int       `gorm:"not null"`
	LastLogin time.Time `gorm:"autoUpdateTime"`
}
