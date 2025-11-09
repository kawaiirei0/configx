package main

import (
	"fmt"
	"log"

	"github.com/kawaiirei0/configx/v2"
)

// ServerConfig 服务器配置
type ServerConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	ReadTimeout    int    `mapstructure:"read_timeout"`
	WriteTimeout   int    `mapstructure:"write_timeout"`
	MaxConnections int    `mapstructure:"max_connections"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	Host            string `mapstructure:"host"`
	Port            int    `mapstructure:"port"`
	Database        string `mapstructure:"database"`
	Username        string `mapstructure:"username"`
	Password        string `mapstructure:"password"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DB           int    `mapstructure:"db"`
	PoolSize     int    `mapstructure:"pool_size"`
	MinIdleConns int    `mapstructure:"min_idle_conns"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

// AppConfig 应用配置（包含嵌套结构）
type AppConfig struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Logging  LoggingConfig  `mapstructure:"logging"`
}

// Clone 实现自定义克隆方法（实现 Cloneable 接口）
// 这比默认的 JSON 序列化方式更高效
func (c AppConfig) Clone() AppConfig {
	// 由于所有字段都是值类型或简单结构体，直接返回副本即可
	// 如果包含 map、slice 等引用类型，需要进行深拷贝
	return AppConfig{
		Server: ServerConfig{
			Host:           c.Server.Host,
			Port:           c.Server.Port,
			ReadTimeout:    c.Server.ReadTimeout,
			WriteTimeout:   c.Server.WriteTimeout,
			MaxConnections: c.Server.MaxConnections,
		},
		Database: DatabaseConfig{
			Driver:          c.Database.Driver,
			Host:            c.Database.Host,
			Port:            c.Database.Port,
			Database:        c.Database.Database,
			Username:        c.Database.Username,
			Password:        c.Database.Password,
			MaxOpenConns:    c.Database.MaxOpenConns,
			MaxIdleConns:    c.Database.MaxIdleConns,
			ConnMaxLifetime: c.Database.ConnMaxLifetime,
		},
		Redis: RedisConfig{
			Host:         c.Redis.Host,
			Port:         c.Redis.Port,
			Password:     c.Redis.Password,
			DB:           c.Redis.DB,
			PoolSize:     c.Redis.PoolSize,
			MinIdleConns: c.Redis.MinIdleConns,
		},
		Logging: LoggingConfig{
			Level:  c.Logging.Level,
			Format: c.Logging.Format,
			Output: c.Logging.Output,
		},
	}
}

func main() {
	fmt.Println("=== 复杂配置示例：嵌套结构与自定义克隆 ===\n")

	// 1. 创建配置管理器
	manager := configx.NewManager(AppConfig{})

	// 2. 设置配置选项
	opts := configx.NewOption()
	opts.Filename.Set("config.yaml")
	opts.Filepath.Set("./example/complex")
	manager.SetOption(opts)

	// 3. 加载配置
	fmt.Println("正在加载复杂配置...")
	if err := manager.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	fmt.Println("✓ 配置加载成功\n")

	// 4. 获取配置（使用自定义 Clone 方法）
	config, err := manager.GetConfig()
	if err != nil {
		log.Fatalf("获取配置失败: %v", err)
	}

	// 5. 显示配置内容
	fmt.Println("📋 服务器配置:")
	fmt.Printf("  地址:         %s:%d\n", config.Server.Host, config.Server.Port)
	fmt.Printf("  读超时:       %d 秒\n", config.Server.ReadTimeout)
	fmt.Printf("  写超时:       %d 秒\n", config.Server.WriteTimeout)
	fmt.Printf("  最大连接数:   %d\n", config.Server.MaxConnections)

	fmt.Println("\n💾 数据库配置:")
	fmt.Printf("  驱动:         %s\n", config.Database.Driver)
	fmt.Printf("  地址:         %s:%d\n", config.Database.Host, config.Database.Port)
	fmt.Printf("  数据库名:     %s\n", config.Database.Database)
	fmt.Printf("  用户名:       %s\n", config.Database.Username)
	fmt.Printf("  最大连接数:   %d\n", config.Database.MaxOpenConns)
	fmt.Printf("  空闲连接数:   %d\n", config.Database.MaxIdleConns)

	fmt.Println("\n🔴 Redis 配置:")
	fmt.Printf("  地址:         %s:%d\n", config.Redis.Host, config.Redis.Port)
	fmt.Printf("  数据库:       %d\n", config.Redis.DB)
	fmt.Printf("  连接池大小:   %d\n", config.Redis.PoolSize)
	fmt.Printf("  最小空闲连接: %d\n", config.Redis.MinIdleConns)

	fmt.Println("\n📝 日志配置:")
	fmt.Printf("  级别:         %s\n", config.Logging.Level)
	fmt.Printf("  格式:         %s\n", config.Logging.Format)
	fmt.Printf("  输出:         %s\n", config.Logging.Output)

	// 6. 演示自定义克隆方法
	fmt.Println("\n🔧 性能优化:")
	fmt.Println("  ✓ 配置结构实现了 Cloneable 接口")
	fmt.Println("  ✓ GetConfig() 使用自定义 Clone() 方法")
	fmt.Println("  ✓ 避免了 JSON 序列化的性能开销")

	// 7. 验证深拷贝
	fmt.Println("\n🧪 验证深拷贝:")
	config2, _ := manager.GetConfig()
	config2.Server.Port = 9999
	
	config3, _ := manager.GetConfig()
	fmt.Printf("  原始端口:     %d\n", config3.Server.Port)
	fmt.Printf("  修改后端口:   %d\n", config2.Server.Port)
	fmt.Println("  ✓ 配置副本互不影响，深拷贝正常工作")

	fmt.Println("\n示例完成！")
}
