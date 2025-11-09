# API 文档

本文档详细说明 ConfigX v2.x 的所有公开 API。

## 目录

- [核心类型](#核心类型)
- [构造函数](#构造函数)
- [Manager 方法](#manager-方法)
- [配置选项](#配置选项)
- [钩子系统](#钩子系统)
- [错误类型](#错误类型)
- [接口](#接口)
- [常量](#常量)

---

## 核心类型

### Manager[T any]

泛型配置管理器，负责配置的加载、访问、监控和热重载。

```go
type Manager[T any] struct {
    // 私有字段，不直接访问
}
```

**类型参数：**
- `T any` - 配置结构体类型，可以是任意类型

**示例：**
```go
type AppConfig struct {
    Port int `mapstructure:"port"`
}

manager := configx.NewManager(AppConfig{})
```

---

## 构造函数

### NewManager[T any]

创建新的泛型配置管理器实例。

```go
func NewManager[T any](defaultConfig T) *Manager[T]
```

**参数：**
- `defaultConfig T` - 配置结构体的零值或默认值，用于类型推断

**返回值：**
- `*Manager[T]` - 配置管理器实例

**示例：**
```go
// 使用零值
manager := configx.NewManager(AppConfig{})

// 使用默认值
manager := configx.NewManager(AppConfig{
    Port: 8080,
    Host: "localhost",
})
```

**注意：**
- `defaultConfig` 参数仅用于类型推断，不会作为实际配置使用
- 实际配置通过 `LoadConfig()` 从文件加载

---

## Manager 方法

### LoadConfig

加载配置文件并解析到泛型类型。

```go
func (m *Manager[T]) LoadConfig() error
```

**返回值：**
- `error` - 如果加载或解析失败则返回错误

**错误类型：**
- `ErrConfigFileNotFound` - 配置文件不存在
- `ErrConfigParseFailed` - 配置解析失败

**示例：**
```go
manager := configx.NewManager(AppConfig{})

opts := configx.NewOption()
opts.Filename.Set("config.yaml")
opts.Filepath.Set("./configs")
manager.SetOption(opts)

if err := manager.LoadConfig(); err != nil {
    log.Fatalf("加载配置失败: %v", err)
}
```

**行为：**
1. 使用写锁保护配置更新
2. 配置 Viper 实例
3. 读取配置文件
4. 解析 YAML 到泛型类型 T
5. 更新内部配置指针

---

### GetConfig

获取配置的深拷贝副本，保证线程安全。

```go
func (m *Manager[T]) GetConfig() (T, error)
```

**返回值：**
- `T` - 配置副本
- `error` - 如果配置未初始化则返回错误

**错误类型：**
- `ErrConfigNotInitialized` - 配置未初始化（未调用 LoadConfig 或 Init）

**示例：**
```go
config, err := manager.GetConfig()
if err != nil {
    log.Fatalf("获取配置失败: %v", err)
}

fmt.Printf("Port: %d\n", config.Port)
```

**行为：**
1. 使用读锁保护并发访问
2. 检查配置是否已初始化
3. 如果配置类型实现了 `Cloneable[T]` 接口，使用自定义 Clone 方法
4. 否则使用 JSON 序列化/反序列化实现深拷贝
5. 返回配置副本

**性能优化：**
```go
// 实现 Cloneable 接口以提升性能
func (c AppConfig) Clone() AppConfig {
    return AppConfig{
        Port: c.Port,
        Host: c.Host,
    }
}
```

---

### Init

初始化配置管理器，加载配置并启动热重载监控。

```go
func (m *Manager[T]) Init(callback func(*Context)) error
```

**参数：**
- `callback func(*Context)` - 配置变更时的回调函数，可以为 nil

**返回值：**
- `error` - 如果初始化失败则返回错误

**示例：**
```go
err := manager.Init(func(ctx *configx.Context) {
    fmt.Println("配置已更新！")
    
    // 获取最新配置
    config, _ := manager.GetConfig()
    fmt.Printf("新端口: %d\n", config.Port)
})

if err != nil {
    log.Fatal(err)
}
```

**行为：**
1. 初始化配置选项（如果未初始化）
2. 调用 `LoadConfig()` 加载初始配置
3. 触发 `InitHook` 钩子
4. 启动文件监控（`monitorConfigChanges`）
5. 配置文件变更时自动重载并执行回调

**注意：**
- `Init` 是启动热重载的推荐方式
- 回调函数在配置成功重载后执行
- 如果重载失败，保持原有配置不变

---

### SetOption

设置配置选项，支持链式调用。

```go
func (m *Manager[T]) SetOption(opts *Option) *Manager[T]
```

**参数：**
- `opts *Option` - 配置选项，如果为 nil 则使用默认选项

**返回值：**
- `*Manager[T]` - 返回管理器实例以支持链式调用

**示例：**
```go
opts := configx.NewOption()
opts.Filename.Set("config.yaml")
opts.Filepath.Set("./configs")
opts.DebounceDur.Set(1000 * configx.OptionDateMillisecond)

manager := configx.NewManager(AppConfig{})
manager.SetOption(opts).SetHook(configx.Info, func(ctx configx.HookContext) {
    log.Println(ctx.Message)
})
```

**注意：**
- 必须在 `LoadConfig()` 或 `Init()` 之前调用
- 如果传入 nil，将使用默认选项

---

### SetHook

设置钩子处理函数，支持链式调用。

```go
func (m *Manager[T]) SetHook(pattern HookPattern, handler HookHandlerFunc) *Manager[T]
```

**参数：**
- `pattern HookPattern` - 钩子级别（InitHook, Debug, Info, Warn, Error）
- `handler HookHandlerFunc` - 钩子处理函数

**返回值：**
- `*Manager[T]` - 返回管理器实例以支持链式调用

**示例：**
```go
manager.SetHook(configx.Info, func(ctx configx.HookContext) {
    log.Printf("[INFO] %s", ctx.Message)
}).SetHook(configx.Error, func(ctx configx.HookContext) {
    log.Printf("[ERROR] %s", ctx.Message)
}).SetHook(configx.Debug, func(ctx configx.HookContext) {
    log.Printf("[DEBUG] %s", ctx.Message)
})
```

**钩子触发时机：**
- `InitHook` - 初始化完成时
- `Debug` - 调试信息
- `Info` - 配置加载、重载成功时
- `Warn` - 警告信息
- `Error` - 配置加载、重载失败时

---

## 配置选项

### Option

配置选项结构体。

```go
type Option struct {
    Filename    OptionString       // 配置文件名
    Filepath    OptionString       // 配置文件路径
    DebounceDur OptionTimeDuration // 防抖间隔
}
```

### NewOption

创建默认配置选项。

```go
func NewOption() *Option
```

**返回值：**
- `*Option` - 配置选项实例

**默认值：**
- `Filename`: "config.yaml"
- `Filepath`: "./configs"
- `DebounceDur`: 800ms

**示例：**
```go
opts := configx.NewOption()
```

---

### OptionString

字符串类型的配置选项。

```go
type OptionString string
```

#### Set

设置字符串选项值。

```go
func (o *OptionString) Set(newStr OptionString, reset ...bool)
```

**参数：**
- `newStr OptionString` - 新值
- `reset ...bool` - 可选，是否强制重置（默认 true）

**示例：**
```go
opts.Filename.Set("myconfig.yaml")
opts.Filepath.Set("./config")
```

#### ToValue

获取字符串选项的值。

```go
func (o *OptionString) ToValue() string
```

**返回值：**
- `string` - 选项值

---

### OptionTimeDuration

时间间隔类型的配置选项。

```go
type OptionTimeDuration time.Duration
```

#### Set

设置时间间隔选项值。

```go
func (o *OptionTimeDuration) Set(newDate OptionTimeDuration, reset ...bool)
```

**参数：**
- `newDate OptionTimeDuration` - 新值
- `reset ...bool` - 可选，是否强制重置（默认 true）

**示例：**
```go
// 设置防抖时间为 1 秒
opts.DebounceDur.Set(1000 * configx.OptionDateMillisecond)

// 设置防抖时间为 500 毫秒
opts.DebounceDur.Set(500 * configx.OptionDateMillisecond)
```

#### ToValue

获取时间间隔选项的值。

```go
func (o *OptionTimeDuration) ToValue() time.Duration
```

**返回值：**
- `time.Duration` - 选项值

---

## 钩子系统

### HookPattern

钩子级别枚举。

```go
type HookPattern int

const (
    InitHook HookPattern = iota  // 初始化钩子
    Debug                        // 调试信息
    Info                         // 一般信息
    Warn                         // 警告信息
    Error                        // 错误信息
)
```

### HookHandlerFunc

钩子处理函数类型。

```go
type HookHandlerFunc func(HookContext)
```

### HookContext

钩子上下文，包含钩子触发时的信息。

```go
type HookContext struct {
    Message string  // 消息内容
}
```

**示例：**
```go
manager.SetHook(configx.Info, func(ctx configx.HookContext) {
    timestamp := time.Now().Format("2006-01-02 15:04:05")
    fmt.Printf("[%s] %s\n", timestamp, ctx.Message)
})
```

---

### Context

配置变更回调的上下文。

```go
type Context struct {
    // 包含管理器引用等信息
}
```

**用途：**
- 在 `Init` 方法的回调函数中使用
- 提供配置变更时的上下文信息

**示例：**
```go
manager.Init(func(ctx *configx.Context) {
    // 使用上下文信息
    config, _ := manager.GetConfig()
    log.Printf("配置已更新: %+v", config)
})
```

---

## 错误类型

### ErrConfigNotInitialized

配置未初始化错误。

```go
var ErrConfigNotInitialized = errors.New("配置未初始化")
```

**触发条件：**
- 在调用 `LoadConfig()` 或 `Init()` 之前调用 `GetConfig()`

**处理方式：**
```go
config, err := manager.GetConfig()
if errors.Is(err, configx.ErrConfigNotInitialized) {
    log.Println("配置未初始化，正在加载...")
    manager.LoadConfig()
}
```

---

### ErrConfigFileNotFound

配置文件不存在错误。

```go
var ErrConfigFileNotFound = errors.New("配置文件不存在")
```

**触发条件：**
- 配置文件路径不存在
- 配置文件名错误

**处理方式：**
```go
if err := manager.LoadConfig(); err != nil {
    if errors.Is(err, configx.ErrConfigFileNotFound) {
        log.Println("配置文件不存在，请检查路径")
    }
}
```

---

### ErrConfigParseFailed

配置解析失败错误。

```go
var ErrConfigParseFailed = errors.New("配置解析失败")
```

**触发条件：**
- YAML 格式错误
- 配置结构与 YAML 不匹配
- 类型转换失败

**处理方式：**
```go
if err := manager.LoadConfig(); err != nil {
    if errors.Is(err, configx.ErrConfigParseFailed) {
        log.Println("配置解析失败，请检查 YAML 格式")
    }
}
```

---

### ErrInvalidConfigType

无效的配置类型错误。

```go
var ErrInvalidConfigType = errors.New("无效的配置类型")
```

**触发条件：**
- 配置类型不符合预期

---

## 接口

### Cloneable[T any]

可克隆接口，用于提供自定义的高效克隆方法。

```go
type Cloneable[T any] interface {
    Clone() T
}
```

**用途：**
- 优化 `GetConfig()` 的性能
- 避免 JSON 序列化的开销

**实现示例：**
```go
type AppConfig struct {
    Port int    `mapstructure:"port"`
    Host string `mapstructure:"host"`
}

// 实现 Cloneable 接口
func (c AppConfig) Clone() AppConfig {
    return AppConfig{
        Port: c.Port,
        Host: c.Host,
    }
}
```

**性能对比：**
- JSON 序列化：~1000 ns/op
- 自定义 Clone：~10 ns/op
- 性能提升：100 倍

**注意：**
- 对于包含 map、slice 等引用类型的结构体，需要进行深拷贝
- 对于简单结构体，直接返回副本即可

**复杂示例：**
```go
type ComplexConfig struct {
    Name    string            `mapstructure:"name"`
    Tags    []string          `mapstructure:"tags"`
    Metadata map[string]string `mapstructure:"metadata"`
}

func (c ComplexConfig) Clone() ComplexConfig {
    // 深拷贝 slice
    tags := make([]string, len(c.Tags))
    copy(tags, c.Tags)
    
    // 深拷贝 map
    metadata := make(map[string]string, len(c.Metadata))
    for k, v := range c.Metadata {
        metadata[k] = v
    }
    
    return ComplexConfig{
        Name:     c.Name,
        Tags:     tags,
        Metadata: metadata,
    }
}
```

---

## 常量

### OptionDateMillisecond

毫秒时间单位常量。

```go
const OptionDateMillisecond = OptionTimeDuration(time.Millisecond)
```

**用途：**
- 设置防抖间隔时使用

**示例：**
```go
// 设置防抖时间为 500 毫秒
opts.DebounceDur.Set(500 * configx.OptionDateMillisecond)

// 设置防抖时间为 1 秒
opts.DebounceDur.Set(1000 * configx.OptionDateMillisecond)
```

---

## 完整使用示例

### 基础示例

```go
package main

import (
    "fmt"
    "log"
    "github.com/kawaiirei0/configx"
)

type AppConfig struct {
    AppName string `mapstructure:"app_name"`
    Version string `mapstructure:"version"`
    Port    int    `mapstructure:"port"`
}

func main() {
    // 1. 创建管理器
    manager := configx.NewManager(AppConfig{})
    
    // 2. 设置选项
    opts := configx.NewOption()
    opts.Filename.Set("config.yaml")
    opts.Filepath.Set("./configs")
    manager.SetOption(opts)
    
    // 3. 加载配置
    if err := manager.LoadConfig(); err != nil {
        log.Fatal(err)
    }
    
    // 4. 获取配置
    config, err := manager.GetConfig()
    if err != nil {
        log.Fatal(err)
    }
    
    // 5. 使用配置
    fmt.Printf("App: %s\n", config.AppName)
    fmt.Printf("Version: %s\n", config.Version)
    fmt.Printf("Port: %d\n", config.Port)
}
```

### 热重载示例

```go
package main

import (
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"
    "github.com/kawaiirei0/configx"
)

type AppConfig struct {
    Port int `mapstructure:"port"`
}

func main() {
    manager := configx.NewManager(AppConfig{})
    
    opts := configx.NewOption()
    opts.Filename.Set("config.yaml")
    opts.Filepath.Set("./configs")
    opts.DebounceDur.Set(500 * configx.OptionDateMillisecond)
    manager.SetOption(opts)
    
    // 设置钩子
    manager.SetHook(configx.Info, func(ctx configx.HookContext) {
        fmt.Printf("[INFO] %s\n", ctx.Message)
    })
    
    // 初始化并启动热重载
    err := manager.Init(func(ctx *configx.Context) {
        config, _ := manager.GetConfig()
        fmt.Printf("配置已更新，新端口: %d\n", config.Port)
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("等待配置变更...")
    
    // 等待退出信号
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    <-sigChan
}
```

### 自定义 Clone 示例

```go
package main

import (
    "fmt"
    "github.com/kawaiirei0/configx"
)

type AppConfig struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
}

type ServerConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}

type DatabaseConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}

// 实现 Cloneable 接口
func (c AppConfig) Clone() AppConfig {
    return AppConfig{
        Server: ServerConfig{
            Host: c.Server.Host,
            Port: c.Server.Port,
        },
        Database: DatabaseConfig{
            Host: c.Database.Host,
            Port: c.Database.Port,
        },
    }
}

func main() {
    manager := configx.NewManager(AppConfig{})
    
    opts := configx.NewOption()
    opts.Filename.Set("config.yaml")
    manager.SetOption(opts)
    
    manager.LoadConfig()
    
    // GetConfig 会自动使用自定义的 Clone 方法
    config, _ := manager.GetConfig()
    fmt.Printf("Server: %s:%d\n", config.Server.Host, config.Server.Port)
}
```

---

## 类型安全说明

ConfigX v2.x 的泛型设计提供了完整的类型安全：

1. **编译时类型检查**：
   ```go
   manager := configx.NewManager(AppConfig{})
   config, _ := manager.GetConfig()
   // config 的类型是 AppConfig，不需要类型断言
   fmt.Println(config.Port)  // 编译器知道 Port 字段存在
   ```

2. **无需类型断言**：
   ```go
   // v1.x 需要类型断言
   cfg, _ := config.GetConfig()
   appCfg := cfg.(*Config)  // 运行时类型断言
   
   // v2.x 不需要
   cfg, _ := manager.GetConfig()
   // cfg 已经是正确的类型
   ```

3. **IDE 支持**：
   - 自动补全
   - 类型提示
   - 重构支持

---

## 最佳实践

1. **使用 Init 方法**：
   ```go
   // 推荐：一次性完成初始化和监控
   manager.Init(callback)
   
   // 不推荐：分开调用
   manager.LoadConfig()
   manager.StartMonitor()
   ```

2. **实现 Clone 方法**：
   ```go
   // 对于复杂配置，实现自定义 Clone 方法
   func (c AppConfig) Clone() AppConfig {
       // 自定义克隆逻辑
   }
   ```

3. **使用钩子记录日志**：
   ```go
   manager.SetHook(configx.Info, func(ctx configx.HookContext) {
       logger.Info(ctx.Message)
   })
   ```

4. **合理设置防抖时间**：
   ```go
   // 开发环境：快速响应
   opts.DebounceDur.Set(200 * configx.OptionDateMillisecond)
   
   // 生产环境：避免频繁重载
   opts.DebounceDur.Set(1000 * configx.OptionDateMillisecond)
   ```

5. **错误处理**：
   ```go
   config, err := manager.GetConfig()
   if err != nil {
       if errors.Is(err, configx.ErrConfigNotInitialized) {
           // 处理未初始化错误
       }
       return err
   }
   ```

---

## 参考资料

- [README.md](./README.md) - 快速开始指南
- [ARCHITECTURE.md](./.docs/ARCHITECTURE.md) - 架构设计文档
- [MIGRATION.md](./MIGRATION.md) - 迁移指南
- [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) - 故障排除指南
- [示例代码](./example/) - 完整示例
