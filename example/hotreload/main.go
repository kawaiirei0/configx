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

// ServerConfig 服务器配置
type ServerConfig struct {
	Host    string `mapstructure:"host"`
	Port    int    `mapstructure:"port"`
	Timeout int    `mapstructure:"timeout"`
	Mode 		string `mapstructure:"mode"`
}

// AppConfig 应用配置
type AppConfig struct {
	Server   ServerConfig `mapstructure:"server"`
	LogLevel string       `mapstructure:"log_level"`
}

func main() {
	fmt.Println("=== 热重载示例：配置文件自动重载 ===\n")

	// 1. 创建配置管理器
	manager := configx.NewManager(AppConfig{})

	// 2. 设置配置选项（包括防抖时间）
	opts := configx.NewOption()
	opts.Filename.Set("config.yaml")
	opts.Filepath.Set("./example/hotreload")
	// 设置防抖时间为 500 毫秒
	opts.DebounceDur.Set(500 * configx.OptionDateMillisecond)
	manager.SetOption(opts)

	// 3. 设置钩子记录配置变更事件
	manager.SetHook(configx.Info, func(ctx configx.HookContext) {
		fmt.Printf("[INFO] %s\n", ctx.Message)
	}).SetHook(configx.Error, func(ctx configx.HookContext) {
		fmt.Printf("[ERROR] %s\n", ctx.Message)
	})

	// 4. 初始化并启动热重载监控
	fmt.Println("正在初始化配置管理器...")
	err := manager.Init(func(ctx *configx.Context) {
		// 配置变更时的回调函数
		fmt.Println("\n🔄 配置文件已更新！")

		// 获取最新配置
		config, err := manager.GetConfig()
		if err != nil {
			fmt.Printf("获取配置失败: %v\n", err)
			return
		}

		// 显示更新后的配置
		fmt.Println("新配置内容:")
		fmt.Printf("  服务器地址: %s:%d\n", config.Server.Host, config.Server.Port)
		fmt.Printf("  超时时间:   %d 秒\n", config.Server.Timeout)
		fmt.Printf("  日志级别:   %s\n", config.LogLevel)
		fmt.Println()
	})

	if err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 5. 显示初始配置
	config, err := manager.GetConfig()
	if err != nil {
		log.Fatalf("获取配置失败: %v", err)
	}

	fmt.Printf("获取到了配置 Config: %v", config)

	fmt.Println("\n当前配置:")
	fmt.Printf("  服务器地址: %s:%d\n", config.Server.Host, config.Server.Port)
	fmt.Printf("  超时时间:   %d 秒\n", config.Server.Timeout)
	fmt.Printf("  日志级别:   %s\n", config.LogLevel)

	// 6. 演示防抖机制
	fmt.Println("\n📝 提示:")
	fmt.Println("  - 修改 example/hotreload/config.yaml 文件来测试热重载")
	fmt.Println("  - 防抖时间设置为 500ms，短时间内的多次修改只会触发一次重载")
	fmt.Println("  - 按 Ctrl+C 退出程序")
	fmt.Println("\n等待配置文件变更...")

	// 7. 模拟应用运行，定期显示当前配置
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			config, err := manager.GetConfig()
			if err != nil {
				continue
			}
			fmt.Printf("\n⏰ [定期检查] 当前端口: %d, 日志级别: %s\n",
				config.Server.Port, config.LogLevel)
		}
	}()

	// 8. 等待退出信号
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("\n\n程序退出")
}
