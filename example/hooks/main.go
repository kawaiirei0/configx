package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kawaiirei0/configx/v2"
)

// AppConfig 应用配置
type AppConfig struct {
	AppName     string `mapstructure:"app_name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
}

// Logger 简单的日志记录器
type Logger struct {
	prefix string
}

func NewLogger(prefix string) *Logger {
	return &Logger{prefix: prefix}
}

func (l *Logger) log(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("[%s] [%s] [%s] %s\n", timestamp, l.prefix, level, message)
}

func (l *Logger) Debug(message string) {
	l.log("DEBUG", message)
}

func (l *Logger) Info(message string) {
	l.log("INFO", message)
}

func (l *Logger) Warn(message string) {
	l.log("WARN", message)
}

func (l *Logger) Error(message string) {
	l.log("ERROR", message)
}

func main() {
	fmt.Println("=== 钩子示例：日志钩子与事件监听 ===\n")

	// 1. 创建日志记录器
	logger := NewLogger("ConfigManager")

	// 2. 创建配置管理器
	manager := configx.NewManager(AppConfig{})

	// 3. 设置配置选项
	opts := configx.NewOption()
	opts.Filename.Set("config.yaml")
	opts.Filepath.Set("./example/hooks")
	opts.DebounceDur.Set(500 * configx.OptionDateMillisecond)
	manager.SetOption(opts)

	// 4. 设置不同级别的钩子
	fmt.Println("正在设置钩子...")

	// 初始化钩子
	manager.SetHook(configx.InitHook, func(ctx configx.HookContext) {
		logger.Info("🚀 " + ctx.Message)
	})

	// Debug 级别钩子
	manager.SetHook(configx.Debug, func(ctx configx.HookContext) {
		logger.Debug("🔍 " + ctx.Message)
	})

	// Info 级别钩子
	manager.SetHook(configx.Info, func(ctx configx.HookContext) {
		logger.Info("ℹ️  " + ctx.Message)
	})

	// Warn 级别钩子
	manager.SetHook(configx.Warn, func(ctx configx.HookContext) {
		logger.Warn("⚠️  " + ctx.Message)
	})

	// Error 级别钩子
	manager.SetHook(configx.Error, func(ctx configx.HookContext) {
		logger.Error("❌ " + ctx.Message)
	})

	fmt.Println("✓ 钩子设置完成\n")

	// 5. 初始化配置管理器（会触发钩子）
	fmt.Println("正在初始化配置管理器...")
	err := manager.Init(func(ctx *configx.Context) {
		// 配置变更回调
		logger.Info("🔄 配置文件已重新加载")
		
		config, err := manager.GetConfig()
		if err != nil {
			logger.Error(fmt.Sprintf("获取配置失败: %v", err))
			return
		}

		// 显示更新后的配置
		fmt.Println("\n📋 更新后的配置:")
		fmt.Printf("  应用名称: %s\n", config.AppName)
		fmt.Printf("  版本号:   %s\n", config.Version)
		fmt.Printf("  环境:     %s\n", config.Environment)
		fmt.Println()
	})

	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 6. 显示当前配置
	config, err := manager.GetConfig()
	if err != nil {
		log.Fatalf("获取配置失败: %v", err)
	}

	fmt.Println("\n📋 当前配置:")
	fmt.Printf("  应用名称: %s\n", config.AppName)
	fmt.Printf("  版本号:   %s\n", config.Version)
	fmt.Printf("  环境:     %s\n", config.Environment)

	// 7. 演示钩子的作用
	fmt.Println("\n📝 钩子说明:")
	fmt.Println("  ✓ InitHook  - 初始化时触发")
	fmt.Println("  ✓ Debug     - 调试信息")
	fmt.Println("  ✓ Info      - 一般信息（配置加载、变更等）")
	fmt.Println("  ✓ Warn      - 警告信息")
	fmt.Println("  ✓ Error     - 错误信息（加载失败、解析错误等）")

	fmt.Println("\n💡 提示:")
	fmt.Println("  - 修改 example/hooks/config.yaml 文件来触发钩子")
	fmt.Println("  - 观察不同级别钩子的输出")
	fmt.Println("  - 按 Ctrl+C 退出程序")
	fmt.Println("\n等待配置文件变更...")

	// 8. 模拟定期操作，展示钩子在实际应用中的使用
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			config, err := manager.GetConfig()
			if err != nil {
				logger.Error(fmt.Sprintf("定期检查失败: %v", err))
				continue
			}
			logger.Info(fmt.Sprintf("定期检查 - 当前环境: %s", config.Environment))
		}
	}()

	// 9. 等待退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("程序正在退出...")
	fmt.Println()
}
