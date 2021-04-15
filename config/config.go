package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var (
	config *Config
	envs   []env
)

// Return a copy of config
func GetConfig() Config {
	tmp := *config

	adminEmails := make([]string, len(config.Admin.Email))
	for i, a := range config.Admin.Email {
		adminEmails[i] = a
	}
	tmp.Admin.Email = adminEmails
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
	if c.Database.Driver == "" {
		return &ConfigError{errType: "database", err: "driver"}
	}

	if c.Database.Username == "" {
		return &ConfigError{errType: "database", err: "username"}
	}

	if c.Database.Password == "" {
		return &ConfigError{errType: "database", err: "password"}
	}

	if len(c.Admin.Email) == 0 {
		return &ConfigError{errType: "admin", err: "email"}
	}

	return nil
}

func (c *Config) setDefaultValue() {
	if c.App.Url == "" {
		c.App.Url = "http://127.0.0.1"
		logrus.Info("app.Url is not set in the config file. Set to default value http://127.0.0.1")
	}

	if c.App.HttpPort == "" && c.App.HttpsPort == "" {
		c.App.HttpPort = "8080"
		logrus.Info("Both app.httpPort and app.httpsPort are not set in the config file. Set app.HttpsPort to default value 8080")
	}

	if c.Database.Host == "" {
		c.Database.Host = "127.0.0.1"
		logrus.Info("database.host is not set in the config file. Set to default value 127.0.0.1")
	}

	if c.Database.Port == "" {
		var p string

		if c.Database.Driver == "mysql" {
			p = "3306"
		}

		c.Database.Port = p
		logrus.Info("database.port is not set in the config file. Set to default value " + p)
	}
}

func (c *Config) setOverwriteValue() {
	envs = []env{
		{"WEB_DB_HOST", "database.host", &config.Database.Host},
		{"WEB_DB_PORT", "database.port", &config.Database.Port},
		{"WEB_APP_URL", "app.url", &config.App.Url},
		{"WEB_APP_HTTP_PORT", "app.httpPort", &config.App.HttpPort},
		{"WEB_APP_HTTPS_PORT", "app.httpsPort", &config.App.HttpsPort},
		{"WEB_EMAIL_SENDER", "app.email.sender", &config.Email.Sender},
		{"WEB_EMAIL_REGION", "app.email.region", &config.Email.Region},
		{"WEB_EMAIL_NUM_PER_DAY", "app.email.numPerDay", &config.Email.NumPerDay},
		{"WEB_EMAIL_NUM_PER_SEC", "app.email.numPerSec", &config.Email.NumPerSec},
	}

	for _, e := range envs {
		value := os.Getenv(e.key)
		if value != "" {
			*e.target = value
			logrus.Info(e.msg + " is overwrote by env " + e.key + ". Set to " + value + ".")
		}
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
