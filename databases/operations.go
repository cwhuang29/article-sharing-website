package databases

import (
	"fmt"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/sirupsen/logrus"
	"time"
)

func IsArticleExists(id int) (succeed bool) {
	var article models.Article
	succeed = true
	tx := db.Where("id = ?", id).First(&article)
	if tx.RowsAffected == 0 {
		succeed = false
	}
	return
}

func InsertArticle(article models.Article) (int, bool) {
	if err := db.Create(&article).Error; err != nil { // Create returns a clone of DB and Error field is set in that clone object
		logrus.Error(err.Error())
		return -1, false
	}

	return article.ID, true
}

func DeleteArticle(id int) bool {
	if err := db.Where("id = ?", id).Delete(&models.Article{}).Error; err != nil {
		logrus.Error(err.Error())
		return false
	}
	return true
}

func ReplaceArticle(article models.Article) (int, bool) {
	// db.Save(&article) will save all fields when performing the Updating SQL (including created_at)
	// if err := db.Model(&article).Where("id = ?", article.ID).Updates(models.Article{ // Where part can be omitted cause article.ID is the primary key
	//     Title:       article.Title,
	//     Subtitle:    article.Subtitle,
	// }).Error; err != nil { }

	if err := db.Model(&article).Where("id = ?", article.ID).Updates(article).Error; err != nil {
		logrus.Error(err.Error())
		return -1, false
	}

	return article.ID, true
}

func GetArticlesInNDays(startTime time.Time) (articles []models.Article) {
	db.Table("articles").Order("id desc").Where("created_at > ?", startTime).Find(&articles)
	return
}

func GetArticlesList(category string, offset int, limit int) (articles []models.Article, err error) {
	if offset < 0 {
		err = fmt.Errorf("<div><p><strong>Invalid parameter</strong></p><p>offset should not be negative.</p></div>")
		return
	} else if limit <= 0 {
		err = fmt.Errorf("<div><p><strong>Invalid parameter</strong></p><p>limit should be greater than zero.</p></div>")
		return
	}

	db.Limit(limit).Offset(offset).Order("id desc").Where("category = ?", category).Find(&articles) // Not found is not an error (just no articles)
	return
}

func GetArticleFullContent(id int) (article models.Article, succeed bool) {
	succeed = true
	tx := db.Where("id = ?", id).First(&article)
	if tx.RowsAffected == 0 {
		succeed = false
	}
	return
}
