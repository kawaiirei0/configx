package config

import "time"

// 默认常量配置
const (
	OptionFilename        = "config.yaml"
	OptionFileType        = "yaml"
	OptionFilepath        = "./configs"
	OptionEnv             = "dev"
	OptionDebounceDur     = 800 * time.Millisecond
	OptionDateMillisecond = OptionTimeDuration(time.Millisecond)
)
