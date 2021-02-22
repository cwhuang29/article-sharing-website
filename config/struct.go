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
}

type App struct {
	Name   string `yaml:"name,omitempty"`
	Port   int    `yaml:"port,omitempty"`
	Env    string `yaml:"env,omitempty"`
	Debug  bool   `yaml:"debug,omitempty"`
	Author string `yaml:"author,omitempty"`
	Log    string `yaml:"log,omitempty"`
}

type Database struct {
	Driver   string `yaml:"driver,omitempty"`
	Host     string `yaml:"host,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Database string `yaml:"database,omitempty"`
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}
