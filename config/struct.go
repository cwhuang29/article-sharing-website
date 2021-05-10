package config

import (
	"fmt"
	"strings"
)

type configError struct {
	errType string
	err     string
}

func (e *configError) Error() string {
	return fmt.Sprintf("%v %v can't be empty", strings.Title(strings.ToLower(e.errType)), strings.Title(strings.ToLower(e.err)))
}

type env struct {
	key    string
	msg    string
	target *string
}

type config struct {
	App      `yaml:"app,omitempty"`
	Database `yaml:"database,omitempty"`
	Admin    `yaml:"admin,omitempty"`
	Email    `yaml:"email,omitempty"`
}

type App struct {
	Name      string `yaml:"name,omitempty"`
	Url       string `yaml:"url,omitempty"`
	HttpPort  string `yaml:"httpPort,omitempty"`
	HttpsPort string `yaml:"httpsPort,omitempty"`
	Debug     string `yaml:"debug,omitempty"`
	Log       string `yaml:"log,omitempty"`
}

type Database struct {
	Driver   string `yaml:"driver,omitempty"`
	Host     string `yaml:"host,omitempty"`
	Port     string `yaml:"port,omitempty"`
	Database string `yaml:"database,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}

type Admin struct {
	Email []string `yaml:"email,omitempty"`
}

type Email struct {
	Sender    string `yaml:"sender,omitempty"`
	Region    string `yaml:"region,omitempty"`
	NumPerDay string `yaml:"numPerDay,omitempty"`
	NumPerSec string `yaml:"numPerSec,omitempty"`
}
