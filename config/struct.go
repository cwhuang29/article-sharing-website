package config

import (
	"fmt"
	"strings"
)

type ConfigError struct {
	errType string
	err     string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("%v %v can't be empty", strings.Title(strings.ToLower(e.errType)), strings.Title(strings.ToLower(e.err)))
}

type Config struct {
	App      `yaml:"app,omitempty"`
	Database `yaml:"database,omitempty"`
	Admin    `yaml:"admin,omitempty"`
}

type App struct {
	Name  string `yaml:"name,omitempty"`
	Port  string `yaml:"port,omitempty"`
	Debug bool   `yaml:"debug,omitempty"`
	Log   string `yaml:"log,omitempty"`
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
