package main

import (
	"fmt"
	"log"

	"github.com/kawaiirei0/configx/v2"
)

// AppConfig 定义简单的应用配置结构体
type AppConfig struct {
	AppName string `mapstructure:"app_name"`
	Version string `mapstructure:"version"`
	Port    int    `mapstructure:"port"`
	Debug   bool   `mapstructure:"debug"`
}

func main() {
	fmt.Println("=== 基础示例：使用泛型配置管理器 ===\n")

	// 1. 创建配置管理器实例
	// 传入配置结构体的零值作为类型参数
	manager := configx.NewManager(AppConfig{})

	// 2. 设置配置选项
	opts := configx.NewOption()
	opts.Filename.Set("config.yaml")
	opts.Filepath.Set("./example/basic")
	manager.SetOption(opts)

	// 3. 加载配置文件
	fmt.Println("正在加载配置文件...")
	if err := manager.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	fmt.Println("✓ 配置加载成功\n")

	// 4. 获取配置（类型安全）
	config, err := manager.GetConfig()
	if err != nil {
		log.Fatalf("获取配置失败: %v", err)
	}

	// 5. 使用配置
	fmt.Println("配置内容:")
	fmt.Printf("  应用名称: %s\n", config.AppName)
	fmt.Printf("  版本号:   %s\n", config.Version)
	fmt.Printf("  端口:     %d\n", config.Port)
	fmt.Printf("  调试模式: %v\n", config.Debug)

	fmt.Println("\n示例完成！")
}
