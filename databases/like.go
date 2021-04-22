package databases

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/sirupsen/logrus"
)

/*
 * 0: User doesn't like this article
 * 1: User has liked this article
 */

func CountUserLikes(userId int) int {
	return int(db.Model(&models.User{ID: userId}).Association("LikedArticles").Count())
}

func GetUserLikedArticles(userId, offset, limit int, isAdmin bool) (articles []models.Article) {
	switch isAdmin {
	case true:
		if err := db.Model(&models.User{ID: userId}).Preload("Tags").Order("id desc").Limit(limit).Offset(offset).Association("LikedArticles").Find(&articles); err != nil {
			logrus.Error(err.Error())
		}
	case false:
		if err := db.Model(&models.User{ID: userId}).Preload("Tags").Order("id desc").Limit(limit).Offset(offset).Where("admin_only = ?", false).Association("LikedArticles").Find(&articles); err != nil {
			logrus.Error(err.Error())
		}
	}
	return
}

func GetLikeStatus(userId, articleId int) int {
	user := models.User{}
	article := models.Article{ID: articleId}

	if err := db.Model(&article).Where("user_id = ?", userId).Association("LikedUsers").Find(&user); err != nil {
		logrus.Error(err.Error())
		return 0
	}

	if user.ID == 0 {
		return 0
	}
	return 1
}

func UpdateLikeStatus(userId, articleId, isLiked int) bool {
	user := models.User{ID: userId}
	article := models.Article{ID: articleId}

	switch isLiked {
	case 0:
		if err := db.Model(&article).Association("LikedUsers").Delete(&user); err != nil {
			logrus.Error(err.Error())
			return false
		}
	case 1:
		if err := db.Model(&article).Association("LikedUsers").Append(&user); err != nil {
			logrus.Error(err.Error())
			return false
		}
	default:
		return false
	}
	return true
}
