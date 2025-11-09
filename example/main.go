package main

import (
	"fmt"
	"github.com/kawaiirei0/configx"

	"github.com/kawaiirei0/configx/utils"
)

type Logger struct {
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) log(level, format string, field ...string) {
	fmt.Printf("[LOG] [%s] %s\n", level, format)
}

func (l *Logger) Debug(format string, field ...string) {
	l.log("DEBUG", format)
}

func (l *Logger) Info(format string, field ...string) {
	l.log("INFO", format)
}

func (l *Logger) Warn(format string, field ...string) {
	l.log("WARN", format)
}

func (l *Logger) Error(format string, field ...string) {
	l.log("ERROR", format)
}

func main() {
	// 实例化日志
	log := NewLogger()

	// setting config hook
	configx.SetHook(configx.Info, func(ctx configx.HookContext) {
		log.Info("正在初始化配置...")
	}).SetHook(configx.Debug, func(ctx configx.HookContext) {
		log.Debug(ctx.Message)
	}).SetHook(configx.Info, func(ctx configx.HookContext) {
		log.Info(ctx.Message)
	}).SetHook(configx.Warn, func(ctx configx.HookContext) {
		log.Warn(ctx.Message)
	}).SetHook(configx.Error, func(ctx configx.HookContext) {
		log.Error(ctx.Message)
	})

	// 设置配置选项
	opts := configx.NewOption()
	opts.Filename.Set("config.dev.json")                      // production | development
	opts.Filepath.Set("./configs")                            // 设置文件夹
	opts.DebounceDur.Set(800 * configx.OptionDateMillisecond) // 设置防抖

	path := utils.ConfigPath(".qwq", "config.toml", true)
	fmt.Println("配置文件路径:", path)

	// 实例化配置管理器
	manager := configx.Default()

	manager.SetOption(opts)

	if err := manager.Init(func(ctx *configx.Context) {
		log.Debug("配置文件更新了")
	}); err != nil {
		log.Error(fmt.Sprintf("初始化配置失败 Error: %s", err.Error()))
	}

	log.Info("初始化成功")

	select {}
}
