package configx

import (
	"time"

	"github.com/kawaiirei0/configx/v2/utils"
)

type Option struct {
	value          string
	pathValue      string
	fileValue      string
	Filename       OptionString
	Filepath       OptionString
	DebounceDur    OptionTimeDuration
	EnvPrefix      OptionString // 环境变量前缀，例如 "APP"
	AutomaticEnv   bool         // 是否自动从环境变量读取
	AllowEmptyEnv  bool         // 是否允许空的环境变量值
	EnvKeyReplacer interface{}  // 环境变量键名转换（*strings.Replacer）
}

// NewOption 创建默认配置
func NewOption() *Option {
	opt := &Option{}
	opt.setDefaultValue()
	return opt
}

// 初始化默认值
func (s *Option) setDefaultValue() *Option {
	s.Filename.Set(OptionFilename, false)
	s.Filepath.Set(OptionFilepath, false)
	s.DebounceDur.Set(OptionTimeDuration(OptionDebounceDur), false)
	s.AutomaticEnv = false  // 默认不启用自动环境变量
	s.AllowEmptyEnv = false // 默认不允许空环境变量
	return s
}

type OptionString string
type OptionTimeDuration time.Duration

func (o *OptionString) Set(newStr OptionString, reset ...bool) {
	if len(reset) == 0 {
		reset = []bool{true}
	}
	if *o != "" && !reset[0] {
		return
	}
	*o = newStr
}

func (o *OptionString) ToValue() string {
	return string(*o)
}

func (o *OptionTimeDuration) Set(newDate OptionTimeDuration, reset ...bool) {
	if len(reset) == 0 {
		reset = []bool{true}
	}
	if *o != 0 && !reset[0] {
		return
	}
	*o = newDate
}

func (o *OptionTimeDuration) ToValue() time.Duration {
	return time.Duration(*o)
}

func (s *Option) File() string {
	if s.fileValue != "" {
		return s.fileValue
	}
	s.fileValue = utils.ConfigPath(s.Filepath.ToValue(), s.Filename.ToValue(), true)
	return s.File()
}

func (s *Option) Path() string {
	if s.pathValue != "" {
		return s.pathValue
	}
	s.pathValue = utils.ConfigPath(s.Filepath.ToValue(), "", true)
	return s.Path()
}

// SetEnvPrefix 设置环境变量前缀
// 例如：SetEnvPrefix("APP") 会使配置项 "database.host" 对应环境变量 "APP_DATABASE_HOST"
func (s *Option) SetEnvPrefix(prefix string) *Option {
	s.EnvPrefix.Set(OptionString(prefix))
	return s
}

// EnableAutomaticEnv 启用自动环境变量读取
// 启用后，配置项会自动从环境变量中读取（如果存在）
func (s *Option) EnableAutomaticEnv(enable bool) *Option {
	s.AutomaticEnv = enable
	return s
}

// SetAllowEmptyEnv 设置是否允许空的环境变量值
func (s *Option) SetAllowEmptyEnv(allow bool) *Option {
	s.AllowEmptyEnv = allow
	return s
}

// SetEnvKeyReplacer 设置环境变量键名转换
// 参数：replacer - *strings.Replacer 实例
// 例如：strings.NewReplacer(".", "_") 将配置键中的 "." 替换为 "_"
func (s *Option) SetEnvKeyReplacer(replacer interface{}) *Option {
	s.EnvKeyReplacer = replacer
	return s
}
