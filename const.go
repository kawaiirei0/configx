package config

import "time"

// 默认常量配置
const (
	OptionFilename        = "config"
	OptionFileType        = "yaml"
	OptionPath            = "./configs"
	OptionEnv             = "dev"
	OptionDebounceDur     = 800 * time.Millisecond
	OptionDateMillisecond = OptionTimeDuration(time.Millisecond)
)
