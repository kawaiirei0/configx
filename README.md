# ConfigX - Go 泛型配置管理器

一个基于 Viper 的轻量级泛型配置管理库，支持自定义配置结构、YAML 文件加载、热更新、防抖处理等功能。

## 功能特性

- 🎯 **泛型设计** - 支持任意自定义配置结构体，类型安全
- 📁 **YAML 支持** - 支持 YAML 格式配置文件
- 🔄 **热更新** - 配置文件变更自动重载
- ⏱️ **防抖机制** - 避免频繁重载，可自定义防抖间隔
- 🔒 **线程安全** - 使用读写锁保证并发访问安全
- 🎣 **钩子系统** - 支持多级别日志钩子（Debug, Info, Warn, Error）
- ⚡ **性能优化** - 支持自定义 Clone 方法优化深拷贝性能

## 快速开始

### 安装

```bash
go get github.com/kawaiirei0/configx
```

### 基本使用

**步骤 1: 定义配置结构体**

```go
package main

import (
    "fmt"
    "log"
    "github.com/kawaiirei0/configx"
)

// 定义你的配置结构体
type AppConfig struct {
    AppName string `mapstructure:"app_name"`
    Version string `mapstructure:"version"`
    Port    int    `mapstructure:"port"`
    Debug   bool   `mapstructure:"debug"`
}
```

**步骤 2: 创建配置管理器**

```go
func main() {
    // 创建泛型配置管理器
    manager := configx.NewManager(AppConfig{})
    
    // 设置配置选项（可选）
    opts := configx.NewOption()
    opts.Filename.Set("config.yaml")
    opts.Filepath.Set("./configs")
    manager.SetOption(opts)
    
    // 加载配置文件
    if err := manager.LoadConfig(); err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }
    
    // 获取配置（类型安全）
    config, err := manager.GetConfig()
    if err != nil {
        log.Fatalf("获取配置失败: %v", err)
    }
    
    // 使用配置
    fmt.Printf("App Name: %s\n", config.AppName)
    fmt.Printf("Version: %s\n", config.Version)
    fmt.Printf("Port: %d\n", config.Port)
}
```

**步骤 3: 创建配置文件**

创建 `configs/config.yaml` 文件：

```yaml
app_name: "MyApplication"
version: "1.0.0"
port: 8080
debug: true
```

## 高级用法

### 自定义配置结构体

ConfigX 支持任意复杂的配置结构，包括嵌套结构：

```go
type ServerConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
    Driver   string `mapstructure:"driver"`
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
}

type AppConfig struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
}

// 使用
manager := configx.NewManager(AppConfig{})
```

### 配置选项

可以通过选项自定义配置行为：

```go
opts := configx.NewOption()
opts.Filename.Set("myconfig.yaml")           // 配置文件名（默认：config.yaml）
opts.Filepath.Set("./config")                // 配置路径（默认：./configs）
opts.DebounceDur.Set(1000 * configx.OptionDateMillisecond)  // 防抖间隔（默认：800ms）

manager := configx.NewManager(AppConfig{})
manager.SetOption(opts)
```

### 热更新

使用 `Init` 方法启动配置监控和热重载：

```go
manager := configx.NewManager(AppConfig{})

// 设置配置选项
opts := configx.NewOption()
opts.Filename.Set("config.yaml")
opts.Filepath.Set("./configs")
manager.SetOption(opts)

// 初始化并启动热重载
err := manager.Init(func(ctx *configx.Context) {
    // 配置变更时的回调函数
    fmt.Println("配置已更新！")
    
    // 获取最新配置
    config, _ := manager.GetConfig()
    fmt.Printf("新端口: %d\n", config.Port)
})

if err != nil {
    log.Fatal(err)
}

// 配置文件变更时会自动重新加载
// 防抖机制会避免频繁重载
```

### 钩子系统

ConfigX 支持多级别的钩子，用于记录配置管理器的各种事件：

```go
manager := configx.NewManager(AppConfig{})

// 设置不同级别的钩子
manager.SetHook(configx.Debug, func(ctx configx.HookContext) {
    fmt.Printf("[DEBUG] %s\n", ctx.Message)
}).SetHook(configx.Info, func(ctx configx.HookContext) {
    fmt.Printf("[INFO] %s\n", ctx.Message)
}).SetHook(configx.Warn, func(ctx configx.HookContext) {
    fmt.Printf("[WARN] %s\n", ctx.Message)
}).SetHook(configx.Error, func(ctx configx.HookContext) {
    fmt.Printf("[ERROR] %s\n", ctx.Message)
})
```

