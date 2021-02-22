package main

import (
	"fmt"
	"github.com/cwhuang29/article-sharing-website/config"
	"github.com/cwhuang29/article-sharing-website/databases"
	"github.com/cwhuang29/article-sharing-website/databases/models"
	"github.com/cwhuang29/article-sharing-website/handlers"
	"github.com/cwhuang29/article-sharing-website/routers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("Parsing config file: %s ...", viper.ConfigFileUsed())
}

func setupConfig() {
	if err := config.Initial(viper.ConfigFileUsed()); err != nil {
		logrus.Fatal(err)
	}
}

func setupDatabase() {
	logrus.Infof("Setting database ...")
	if err := databases.Initial(); err != nil {
		logrus.Fatal(err)
	}
}

func createFakeUser() {
	pwd, _ := handlers.HashPassword("a1234567")
	user := models.User{
		Name:      "admin",
		Password:  string(pwd),
		LastName:  "Huang",
		FirstName: "Jimmy",
		Gender:    "male",
		Email:     "admin@gmail.com",
		Major:     "Pharmacy",
	}
	_, err := databases.InsertUserToDB(user)
	if err != nil {
		fmt.Println(err)
		return
	}
}

var rootCmd = &cobra.Command{
	Use:   "article-sharing-website",
	Short: "article-sharing-website",
	Long:  "An article sharing website developed by Go.",
	Run: func(cmd *cobra.Command, args []string) {
		setupConfig()
		setupDatabase()
		createFakeUser()
		routers.Router()
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err) // This will call os.Exit(255)
	}
}
