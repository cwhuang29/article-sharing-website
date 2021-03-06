package models

import (
	"time"
)

type Login struct {
	Email     string    `gorm:"not null"`
	Token     string    `gorm:"unique;not null;size:64"`
	MaxAge    int       `gorm:"not null"`
	LastLogin time.Time `gorm:"not null"`
}
