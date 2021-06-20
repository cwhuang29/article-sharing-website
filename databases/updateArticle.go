package databases

import (
	"errors"

	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/sirupsen/logrus"

	"gorm.io/gorm"
)

func UpdateTagsStats(tagValue string) {
	var tag models.Tag

	db.Where("value = ?", tagValue).First(&tag)
	if err := db.Model(&tag).Updates(models.Tag{Views: tag.Views + 1}).Error; err != nil {
		logrus.Error(err.Error())
	}
}

func DeleteArticle(id int, isAdmin bool) bool {
	article := GetArticleWithoutTags(id, isAdmin)
	if article.ID == 0 {
		return false
	}

	// The following query won't remove any tags (the primary key, i.e. ID field, of the struct should NOT be empty):
	//     db.Where("id = ?", id).Select("Tags").Delete(&models.Article{})
	if err := db.Select("Tags").Delete(&article).Error; err != nil {
		logrus.Error(err.Error())
		return false
	}
	return true
}

func insertArticlesAndTagsAssociation(tx *gorm.DB, article models.Article, tag models.Tag) error {
	// Ths association created (a new record in the join table articles_tags) is based on primary keys (article.ID and tag.ID)
	return tx.Model(&article).Association("Tags").Append(&tag)
}

func getArticlesAndTagsAssociation(id int) []models.Tag {
	article := models.Article{ID: id}

	db.Table("articles").Preload("Tags").Find(&article)
	return article.Tags
}

func deleteArticlesAndTagsAssociation(article models.Article, tag models.Tag) bool {
	if err := db.Model(&article).Association("Tags").Delete(&tag); err != nil {
		logrus.Error(err.Error())
		return false
	}
	return true
}

func insertTag(tx *gorm.DB, tag models.Tag) (models.Tag, error) {
	if err := tx.FirstOrCreate(&tag, models.Tag{Value: tag.Value}).Error; err != nil {
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
	// When update with struct, GORM will only update non-zero fields. Use map type variable to update or Select() to specify fields to update
	// Where clause can be omitted since article.ID is the primary key
	a := map[string]interface{}{
		"title":        article.Title,
		"subtitle":     article.Subtitle,
		"authors":      article.Authors,
		"release_date": article.ReleaseDate,
		"category":     article.Category,
		"outline":      article.Outline,
		"content":      article.Content,
		"admin_only":   article.AdminOnly,
	}

	// If user doesn't upload new cover photo while editing articles, we have to keep the original one
	if article.CoverPhoto != "" {
		a["cover_photo"] = article.CoverPhoto
	}

	if err := db.Model(&article).Where("id = ?", article.ID).Updates(a).Error; err != nil {
		logrus.Error(err.Error())
		return models.Article{}, err
	}

	return article, nil
}

func SubmitArticle(article models.Article, action string) (newArticleID int, succeed bool) {
	/*
	 * We can create articles, tags and their associations (in the articles_tags table) directly by db.Create(&article).
	 * However, when different articles have same tags, those tags will have duplicate entries in the tag table
	 * Follow up: the duplicate values in tag table may be fixed by adding `gorm:"unique"` into struct's definition,
	 * but I prefer to keep this function as it is a good example for how to use transactions
	 */

	newArticleID = -1
	succeed = false

	var err error
	var newArticle models.Article
	var currTags, prevTags []models.Tag

	tx := db.Begin() // Use 'tx' from this point, not 'db'
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return
	}

	currTags = article.Tags
	article.Tags = []models.Tag{}

	switch action {
	case "create":
		newArticle, err = insertArticle(article)
	case "update":
		prevTags = getArticlesAndTagsAssociation(article.ID)
		newArticle, err = updateArticle(article)
	default:
		err = errors.New("")
	}

	if err != nil {
		return
	}

	for _, tag := range currTags {
		var t models.Tag
		t, err = insertTag(tx, tag)
		if err != nil {
			tx.Rollback()
			return
		}

		err = insertArticlesAndTagsAssociation(tx, newArticle, t)
		if err != nil {
			tx.Rollback()
			return
		}
	}

	res := tx.Commit()
	if res.Error != nil {
		return
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
