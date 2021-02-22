package models

import (
	"time"
)

type Info struct {
	About         string    `gorm:"not null"`
	Footer        string    `gorm:"not null"`
	FacebookLink  string    `gorm:"not null"`
	InstagramLink string    `gorm:"not null"`
	TwitterLink   string    `gorm:"not null"`
	YoutubeLink   string    `gorm:"not null"`
	LinkedinLink  string    `gorm:"not null"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"-"`
}
