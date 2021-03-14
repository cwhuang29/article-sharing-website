package models

import (
	"gorm.io/gorm"
)

/*
 * func (u *User) AfterDelete(tx *gorm.DB) (err error) {
 *   if u.Confirmed {
 *     tx.Model(&Address{}).Where("user_id = ?", u.ID).Update("invalid", false)
 *   }
 *   return
 * }
 */

func (u *User) BeforeDelete(tx *gorm.DB) (err error) {
	/*
	 * Note: how the struct been deleted does matters
	 * db.Delete(&models.User{}, "email = ?", user.Email) - u is an empty struct
	 * db.Delete(&models.User{Email: user.Email}) - u has a non-zero value field Email
	 * So it is better to pass an non-zero struct into db.Delete()
	 */

	// By using Migrator().CreateConstraint() to create foreign keys, Login records will be deleted when the corresponding user be deleted
	// return tx.Model(&Login{}).Where("user_id = ?", u.ID).Delete(&Login{}).Error
	return
}
