package databases

import (
	"strconv"

	"github.com/cwhuang29/article-sharing-website/config"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/sirupsen/logrus"
)

func GetResetPasswordToken(token string) (passwordToken models.Password) {
	db.Table("passwords").Where("token = ?", token).Preload("User").Find(&passwordToken)
	return
}

func CountUserResetPasswordTokens(id int) int {
	tx := db.Table("passwords").Where("user_id = ?", id).Find(&[]models.Password{})
	return int(tx.RowsAffected)
}

func InsertResetPasswordToken(id int, token string, maxAge int) bool {
	var passwordToken models.Password

	if tx := db.Table("passwords").Where("token  = ?", token).Find(&passwordToken); tx.RowsAffected != 0 {
		return false // There are duplicate tokens
	}

	passwordToken = models.Password{
		UserID: strconv.Itoa(id),
		Token:  token,
		MaxAge: maxAge,
	}

	if err := db.Create(&passwordToken).Error; err != nil {
		logrus.Error(err.Error())
		return false
	}
	return true
}

func DeletePasswordToken(token string) {
	db.Delete(&models.Password{}, "token = ?", token)
}

func DeleteExpiredPasswordTokens(id int) {
	driver := config.GetCopy().Driver
	switch driver {
	case "mysql":
		db.Exec("DELETE FROM passwords WHERE user_id = \"" + strconv.Itoa(id) + "\" AND created_at + max_age - now() < 0")
	case "sqlite":
		db.Exec("DELETE FROM passwords WHERE user_id = \"" + strconv.Itoa(id) + "\" AND created_at + max_age - strftime('%s', 'now') < 0")
	default:
		panic("DB driver is incorrect!")
	}
}
