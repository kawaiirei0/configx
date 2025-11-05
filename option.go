package config

import "time"

type Option struct {
	Filename    OptionString
	FileType    OptionString
	Path        OptionString
	Env         OptionString
	DebounceDur OptionTimeDuration
}

// NewOption 创建默认配置
func NewOption() *Option {
	opt := &Option{}
	opt.defaultValueInit()
	return opt
}

// 初始化默认值
func (s *Option) defaultValueInit() *Option {
	s.Filename.Set(OptionFilename)
	s.FileType.Set(OptionFileType)
	s.Path.Set(OptionPath)
	s.Env.Set(OptionEnv)
	s.DebounceDur.Set(OptionTimeDuration(OptionDebounceDur))
	return s
}

type OptionString string
type OptionTimeDuration time.Duration

func (o *OptionString) Set(newStr OptionString) {
	if *o != "" {
		return
	}
	*o = newStr
}

func (o *OptionString) ToValue() string {
	return string(*o)
}

func (o *OptionTimeDuration) Set(newDate OptionTimeDuration) {
	if *o != 0 {
		return
	}
	*o = newDate
}

func (o *OptionTimeDuration) ToValue() time.Duration {
	return time.Duration(*o)
}
