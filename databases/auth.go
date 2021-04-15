package databases

import (
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/sirupsen/logrus"
)

func IsAdminUser(email string) bool {
	var user models.User

	if tx := db.Table("admins").Where("email = ?", email).Find(&user); tx.RowsAffected == 0 {
		return false
	}
	return true
}

func GetUser(email string) (user models.User) {
	db.Table("users").Where("email = ?", email).Find(&user)
	return
}

func InsertUser(user models.User) (int, bool) {
	if err := db.Create(&user).Error; err != nil {
		logrus.Error(err.Error())
		return -1, false
	} // Create returns a clone of DB and Error field is set in that clone object

	return user.ID, true
}

func UpdatePassword(user models.User, password, token string) bool {
	if err := db.Model(&user).Updates(models.User{Password: password}).Error; err != nil {
		logrus.Error(err.Error())
		return false
	}
	return true
}
