package models

import (
	"time"
)

type User struct {
	ID                 int       `gorm:"primaryKey" json:"id" yaml:"id"`
	FirstName          string    `gorm:"size:50" json:"first_name" yaml:"firstName"`
	LastName           string    `gorm:"size:50" json:"last_name" yaml:"lastName"`
	Password           string    `gorm:"not null;size:100" json:"password" yaml:"password"`
	Email              string    `gorm:"unique;not null;size:100" json:"email" yaml:"email"`
	Gender             string    `gorm:"default:other" json:"gender" yaml:"gender"`
	Major              string    `gorm:"default:other" json:"major" yaml:"major"`
	SubscribeEmail     bool      `gorm:"default:true" json:"subscribe_email" yaml:"subscribeEmail"`
	Admin              bool      `gorm:"default:false" json:"admin" yaml:"admin"`
	BookmarkedArticles []Article `gorm:"many2many:users_articles_bookmark;"`
	LikedArticles      []Article `gorm:"many2many:users_articles_like;"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"-"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"-"`
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

3. Without specifying size, the column type will be longtext
*/
