package main

import (
	"fmt"
	"log"

	"github.com/kawaiirei0/configx/v2"
)

// AppConfig 简单的应用配置
type AppConfig struct {
	AppName string `mapstructure:"app_name"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
	Debug   bool   `mapstructure:"debug"`
}

func main() {
	fmt.Println("=== 多格式配置文件示例 ===\n")

	formats := []string{"yaml", "json", "toml"}

	for _, format := range formats {
		fmt.Printf("--- 测试 %s 格式 ---\n", format)
		testFormat(format)
		fmt.Println()
	}

	fmt.Println("=== 示例完成 ===")
	fmt.Println("\n✓ ConfigX 支持多种配置文件格式")
	fmt.Println("✓ 自动根据文件扩展名识别格式")
	fmt.Println("✓ 无需额外配置")
}

func testFormat(format string) {
	// 创建配置管理器
	manager := configx.NewManager(AppConfig{})

	// 设置配置选项
	opts := configx.NewOption()
	opts.Filename.Set(configx.OptionString(fmt.Sprintf("config.%s", format)))
	opts.Filepath.Set("./example/multi-format")
	manager.SetOption(opts)

	// 加载配置
	if err := manager.LoadConfig(); err != nil {
		log.Printf("加载 %s 配置失败: %v\n", format, err)
		return
	}

	// 获取配置
	config, err := manager.GetConfig()
	if err != nil {
		log.Printf("获取配置失败: %v\n", err)
		return
	}

	// 显示配置
	fmt.Printf("  应用名称: %s\n", config.AppName)
	fmt.Printf("  版本号:   %s\n", config.Version)
	fmt.Printf("  端口:     %d\n", config.Port)
	fmt.Printf("  调试模式: %v\n", config.Debug)
	fmt.Printf("  ✓ %s 格式加载成功\n", format)
}
