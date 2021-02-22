package databases

import (
	"os"
	"strconv"

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
}

func Initial() error {
	var err error

	config := config.GetConfigDatabase()

	switch driver := config.Driver; driver {
	case "mysql":
		// dsn := "user:pwd@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
		host := os.Getenv("DB_HOST")
		if host != "" {
			config.Host = host
		}
		dsn := config.Username + ":" + config.Password + "@tcp(" + config.Host + ":" + strconv.Itoa(config.Port) + ")/" + config.Database + "?charset=utf8mb4&parseTime=True"
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
	return nil
}
