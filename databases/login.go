package databases

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/sirupsen/logrus"
	"strconv"
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

func DeleteLoginToken(token string) {
	// Too complicated (login token is not a critical field so we don't have to treat it such carefully).
	// Besides, the probability of duplicate tokens is super small
	// user := GetUser(email)
	// db.Delete(&models.Login{}, "user_id = ? and token = ?", user.ID, token)
	db.Delete(&models.Login{}, "token = ?", token)
}

func DeleteExpiredLoginTokens(id int) {
	db.Exec("DELETE FROM logins WHERE user_id = \"" + strconv.Itoa(id) + "\" AND last_login + max_age - now() < 0")
}
