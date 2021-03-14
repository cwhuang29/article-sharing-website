package models

type Tag struct {
	ID       int       `gorm:"primaryKey"`
	Value    string    `gorm:"not null;size:50"`
	Views    uint64    `gorm:"default:0"`
	Articles []Article `gorm:"many2many:articles_tags;"`
}
