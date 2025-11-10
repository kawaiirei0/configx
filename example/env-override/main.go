package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kawaiirei0/configx/v2"
)

// AppConfig 应用配置结构
type AppConfig struct {
	App      AppInfo      `mapstructure:"app"`
	Server   ServerConfig `mapstructure:"server"`
	Database DBConfig     `mapstructure:"database"`
	API      APIConfig    `mapstructure:"api"`
}

type AppInfo struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
	Debug   bool   `mapstructure:"debug"`
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type APIConfig struct {
	Key     string `mapstructure:"key"`
	Secret  string `mapstructure:"secret"`
	Timeout int    `mapstructure:"timeout"`
}

func main() {
	fmt.Println("=== 环境变量覆盖示例 ===\n")

	// 演示不同的配置文件格式
	formats := []string{"yaml", "json", "toml"}

	for _, format := range formats {
		fmt.Printf("\n--- 使用 %s 格式 ---\n", strings.ToUpper(format))
		demonstrateFormat(format)
	}

	fmt.Println("\n=== 示例完成 ===")
}

func demonstrateFormat(format string) {
	// 1. 创建配置管理器
	manager := configx.NewManager(AppConfig{})

	// 2. 设置配置选项
	opts := configx.NewOption()
	opts.Filename.Set(configx.OptionString(fmt.Sprintf("config.%s", format)))
	opts.Filepath.Set("./example/env-override")

	// 3. 配置环境变量支持
	opts.SetEnvPrefix("MYAPP")                            // 环境变量前缀
	opts.EnableAutomaticEnv(true)                         // 启用自动环境变量
	opts.SetAllowEmptyEnv(false)                          // 不允许空环境变量
	opts.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // 将 . 替换为 _

	manager.SetOption(opts)

	// 4. 设置一些环境变量（模拟实际使用）
	// 注意：环境变量名格式为 PREFIX_SECTION_KEY
	os.Setenv("MYAPP_DATABASE_PASSWORD", "env_password_123")
	os.Setenv("MYAPP_API_KEY", "env_api_key_xyz")
	os.Setenv("MYAPP_API_SECRET", "env_secret_abc")
	os.Setenv("MYAPP_SERVER_PORT", "9090")
	os.Setenv("MYAPP_APP_DEBUG", "true")

	defer func() {
		os.Unsetenv("MYAPP_DATABASE_PASSWORD")
		os.Unsetenv("MYAPP_API_KEY")
		os.Unsetenv("MYAPP_API_SECRET")
		os.Unsetenv("MYAPP_SERVER_PORT")
		os.Unsetenv("MYAPP_APP_DEBUG")
	}()

	// 5. 加载配置
	if err := manager.LoadConfig(); err != nil {
		log.Printf("加载配置失败: %v\n", err)
		return
	}

	// 6. 获取配置
	config, err := manager.GetConfig()
	if err != nil {
		log.Printf("获取配置失败: %v\n", err)
		return
	}

	// 7. 显示配置（注意敏感信息被环境变量覆盖）
	fmt.Printf("应用名称: %s\n", config.App.Name)
	fmt.Printf("应用版本: %s\n", config.App.Version)
	fmt.Printf("服务器端口: %d (环境变量覆盖: %s)\n",
		config.Server.Port,
		getEnvStatus(config.Server.Port != 8080))

	fmt.Printf("\n数据库配置:\n")
	fmt.Printf("  主机: %s\n", config.Database.Host)
	fmt.Printf("  端口: %d\n", config.Database.Port)
	fmt.Printf("  用户名: %s\n", config.Database.Username)
	fmt.Printf("  密码: %s (环境变量覆盖: %s)\n",
		maskPassword(config.Database.Password),
		getEnvStatus(config.Database.Password != "default_password"))

	fmt.Printf("\nAPI 配置:\n")
	fmt.Printf("  API Key: %s (环境变量覆盖: %s)\n",
		maskSecret(config.API.Key),
		getEnvStatus(config.API.Key != "default_api_key"))
	fmt.Printf("  API Secret: %s (环境变量覆盖: %s)\n",
		maskSecret(config.API.Secret),
		getEnvStatus(config.API.Secret != "default_secret"))
	fmt.Printf("  超时时间: %d 秒\n", config.API.Timeout)
}

func maskPassword(password string) string {
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}

func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "********"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}

func getEnvStatus(isOverridden bool) string {
	if isOverridden {
		return "✓ 已覆盖"
	}
	return "✗ 使用默认值"
}
