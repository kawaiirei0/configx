# 迁移指南：从 v1.x 到 v2.x

本指南帮助你从 ConfigX v1.x 迁移到 v2.x。v2.x 引入了泛型设计，带来了一些破坏性变更，但也提供了更好的类型安全和灵活性。

## 概述

### 主要变更

v2.x 的核心变更是从硬编码配置模型转变为泛型设计：

- ✅ **泛型支持** - Manager 现在是泛型类型 `Manager[T any]`
- ✅ **类型安全** - 编译时类型检查，无需类型断言
- ✅ **自定义配置** - 在你的项目中定义配置结构体
- ❌ **移除全局单例** - 不再提供 `Default()` 函数
- ❌ **移除硬编码配置** - 删除了库内部的 `Config` 和 `configure.App`

### 为什么要迁移？

- **更好的类型安全**：编译时捕获类型错误
- **更大的灵活性**：自定义任意配置结构
- **更好的性能**：可选的自定义 Clone 方法
- **更清晰的 API**：泛型提供更直观的接口

## 破坏性变更清单

### 1. 全局单例函数已移除

**v1.x:**
```go
manager := config.Default()
```

**v2.x:**
```go
manager := configx.NewManager(AppConfig{})
```

**原因：** Go 泛型不支持泛型全局变量，且鼓励更好的依赖管理模式。

---

### 2. GetConfig 现在是 Manager 的方法

**v1.x:**
```go
cfg, err := config.GetConfig()
```

**v2.x:**
```go
cfg, err := manager.GetConfig()
```

**原因：** 泛型设计要求配置访问通过 Manager 实例。

---

### 3. 配置结构体需要自定义

**v1.x:**
```go
// 使用库内置的 Config 结构
type Config struct {
    App configure.App `mapstructure:"app"`
}
```

**v2.x:**
```go
// 在你的项目中定义配置结构
type AppConfig struct {
    AppName string `mapstructure:"app_name"`
    Version string `mapstructure:"version"`
    Port    int    `mapstructure:"port"`
}
```

**原因：** 泛型设计允许任意配置结构，不再限制于库内部定义。

---

### 4. 配置文件结构可能需要调整

**v1.x (config.yaml):**
```yaml
app:
  name: "MyApp"
  version: "1.0.0"
  description: "My application"
```

**v2.x (config.yaml):**
```yaml
app_name: "MyApp"
version: "1.0.0"
port: 8080
```

**原因：** 配置结构由你定义，YAML 结构需要匹配你的结构体。

---

### 5. 初始化方式变更

**v1.x:**
```go
manager := config.Default()
manager.LoadConfig()
manager.StartMonitor()
```

**v2.x:**
```go
manager := configx.NewManager(AppConfig{})
manager.SetOption(opts)
manager.Init(callback)  // 包含 LoadConfig 和监控启动
```

**原因：** 简化 API，`Init` 方法一次性完成初始化和监控启动。

## 逐步迁移步骤

### 步骤 1: 更新依赖

```bash
# 更新到 v2.x
go get github.com/kawaiirei0/configx@v2
```

### 步骤 2: 定义配置结构体

在你的项目中创建配置结构体文件（例如 `config/types.go`）：

```go
package config

// AppConfig 应用配置
type AppConfig struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
    Host string `mapstructure:"host"`
    Port int    `mapstructure:"port"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
    Driver   string `mapstructure:"driver"`
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Username string `mapstructure:"username"`
    Password string `mapstructure:"password"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
    Host     string `mapstructure:"host"`
    Port     int    `mapstructure:"port"`
    Password string `mapstructure:"password"`
    DB       int    `mapstructure:"db"`
}
```

### 步骤 3: 更新配置文件

根据新的配置结构体更新 YAML 文件：

```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  driver: "mysql"
  host: "localhost"
  port: 3306
  username: "root"
  password: "password"

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0
```

### 步骤 4: 创建配置管理器实例

在你的应用中创建全局 Manager 实例（例如 `config/manager.go`）：

```go
package config

import (
    "github.com/kawaiirei0/configx"
)

var manager *configx.Manager[AppConfig]

// InitConfig 初始化配置管理器
func InitConfig() error {
    manager = configx.NewManager(AppConfig{})
    
    // 设置配置选项
    opts := configx.NewOption()
    opts.Filename.Set("config.yaml")
    opts.Filepath.Set("./configs")
    manager.SetOption(opts)
    
    // 初始化并启动热重载
    return manager.Init(func(ctx *configx.Context) {
        // 配置变更回调
        log.Println("配置已更新")
    })
}

// GetConfig 获取配置
func GetConfig() (AppConfig, error) {
    return manager.GetConfig()
}

// GetManager 获取管理器实例
func GetManager() *configx.Manager[AppConfig] {
    return manager
}
```

### 步骤 5: 更新应用代码

**v1.x:**
```go
package main

import "config"

func main() {
    manager := config.Default()
    manager.LoadConfig()
    
    cfg, err := config.GetConfig()
    if err != nil {
        panic(err)
    }
    
    println(cfg.App.Name)
}
```

**v2.x:**
```go
package main

import "yourproject/config"

func main() {
    if err := config.InitConfig(); err != nil {
        panic(err)
    }
    
    cfg, err := config.GetConfig()
    if err != nil {
        panic(err)
    }
    
    println(cfg.Server.Host)
}
```

### 步骤 6: 更新钩子设置（如果使用）

**v1.x:**
```go
logger := config.NewLogger()
logger.SetHook(func(msg string) {
    log.Println(msg)
})
manager.SetLogger(logger)
```

