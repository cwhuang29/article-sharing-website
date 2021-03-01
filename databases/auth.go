package databases

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/sirupsen/logrus"
	"time"
)

func GetUser(email string) (user models.User) {
	db.Table("users").Where("email = ?", email).Find(&user)
	return
}

func IsAdminUser(email string) bool {
	var user models.User

	tx := db.Table("admins").Where("email = ?", email).Find(&user)
	if tx.RowsAffected == 0 {
		return false
	}
	return true
}

func InsertUserToDB(user models.User) (int, bool) {
	if err := db.Create(&user).Error; err != nil {
		logrus.Error(err.Error())
		return -1, false
	} // Create returns a clone of DB and Error field is set in that clone object

	return user.ID, true
}

func InsertLoginToken(email string, token string, maxAge int) {
	var loginToken models.Login

	tx := db.Table("logins").Where("email  = ?", email).Find(&loginToken)
	if tx.RowsAffected == 0 {
		newTx := models.Login{
			Email:     email,
			Token:     token,
			MaxAge:    maxAge,
			LastLogin: time.Now().UTC(),
		}
		db.Create(&newTx)
	}
	loginToken.Token = token
	loginToken.MaxAge = maxAge
	loginToken.LastLogin = time.Now().UTC()
	db.Save(&loginToken)
}

func DeleteLoginToken(email string) {
	var loginToken models.Login

	tx := db.Table("logins").Where("email  = ?", email).Find(&loginToken)
	if tx.RowsAffected != 0 {
		db.Delete(&models.Login{}, "email = ?", email)
	}
}

func GetLoginCredentials(email string) (loginSession models.Login) {
	db.Table("logins").Where("email  = ?", email).Find(&loginSession)
	return
}
