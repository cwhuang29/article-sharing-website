package databases

import (
	"github.com/cwhuang29/article-sharing-website/config"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func GetDB() *gorm.DB {
	return db
}

func createTable() {
	if !(db.Migrator().HasTable(&models.User{})) {
		db.Migrator().CreateTable(&models.User{})
	}
	if !(db.Migrator().HasTable(&models.Article{})) {
		db.Migrator().CreateTable(&models.Article{})
	}
	if !(db.Migrator().HasTable(&models.Info{})) {
		db.Migrator().CreateTable(&models.Info{})
	}
	if !(db.Migrator().HasTable(&models.Login{})) {
		db.Migrator().CreateTable(&models.Login{})
	}
	if !(db.Migrator().HasTable(&models.Admin{})) {
		db.Migrator().CreateTable(&models.Admin{})
	}
}

func registerAdminEmail(emails []string) {
	for _, email := range emails {
		obj := models.Admin{Email: email}
		db.Create(&obj)
	}
}

func Initial() error {
	var err error

	cfg := config.GetConfig()

	dbConfig := cfg.Database
	switch driver := dbConfig.Driver; driver {
	case "mysql":
		// dsn := "user:pwd@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		dsn := dbConfig.Username + ":" + dbConfig.Password + "@tcp(" + dbConfig.Host + ":" + dbConfig.Port + ")/" + dbConfig.Database + "?charset=utf8mb4&parseTime=True"
		if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}); err != nil {
			return err
		}
	case "sqlite":
		if db, err = gorm.Open(sqlite.Open("tmp.db"), &gorm.Config{}); err != nil {
			return err
		}
	default:
		panic("Please select a correct database driver.")
	}

	createTable()
	registerAdminEmail(cfg.Admin.Email)
	return nil
}
