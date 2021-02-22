package models

import (
	"time"
)

type Login struct {
	Email     string    `gorm:"primaryKey;unique"`
	Token     string    `gorm:"unique;not null;size:64"`
	MaxAge    int       `gorm:"not null"`
	LastLogin time.Time `gorm:"not null"`
}
