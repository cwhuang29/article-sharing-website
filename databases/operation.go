package databases

import (
	"fmt"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"time"
)

func InsertArticleToDB(article models.Article) (int, bool) {
	// if err := db.Create(&Article).Error; err != nil { } // Create returns a clone of DB and Error field is set in that clone object
	tx := db.Create(&article)

	if tx.RowsAffected > 0 {
		return article.ID, true
	} else {
		return -1, false
	}
}

func GetArticlesInNDays(startTime time.Time) (articles []models.Article) {
	db.Table("articles").Order("id desc").Where("created_at > ?", startTime).Find(&articles)
	return
}

func GetArticlesList(category string, offset int, limit int) (articles []models.Article, err error) {
	fmt.Println(err, err == nil)
	err = nil
	if offset < 0 {
		err = fmt.Errorf("<div><p><strong>Invalid parameter</strong></p><p>offset should not be negative</p></div>")
		return
	} else if limit <= 0 {
		err = fmt.Errorf("<div><p><strong>Invalid parameter</strong></p><p>limit should be greater than zero</p></div>")
		return
	}

	db.Limit(limit).Offset(offset).Order("id desc").Where("category = ?", category).Find(&articles) // Not found is not an error (just no articles)
	return
}

func GetArticleFullContent(id int) (article models.Article, err error) {
	// db.Table("hits").Select("created_at").Where("service_group_id = ?", id).Where("created_at > ?", thirtyDaysAgo).Find(&hits)
	tx := db.Where("id = ?", id).First(&article)
	if tx.RowsAffected == 0 {
		err = fmt.Errorf("<div><p><strong>Article Not Found</strong></p><p>Article with ID %d not found</p></div>", id)
	}
	return
}
