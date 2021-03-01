package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var (
	config *Config
)

func GetConfig() Config {
	tmp := *config

	adminEmails := make([]string, len(config.Admin.Email))
	for i, a := range config.Admin.Email {
		adminEmails[i] = a
	}
	return tmp
}

func (c *Config) load(configFilePath string) error {
	file, err := os.Open(configFilePath)
	if err != nil {
		return err
	}

	b, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if len(b) != 0 {
		return yaml.Unmarshal(b, c)
	}
	return nil
}

func (c *Config) check() *ConfigError {
	if c.App.Name == "" {
		return &ConfigError{errType: "app", err: "name"}
	} else if c.Database.Driver == "" {
		return &ConfigError{errType: "database", err: "driver"}
	} else if c.Database.Username == "" {
		return &ConfigError{errType: "database", err: "username"}
	} else if c.Database.Password == "" {
		return &ConfigError{errType: "database", err: "password"}
	} else if len(c.Admin.Email) == 0 {
		return &ConfigError{errType: "admin", err: "email"}
	}
	return nil
}

func (c *Config) setDefaultValue() {
	if c.App.Port == "0" {
		c.App.Port = "8080"
		logrus.Infof("App port is not set in the config file. Set to default value 8080")
	}

	if c.Database.Host == "" {
		c.Database.Host = "127.0.0.1"
		logrus.Infof("Database host is not set in the config file. Set to default value 127.0.0.1")
	}

	if c.Database.Port == "0" {
		c.Database.Port = "3306"
		logrus.Infof("Database port is not set in the config file. Set to default value 3306")
	}
}

func (c *Config) setOverwriteValue() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost != "" {
		logrus.Info("Database host is overwrote by env DB_HOST. Set to " + dbHost + ".")
		c.Database.Host = dbHost
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort != "" {
		logrus.Info("Database port is overwrote by env DB_PORT. Set to " + dbPort + ".")
		c.Database.Port = dbPort
	}

	appPort := os.Getenv("APP_PORT")
	if appPort != "" {
		logrus.Info("App port is overwrote by env APP_PORT. Set to " + appPort + ".")
		c.App.Port = appPort
	}

}

func Initial(configFilePath string) error {
	config = &Config{}

	if err := config.load(configFilePath); err != nil {
		return err
	}

	if err := config.check(); err != nil {
		return err
	}

	config.setOverwriteValue()
	config.setDefaultValue()
	return nil
}
