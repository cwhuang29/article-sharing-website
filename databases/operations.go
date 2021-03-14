package databases

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func IsArticleExists(id int) (succeed bool) {
	if tx := db.Where("id = ?", id).First(&models.Article{}); tx.RowsAffected == 0 {
		return
	}
	succeed = true
	return
}

func GetArticle(id int) (article models.Article) {
	db.Preload("Tags").Where("id = ?", id).First(&article)
	return
}

func GetArticleWithoutTags(id int) (article models.Article) {
	db.Where("id = ?", id).First(&article)
	return
}

func GetArticlesInATimePeriod(start, end time.Time) (articles []models.Article) {
	db.Order("id desc").Preload("Tags").Where("created_at >= ? and created_at < ?", start, end).Find(&articles)
	return
}

func GetArticlesList(category string, offset int, limit int) (articles []models.Article) {
	db.Limit(limit).Offset(offset).Order("id desc").Preload("Tags").Where("category = ?", category).Find(&articles) // No articles is not an error
	return
}

func insertArticlesAndTagsAssociation(tx *gorm.DB, article models.Article, tag models.Tag) error {
	// Ths association created is based on primary keys (which are article.ID and tag.ID)
	return tx.Model(&article).Association("Tags").Append(&tag)
}

func getArticlesAndTagsAssociation(id int) []models.Tag {
	article := models.Article{ID: id}

	db.Table("articles").Preload("Tags").Find(&article)
	return article.Tags
}

func deleteArticlesAndTagsAssociation(article models.Article, tag models.Tag) bool {
	if err := db.Model(&article).Association("Tags").Delete(tag); err != nil {
		logrus.Error(err.Error())
		return false
	}
	return true
}

func insertTag(tag models.Tag) (models.Tag, error) {
	if err := db.FirstOrCreate(&tag, models.Tag{Value: tag.Value}).Error; err != nil {
		logrus.Error(err.Error())
		return models.Tag{}, err
	}

	return tag, nil
}

func insertArticle(article models.Article) (models.Article, error) {
	if err := db.Create(&article).Error; err != nil {
		logrus.Error(err.Error())
		return models.Article{}, err
	}

	return article, nil
}

func updateArticle(article models.Article) (models.Article, error) {
	// db.Save(&article) will save all fields when performing the Updating SQL (including created_at)
	// if err := db.Model(&article).Where("id = ?", article.ID).Updates(models.Article{ // Update only several keys
	//     Title:       article.Title,
	//     Subtitle:    article.Subtitle,
	// }).Error; err != nil { }

	if err := db.Model(&article).Where("id = ?", article.ID).Updates(article).Error; err != nil { //  Where clause can be omitted cause article.ID is the primary key
		logrus.Error(err.Error())
		return models.Article{}, err
	}

	return article, nil
}

func SubmitArticle(article models.Article, action string) (int, bool) {
	/*
	 * We can create articles, tags and their associations (in the articles_tags table) directly by db.Create(&article).
	 * However, when different articles have same tags, those tags will have duplicate entries in the tag table
	 * Follow up: the duplicate values in tag table may be fixed by adding `gorm:"unique"` into struct's definition,
	 * but I prefer to keep this function as it is a good example for how to use transactions
	 */

	if action != "create" && action != "update" {
		return -1, false
	}

	var err error
	var newArticle models.Article
	var currTags, prevTags []models.Tag

	tx := db.Begin()
	if tx.Error != nil {
		return -1, false
	}

	currTags = article.Tags
	article.Tags = []models.Tag{}

	if action == "create" {
		newArticle, err = insertArticle(article)
	} else if action == "update" {
		prevTags = getArticlesAndTagsAssociation(article.ID)
		newArticle, err = updateArticle(article)
	}
	if err != nil {
		return -1, false
	}

	for _, tag := range currTags {
		var t models.Tag
		t, err = insertTag(tag)
		if err != nil {
			tx.Rollback()
			return -1, false
		}

		err = insertArticlesAndTagsAssociation(tx, newArticle, t)
		if err != nil {
			tx.Rollback()
			return -1, false
		}
	}

	res := tx.Commit()
	if res.Error != nil {
		return -1, false
	}

	if action == "update" {
		for _, prevTag := range prevTags {
			find := false
			for _, currTag := range currTags {
				if prevTag.Value == currTag.Value {
					find = true
					break
				}
			}
			if find == false {
				deleteArticlesAndTagsAssociation(newArticle, prevTag)
			}
		}
	}

	return newArticle.ID, true
}

func DeleteArticle(id int) bool {
	article := GetArticleWithoutTags(id)
	if article.ID == 0 {
		return false
	}

	// This won't remove any tags (the ID field should not be empty): db.Where("id = ?", id).Select("Tags").Delete(&models.Article{})
	if err := db.Select("Tags").Delete(&article).Error; err != nil {
		logrus.Error(err.Error())
		return false
	}
	return true
}
