package databases

import (
	"time"

	"github.com/cwhuang29/article-sharing-website/databases/models"
)

func IsArticleExists(id int, isAdmin bool) (succeed bool) {
	var article models.Article

	if tx := db.Where("id = ?", id).First(&article); tx.RowsAffected == 0 {
		return
	}

	if isAdmin == true || isAdmin == false && article.AdminOnly == false {
		succeed = true
	}
	return
}

func GetArticleTags(article models.Article) (tags []models.Tag) {
	db.Model(&article).Association("Tags").Find(&tags)
	return
}

func GetArticle(id int, isAdmin bool) (article models.Article) {
	db.Preload("Tags").Where("id = ?", id).First(&article)

	if isAdmin == false && article.AdminOnly == true {
		article = models.Article{}
	}
	return
}

func GetArticleWithoutTags(id int, isAdmin bool) (article models.Article) {
	db.Where("id = ?", id).First(&article)

	if isAdmin == false && article.AdminOnly == true {
		article = models.Article{}
	}
	return
}

func GetArticlesInATimePeriod(start, end time.Time, isAdmin bool) (articles []models.Article) {
	switch isAdmin {
	case true:
		db.Preload("Tags").Order("id desc").Where("updated_at >= ? and updated_at < ?", start, end).Find(&articles)
	case false:
		db.Preload("Tags").Order("id desc").Where("updated_at >= ? and updated_at < ? and admin_only = ?", start, end, false).Find(&articles)
	}
	return
}

func GetArticles(offset, limit int, isAdmin bool) (articles []models.Article) {
	switch isAdmin {
	case true:
		db.Preload("Tags").Order("id desc").Limit(limit).Offset(offset).Find(&articles)
	case false:
		db.Preload("Tags").Order("id desc").Limit(limit).Offset(offset).Where("admin_only = ?", false).Find(&articles)
	}
	return
}

func GetSameCategoryArticles(category string, offset, limit int, isAdmin bool) (articles []models.Article) {
	switch isAdmin {
	case true:
		db.Preload("Tags").Order("id desc").Limit(limit).Offset(offset).Where("category = ?", category).Find(&articles)
	case false:
		db.Preload("Tags").Order("id desc").Limit(limit).Offset(offset).Where("category = ? and admin_only = ?", category, false).Find(&articles)
	}
	return
}

func GetSameTagArticles(tagValue string, offset, limit int, isAdmin bool) (articles []models.Article) {
	var tags models.Tag

	db.Preload("Articles").Where("value = ?", tagValue).First(&tags)

	for i := len(tags.Articles) - 1; i >= 0; i-- {
		if offset > 0 && (isAdmin || !tags.Articles[i].AdminOnly) {
			offset -= 1
			continue
		}

		if isAdmin || !tags.Articles[i].AdminOnly {
			articles = append(articles, tags.Articles[i])
		}
		if len(articles) == limit {
			break
		}
	}
	return
}
