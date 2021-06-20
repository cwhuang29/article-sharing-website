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

func getMysqlDSN(cfg config.Database) string {
	// Example: "user:pwd@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	return cfg.Username + ":" + cfg.Password + "@tcp(" + cfg.Host + ":" + cfg.Port + ")/" + cfg.Database + "?charset=utf8mb4&parseTime=True"
}

func connect(cfg config.Database) (err error) {
	switch driver := cfg.Driver; driver {
	case "mysql":
		dsn := getMysqlDSN(cfg)
		if db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}); err != nil {
			return err
		}
	case "sqlite":
		if db, err = gorm.Open(sqlite.Open("tmp.db"), &gorm.Config{}); err != nil {
			return err
		}
	default:
		panic("Please select a correct database driver (mysql or sqlite).")
	}
	return
}

func createTables() {
	// See https://gorm.io/docs/migration.html
	if !(db.Migrator().HasTable(&models.User{}) && db.Migrator().HasTable(&models.Article{}) && db.Migrator().HasTable(&models.Tag{})) {
		db.AutoMigrate(&models.User{}, &models.Article{}, &models.Tag{})
	}
	if !(db.Migrator().HasTable(&models.Login{})) {
		db.Migrator().CreateTable(&models.Login{})
	}
	if !(db.Migrator().HasTable(&models.Password{})) {
		db.Migrator().CreateTable(&models.Password{})
	}
	if !(db.Migrator().HasTable(&models.Admin{})) {
		db.Migrator().CreateTable(&models.Admin{})
	}
}

func createConstraints() {
	if !db.Migrator().HasConstraint(&models.Login{}, "User") {
		db.Migrator().CreateConstraint(&models.Login{}, "User")
	}
	if !db.Migrator().HasConstraint(&models.Password{}, "User") {
		db.Migrator().CreateConstraint(&models.Password{}, "User")
	}
}

func registerAdminEmail(emails []string) {
	for _, email := range emails {
		obj := models.Admin{Email: email}
		db.Create(&obj)
	}
}

func Initial() (err error) {
	cfg := config.GetCopy()

	if err = connect(cfg.Database); err != nil {
		return
	}
	createTables()
	createConstraints()
	registerAdminEmail(cfg.Admin.Email)
	return
}
