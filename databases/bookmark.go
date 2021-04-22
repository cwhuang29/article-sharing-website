package databases

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/sirupsen/logrus"
)

/*
 * 0: User doesn't bookmark this article
 * 1: User has bookmarked this article
 */

func CountUserBookmarks(userId int) int {
	// db.Model(&user).Where("user_id = ?", userId) doesn't work since the primary key of user struct is empty
	return int(db.Model(&models.User{ID: userId}).Association("BookmarkedArticles").Count())
}

func GetUserBookmarkArticles(userId, offset, limit int, isAdmin bool) (articles []models.Article) {
	// Note that Limit(), Offset(), and Where() have to put in front of Association()
	switch isAdmin {
	case true:
		if err := db.Model(&models.User{ID: userId}).Preload("Tags").Order("id desc").Limit(limit).Offset(offset).Association("BookmarkedArticles").Find(&articles); err != nil {
			logrus.Error(err.Error())
		}
	case false:
		if err := db.Model(&models.User{ID: userId}).Preload("Tags").Order("id desc").Limit(limit).Offset(offset).Where("admin_only = ?", false).Association("BookmarkedArticles").Find(&articles); err != nil {
			logrus.Error(err.Error())
		}
	}
	return
}

func GetBookmarkStatus(userId, articleId int) int {
	user := models.User{}
	article := models.Article{ID: articleId}

	// Note: The argument of Association() is "BookmarkedUsers" (the name of the member in models.Article)
	if err := db.Model(&article).Where("user_id = ?", userId).Association("BookmarkedUsers").Find(&user); err != nil {
		logrus.Error(err.Error())
		return 0
	}

	if user.ID == 0 { // If there is no matched record, the argument of Find() remains untouched
		return 0
	}
	return 1
}

func UpdateBookmarkStatus(userId, articleId, isBookmarked int) bool {
	user := models.User{ID: userId}
	article := models.Article{ID: articleId}

	switch isBookmarked {
	case 0:
		// Remove the relationship between source & arguments if exists, only delete the reference
		if err := db.Model(&article).Association("BookmarkedUsers").Delete(&user); err != nil {
			logrus.Error(err.Error())
			return false
		}
	case 1:
		// Append new associations for many to many, has many, replace current association for has one, belongs to
		if err := db.Model(&article).Association("BookmarkedUsers").Append(&user); err != nil {
			logrus.Error(err.Error())
			return false
		}
	default:
		return false
	}
	return true
}
