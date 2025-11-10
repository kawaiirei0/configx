package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kawaiirei0/configx/v2"
)

// SecureConfig 包含敏感信息的配置
type SecureConfig struct {
	App      AppInfo     `mapstructure:"app"`
	Database DBConfig    `mapstructure:"database"`
	Redis    RedisConfig `mapstructure:"redis"`
	JWT      JWTConfig   `mapstructure:"jwt"`
	AWS      AWSConfig   `mapstructure:"aws"`
}

type AppInfo struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type DBConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
}

type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expiry int    `mapstructure:"expiry"`
}

type AWSConfig struct {
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Region    string `mapstructure:"region"`
}

func main() {
	fmt.Println("=== 特定环境变量绑定示例 ===\n")

	// 1. 设置环境变量（模拟生产环境）
	setupEnvironmentVariables()
	defer cleanupEnvironmentVariables()

	// 2. 创建配置管理器
	manager := configx.NewManager(SecureConfig{})

	// 3. 设置配置选项
	opts := configx.NewOption()
	opts.Filename.Set("config.yaml")
	opts.Filepath.Set("./example/env-bind")
	manager.SetOption(opts)

	// 4. 绑定特定的敏感配置到环境变量
	// 这种方式比 AutomaticEnv 更精确，只绑定需要的配置项
	fmt.Println("绑定敏感配置到环境变量...")

	if err := manager.BindEnv("database.password", "DB_PASSWORD"); err != nil {
		log.Printf("绑定 database.password 失败: %v\n", err)
	}

	if err := manager.BindEnv("redis.password", "REDIS_PASSWORD"); err != nil {
		log.Printf("绑定 redis.password 失败: %v\n", err)
	}

	if err := manager.BindEnv("jwt.secret", "JWT_SECRET"); err != nil {
		log.Printf("绑定 jwt.secret 失败: %v\n", err)
	}

	if err := manager.BindEnv("aws.access_key", "AWS_ACCESS_KEY_ID"); err != nil {
		log.Printf("绑定 aws.access_key 失败: %v\n", err)
	}

	if err := manager.BindEnv("aws.secret_key", "AWS_SECRET_ACCESS_KEY"); err != nil {
		log.Printf("绑定 aws.secret_key 失败: %v\n", err)
	}

	fmt.Println("✓ 环境变量绑定完成\n")

	// 5. 加载配置
	fmt.Println("加载配置文件...")
	if err := manager.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	fmt.Println("✓ 配置加载成功\n")

	// 6. 获取配置
	config, err := manager.GetConfig()
	if err != nil {
		log.Fatalf("获取配置失败: %v", err)
	}

	// 7. 显示配置（敏感信息已被环境变量覆盖）
	displayConfig(config)

	fmt.Println("\n=== 示例完成 ===")
	fmt.Println("\n💡 提示:")
	fmt.Println("  - 敏感信息（密码、密钥）已从环境变量读取")
	fmt.Println("  - 配置文件中的默认值被安全地覆盖")
	fmt.Println("  - 这是生产环境的推荐做法")
}

func setupEnvironmentVariables() {
	fmt.Println("设置环境变量（模拟生产环境）...")
	os.Setenv("DB_PASSWORD", "prod_db_password_secure_123")
	os.Setenv("REDIS_PASSWORD", "prod_redis_password_secure_456")
	os.Setenv("JWT_SECRET", "prod_jwt_secret_very_long_and_secure_789")
	os.Setenv("AWS_ACCESS_KEY_ID", "AKIAIOSFODNN7EXAMPLE")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY")
	fmt.Println("✓ 环境变量设置完成\n")
}

func cleanupEnvironmentVariables() {
	os.Unsetenv("DB_PASSWORD")
	os.Unsetenv("REDIS_PASSWORD")
	os.Unsetenv("JWT_SECRET")
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
}

func displayConfig(config SecureConfig) {
	fmt.Println("📋 配置信息:")
	fmt.Println()

	fmt.Printf("应用信息:\n")
	fmt.Printf("  名称: %s\n", config.App.Name)
	fmt.Printf("  版本: %s\n", config.App.Version)
	fmt.Println()

	fmt.Printf("数据库配置:\n")
	fmt.Printf("  主机: %s\n", config.Database.Host)
	fmt.Printf("  端口: %d\n", config.Database.Port)
	fmt.Printf("  用户名: %s\n", config.Database.Username)
	fmt.Printf("  密码: %s ✓ 从环境变量读取\n", maskSecret(config.Database.Password))
	fmt.Println()

	fmt.Printf("Redis 配置:\n")
	fmt.Printf("  主机: %s\n", config.Redis.Host)
	fmt.Printf("  端口: %d\n", config.Redis.Port)
	fmt.Printf("  密码: %s ✓ 从环境变量读取\n", maskSecret(config.Redis.Password))
	fmt.Println()

	fmt.Printf("JWT 配置:\n")
	fmt.Printf("  密钥: %s ✓ 从环境变量读取\n", maskSecret(config.JWT.Secret))
	fmt.Printf("  过期时间: %d 秒\n", config.JWT.Expiry)
	fmt.Println()

	fmt.Printf("AWS 配置:\n")
	fmt.Printf("  Access Key: %s ✓ 从环境变量读取\n", maskSecret(config.AWS.AccessKey))
	fmt.Printf("  Secret Key: %s ✓ 从环境变量读取\n", maskSecret(config.AWS.SecretKey))
	fmt.Printf("  区域: %s\n", config.AWS.Region)
}

func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "********"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}
