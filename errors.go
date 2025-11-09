package configx

import "errors"

var (
	// ErrConfigNotInitialized 配置未初始化错误
	ErrConfigNotInitialized = errors.New("配置未初始化")
	
	// ErrConfigFileNotFound 配置文件不存在错误
	ErrConfigFileNotFound = errors.New("配置文件不存在")
	
	// ErrConfigParseFailed 配置解析失败错误
	ErrConfigParseFailed = errors.New("配置解析失败")
	
	// ErrInvalidConfigType 无效的配置类型错误
	ErrInvalidConfigType = errors.New("无效的配置类型")
)
