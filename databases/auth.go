package databases

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/sirupsen/logrus"
	"strconv"
)

func IsAdminUser(email string) bool {
	var user models.User

	tx := db.Table("admins").Where("email = ?", email).Find(&user)
	if tx.RowsAffected == 0 {
		return false
	}
	return true
}

func GetUser(email string) (user models.User) {
	db.Table("users").Where("email = ?", email).Find(&user)
	return
}

func GetLoginCredentials(email string) (loginSession []models.Login) {
	user := GetUser(email)
	db.Table("logins").Where("user_id  = ?", user.ID).Preload("User").Find(&loginSession) // Preload users when find loginSession
	return
}

func InsertUser(user models.User) (int, bool) {
	if err := db.Create(&user).Error; err != nil {
		logrus.Error(err.Error())
		return -1, false
	} // Create returns a clone of DB and Error field is set in that clone object

	return user.ID, true
}

func InsertLoginToken(email, token string, maxAge int) {
	user := GetUser(email)
	loginToken := models.Login{User: user, Token: token, MaxAge: maxAge}
	db.Create(&loginToken)

	// Notice: In the beginning there is only one token per user,
	// which is not user-friendly because when user login with cellphone, their web accounts will be logout.
	// var loginToken models.Login
	// tx := db.Table("logins").Where("email  = ?", email).Find(&loginToken)
	// if tx.RowsAffected == 0 {
	//     newTx := models.Login{
	//         Email:     email,
	//         Token:     token,
	//         MaxAge:    maxAge,
	//         LastLogin: time.Now().UTC(),
	//     }
	//     db.Create(&newTx)
	// }
	// loginToken.Token = token
	// loginToken.MaxAge = maxAge
	// loginToken.LastLogin = time.Now().UTC()
	// db.Save(&loginToken)
}

func DeleteLoginToken(email, token string) {
	user := GetUser(email)
	db.Delete(&models.Login{}, "user_id = ? and token = ?", user.ID, token)
}

func DeleteExpiredLoginTokens(email string) {
	user := GetUser(email)
	db.Exec("delete from logins where user_id = \"" + strconv.Itoa(user.ID) + "\" and last_login + max_age - now() < 0")
}
