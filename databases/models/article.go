package models

import (
	"gorm.io/gorm"
	"time"
)

// For JavaScript, length of a Mandarin word is 1. But for Go, the len() value is based on bytes
// This leads to inconvenience when counting words. Currently, I'll just give a larger word count limit to tolerate this issue
type Article struct {
	ID              int            `gorm:"primaryKey" json:"id"`
	Title           string         `gorm:"not null;size:255" json:"title" yaml:"title"`
	Subtitle        string         `gorm:"size:255" json:"subtitle" yaml:"subtitle"`
	Authors         string         `gorm:"size:50" json:"author" yaml:"author"`
	ReleaseDate     time.Time      `json:"release_date" yaml:"releaseDate"`
	Category        string         `gorm:"size:50;not null" json:"category" yaml:"category"`
	Tags            []Tag          `gorm:"many2many:articles_tags;"`
	Outline         string         `gorm:"size:800" json:"outline" yaml:"outline"`
	CoverPhoto      string         `gorm:"size:300" json:"cover_photo" yaml:"coverPhoto"`
	Content         string         `gorm:"not null" json:"content" yaml:"content"`           // Without size, the type will be longtext
	AdminOnly       bool           `gorm:"default:false" json:"admin_only" yaml:"adminOnly"` // Database stored 0/1
	BookmarkedUsers []User         `gorm:"many2many:users_articles_bookmark;"`
	LikedUsers      []User         `gorm:"many2many:users_articles_like;"`
	CreatedAt       time.Time      `json:"-"`
	UpdatedAt       time.Time      `json:"-"`
	DeletedAt       gorm.DeletedAt `json:"-"`
}
