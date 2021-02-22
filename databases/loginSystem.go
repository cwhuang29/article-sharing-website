package databases

import (
	"fmt"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"time"
)

func InsertUserToDB(user models.User) (int, error) {
	var err error

	tx := db.Create(&user)
	if tx.RowsAffected > 0 {
		return user.ID, nil
	} else {
		err = fmt.Errorf("<div><p><strong>Failed to Register</strong></p><p>Duplicate email</p></div>")
		return -1, err
	}
}

func GetUser(email string) (user models.User, err error) {
	tx := db.Table("users").Where("email = ?", email).Find(&user)
	if tx.RowsAffected == 0 {
		err = fmt.Errorf("<div><p><strong>User not Found</strong></p><p>Please try again</p></div>")
	}
	return
}

func InsertLoginToken(email string, token string, maxAge int) {
	var loginToken models.Login

	tx := db.Table("logins").Where("email  = ?", email).Find(&loginToken)
	if tx.RowsAffected == 0 {
		newTx := models.Login{
			Email:     email,
			Token:     token,
			MaxAge:    maxAge,
			LastLogin: time.Now(),
		}
		db.Create(&newTx)
	}
	loginToken.Token = token
	loginToken.MaxAge = maxAge
	loginToken.LastLogin = time.Now()
	db.Save(&loginToken)
}

func GetLoginToken(email string) string {
	var loginSession models.Login
	db.Table("logins").Where("email  = ?", email).Find(&loginSession)
	return loginSession.Token
}
