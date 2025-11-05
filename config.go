package config

import (
	"config/configure"
)

// Config 应用配置结构
type Config struct {
	App configure.App `mapstructure:"app"`
}

func NewConfig() *Config {
	//configure.DBConfig{}
	return &Config{}
}