### 性能优化：自定义 Clone 方法

默认情况下，`GetConfig()` 使用 JSON 序列化实现深拷贝。你可以实现 `Cloneable` 接口来提供更高效的克隆方法：

```go
type AppConfig struct {
    Name    string `mapstructure:"name"`
    Version string `mapstructure:"version"`
}

// 实现 Cloneable 接口
func (c AppConfig) Clone() AppConfig {
    return AppConfig{
        Name:    c.Name,
        Version: c.Version,
    }
}

// GetConfig() 会自动使用自定义的 Clone() 方法
manager := configx.NewManager(AppConfig{})
config, _ := manager.GetConfig()  // 使用高效的自定义克隆
```

## API 参考

### 核心类型

```go
// Manager 泛型配置管理器
type Manager[T any] struct { ... }

// 创建配置管理器
func NewManager[T any](defaultConfig T) *Manager[T]
```

### 核心方法

```go
// 加载配置文件
func (m *Manager[T]) LoadConfig() error

// 获取配置副本（类型安全，返回深拷贝）
func (m *Manager[T]) GetConfig() (T, error)

// 初始化并启动热重载监控
func (m *Manager[T]) Init(callback func(*Context)) error

// 设置配置选项（支持链式调用）
func (m *Manager[T]) SetOption(opts *Option) *Manager[T]

// 设置钩子（支持链式调用）
func (m *Manager[T]) SetHook(pattern HookPattern, handler HookHandlerFunc) *Manager[T]
```

### 配置选项

```go
type Option struct {
    Filename    OptionString       // 配置文件名
    Filepath    OptionString       // 配置文件路径
    DebounceDur OptionTimeDuration // 防抖间隔
}

// 创建默认配置选项
func NewOption() *Option
```

### 钩子级别

```go
const (
    InitHook HookPattern = iota  // 初始化钩子
    Debug                        // 调试信息
    Info                         // 一般信息
    Warn                         // 警告信息
    Error                        // 错误信息
)
```

### 错误类型

```go
var (
    ErrConfigNotInitialized error  // 配置未初始化
    ErrConfigFileNotFound   error  // 配置文件不存在
    ErrConfigParseFailed    error  // 配置解析失败
    ErrInvalidConfigType    error  // 无效的配置类型
)
```

### Cloneable 接口

```go
// 实现此接口以提供自定义的高效克隆方法
type Cloneable[T any] interface {
    Clone() T
}
```

## 默认配置

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| 文件名 | config.yaml | 配置文件名 |
| 配置路径 | ./configs | 配置文件目录 |
| 防抖间隔 | 800ms | 文件变更防抖时间 |

## 示例项目

查看 `example/` 目录获取完整示例：

- **basic** - 基础用法示例
- **complex** - 复杂嵌套配置和自定义 Clone 方法
- **hotreload** - 热重载和防抖机制演示
- **hooks** - 钩子系统使用示例

运行示例：

```bash
# 基础示例
go run example/basic/main.go

# 复杂配置示例
go run example/complex/main.go

# 热重载示例
go run example/hotreload/main.go

# 钩子示例
go run example/hooks/main.go
```

## 最佳实践

1. **使用 mapstructure 标签** - 确保配置字段正确映射
2. **实现 Clone 方法** - 对于复杂配置结构，实现自定义 Clone 方法以提升性能
3. **使用钩子记录日志** - 通过钩子系统集成你的日志框架
4. **合理设置防抖时间** - 根据实际需求调整防抖间隔
5. **管理 Manager 实例** - 在应用中创建全局 Manager 实例或使用依赖注入

## 从 v1.x 迁移

如果你正在使用旧版本的 configx，请查看 [MIGRATION.md](MIGRATION.md) 获取详细的迁移指南。

主要变更：
- 不再提供全局单例 `Default()` 函数
- 需要在项目中定义自己的配置结构体
- 使用 `NewManager[T]()` 创建泛型管理器
- `GetConfig()` 现在是 Manager 的方法，返回泛型类型

## 依赖库

- [spf13/viper](https://github.com/spf13/viper) - 配置解析
- [fsnotify/fsnotify](https://github.com/fsnotify/fsnotify) - 文件监控
- [go-viper/mapstructure](https://github.com/go-viper/mapstructure) - 结构体映射