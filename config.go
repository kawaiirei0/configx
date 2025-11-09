package configx

import (
	"github.com/kawaiirei0/configx/configure"
)

// Config 应用配置结构
type Config struct {
	App configure.App `mapstructure:"app"`
}

func NewConfig() *Config {
	//configure.DBConfig{}
	return &Config{
		App: configure.App{
			Name:        "qwq",
			Version:     "1.0.0",
			Description: "QWQ App.",
		},
	}
}
