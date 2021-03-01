package models

import (
	"time"
)

type Admin struct {
	Email     string    `gorm:"primaryKey;unique;not null;size:100" json:"email" yaml:"email"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"-"`
}
