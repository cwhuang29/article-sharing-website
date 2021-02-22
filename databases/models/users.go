package models

import (
	"time"
)

type User struct {
	ID               int       `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"unique;not null;size:20" json:"name" yaml:"name"`
	Password         string    `gorm:"not null;size:100" json:"password" yaml:"password"`
	LastName         string    `gorm:"size:20" json:"lastname" yaml:"lastname"`
	FirstName        string    `gorm:"size:20" json:"firstname" yaml:"firstname"`
	Gender           string    `gorm:"default:other" json:"gender" yaml:"gender"`
	Email            string    `gorm:"unique;not null;size:255" json:"email" yaml:"email"`
	Major            string    `gorm:"default:other" json:"major" yaml:"major"`
	HighestEducation string    `gorm:"default:bachelor" json:"highest_education" yaml:"highestEducation"`
	SubscribeEmail   bool      `gorm:"default:true" json:"subscribe_email" yaml:"subscribeEmail"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"-"`
}

/*
1. The `gorm:"primaryKey" tag:
| index_name | bigint(20) unsigned | NO   | PRI | NULL    | auto_increment |

2. gorm.Model creates these fields:
+------------+---------------------+------+-----+---------+----------------+
| Field      | Type                | Null | Key | Default | Extra          |
+------------+---------------------+------+-----+---------+----------------+
| id         | bigint(20) unsigned | NO   | PRI | NULL    | auto_increment |
| created_at | datetime(3)         | YES  |     | NULL    |                |
| updated_at | datetime(3)         | YES  |     | NULL    |                |
| deleted_at | datetime(3)         | YES  | MUL | NULL    |                |
+------------+---------------------+------+-----+---------+----------------+
*/
