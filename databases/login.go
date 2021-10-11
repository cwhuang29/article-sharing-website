package databases

import (
	"strconv"

	"github.com/cwhuang29/article-sharing-website/config"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/sirupsen/logrus"
)

func GetLoginCredentials(id int) (loginSession []models.Login) {
	db.Table("logins").Preload("User").Where("user_id  = ?", id).Find(&loginSession) // Preload users cause we'll need them for validation later
	return
}

func InsertLoginToken(id int, token string, maxAge int) bool {
	var loginToken models.Login

	if tx := db.Table("logins").Where("token  = ?", token).Find(&loginToken); tx.RowsAffected != 0 {
		return false
	}

	loginToken = models.Login{
		UserID: strconv.Itoa(id),
		Token:  token,
		MaxAge: maxAge,
	}

	if err := db.Create(&loginToken).Error; err != nil {
		logrus.Error(err.Error())
		return false
	}
	return true
}

func DeleteLoginToken(userID int, token string) {
	db.Delete(&models.Login{}, "user_id = ? and token = ?", userID, token)
}

func DeleteExpiredLoginTokens(id int) {
	driver := config.GetCopy().Driver
	switch driver {
	case "mysql":
		db.Exec("DELETE FROM logins WHERE user_id = \"" + strconv.Itoa(id) + "\" AND last_login + max_age - now() < 0")
	case "sqlite":
		db.Exec("DELETE FROM logins WHERE user_id = \"" + strconv.Itoa(id) + "\" AND last_login + max_age - strftime('%s', 'now') < 0")
	default:
		panic("DB driver is incorrect!")
	}
}