**v2.x:**
```go
manager.SetHook(configx.Info, func(ctx configx.HookContext) {
    log.Printf("[INFO] %s", ctx.Message)
}).SetHook(configx.Error, func(ctx configx.HookContext) {
    log.Printf("[ERROR] %s", ctx.Message)
})
```

### 步骤 7: 测试和验证

1. 运行应用，确保配置正确加载
2. 修改配置文件，验证热重载功能
3. 检查所有使用配置的地方是否正常工作

## 代码对比示例

### 完整示例对比

#### v1.x 代码

```go
package main

import (
    "fmt"
    "log"
    "config"
)

func main() {
    // 获取默认管理器
    manager := config.Default()
    
    // 设置选项
    option := config.NewOption()
    option.Filename.Set("config")
    option.Path.Set("./configs")
    manager.SetOption(option)
    
    // 加载配置
    if err := manager.LoadConfig(); err != nil {
        log.Fatal(err)
    }
    
    // 启动监控
    if err := manager.StartMonitor(); err != nil {
        log.Fatal(err)
    }
    
    // 获取配置
    cfg, err := config.GetConfig()
    if err != nil {
        log.Fatal(err)
    }
    
    // 使用配置
    fmt.Printf("App: %s\n", cfg.App.Name)
    fmt.Printf("Version: %s\n", cfg.App.Version)
}
```

#### v2.x 代码

```go
package main

import (
    "fmt"
    "log"
    "github.com/kawaiirei0/configx"
)

// 定义配置结构
type AppConfig struct {
    AppName string `mapstructure:"app_name"`
    Version string `mapstructure:"version"`
    Port    int    `mapstructure:"port"`
}

func main() {
    // 创建泛型管理器
    manager := configx.NewManager(AppConfig{})
    
    // 设置选项
    opts := configx.NewOption()
    opts.Filename.Set("config.yaml")
    opts.Filepath.Set("./configs")
    manager.SetOption(opts)
    
    // 初始化（包含加载和监控）
    if err := manager.Init(func(ctx *configx.Context) {
        fmt.Println("配置已更新")
    }); err != nil {
        log.Fatal(err)
    }
    
    // 获取配置（类型安全）
    cfg, err := manager.GetConfig()
    if err != nil {
        log.Fatal(err)
    }
    
    // 使用配置
    fmt.Printf("App: %s\n", cfg.AppName)
    fmt.Printf("Version: %s\n", cfg.Version)
    fmt.Printf("Port: %d\n", cfg.Port)
}
```

## 常见问题解答

### Q1: 为什么移除全局单例？

**A:** Go 泛型不支持泛型全局变量。强行实现会增加复杂度和运行时开销。我们鼓励在应用层实现单例模式或使用依赖注入。

### Q2: 如何在 v2.x 中实现单例模式？

**A:** 在你的项目中创建全局 Manager 实例：

```go
package config

var manager = configx.NewManager(AppConfig{})

func GetManager() *configx.Manager[AppConfig] {
    return manager
}
```

### Q3: 配置结构体必须使用 mapstructure 标签吗？

**A:** 是的。Viper 使用 mapstructure 库进行结构体映射，标签是必需的。

### Q4: 如何优化 GetConfig 的性能？

**A:** 实现 `Cloneable[T]` 接口：

```go
func (c AppConfig) Clone() AppConfig {
    return AppConfig{
        AppName: c.AppName,
        Version: c.Version,
        Port:    c.Port,
    }
}
```

这比默认的 JSON 序列化快 10-100 倍。

### Q5: v1.x 和 v2.x 可以共存吗？

**A:** 可以，但不推荐。如果必须共存，使用不同的导入别名：

```go
import (
    configv1 "github.com/kawaiirei0/configx"
    configv2 "github.com/kawaiirei0/configx/v2"
)
```

### Q6: 迁移需要多长时间？

**A:** 对于小型项目，通常 1-2 小时。对于大型项目，可能需要半天到一天，主要时间花在定义配置结构体和更新配置文件上。

### Q7: 有没有自动化迁移工具？

**A:** 目前没有。由于配置结构的多样性，自动化迁移很困难。但迁移过程相对简单，按照本指南逐步操作即可。

### Q8: 如果遇到问题怎么办？

**A:** 
1. 查看 [示例代码](./example/)
2. 阅读 [API 文档](./API.md)
3. 查看 [故障排除指南](./TROUBLESHOOTING.md)
4. 在 GitHub 提交 Issue

## 迁移检查清单

使用此清单确保迁移完整：

- [ ] 更新依赖到 v2.x
- [ ] 定义自定义配置结构体
- [ ] 更新配置文件格式
- [ ] 创建 Manager 实例管理代码
- [ ] 更新所有使用 `config.Default()` 的地方
- [ ] 更新所有使用 `config.GetConfig()` 的地方
- [ ] 更新钩子设置代码（如果使用）
- [ ] 更新配置选项设置代码
- [ ] 测试配置加载功能
- [ ] 测试热重载功能
- [ ] 测试并发访问场景
- [ ] 更新相关文档和注释
- [ ] 运行完整测试套件

## 获取帮助

如果在迁移过程中遇到问题：

- 📖 查看 [README.md](./README.md) 了解基本用法
- 🏗️ 查看 [ARCHITECTURE.md](./.docs/ARCHITECTURE.md) 了解架构设计
- 📚 查看 [API.md](./API.md) 了解详细 API
- 🔧 查看 [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) 解决常见问题
- 💡 查看 [example/](./example/) 目录的示例代码
- 🐛 在 GitHub 提交 Issue

祝迁移顺利！
