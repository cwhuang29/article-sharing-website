package models

import (
	"gorm.io/gorm"
	"time"
)

type Article struct {
	ID          int            `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"not null;size:255" json:"title" yaml:"title"`
	Subtitle    string         `gorm:"size:255" json:"subtitle" yaml:"subtitle"`
	Authors     string         `gorm:"size:50" json:"author" yaml:"author"`
	ReleaseDate time.Time      `json:"release_date" yaml:"releaseDate"`
	Category    string         `gorm:"not null" json:"category" yaml:"category"`
	Tags        []Tag          `gorm:"many2many:articles_tags;"`
	Content     string         `gorm:"not null" json:"content" yaml:"content"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   gorm.DeletedAt `json:"-"`
}
