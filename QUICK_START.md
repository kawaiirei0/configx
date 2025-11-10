# ConfigX v2 快速开始

## 安装

```bash
go get github.com/kawaiirei0/configx/v2
```

---

## 基础使用

### 1. 定义配置结构

```go
type AppConfig struct {
    AppName string `mapstructure:"app_name"`
    Version string `mapstructure:"version"`
    Port    int    `mapstructure:"port"`
}
```

### 2. 创建配置管理器

```go
manager := configx.NewManager(AppConfig{})
```

### 3. 设置配置文件

```go
opts := configx.NewOption()
opts.Filename.Set("config.yaml")  // 支持 .yaml, .json, .toml 等
opts.Filepath.Set("./configs")
manager.SetOption(opts)
```

### 4. 加载配置

```go
if err := manager.LoadConfig(); err != nil {
    log.Fatal(err)
}
```

### 5. 获取配置

```go
config, err := manager.GetConfig()
if err != nil {
    log.Fatal(err)
}

fmt.Println(config.AppName)
```

---

## 使用环境变量（推荐）

### 方式 1: 精确绑定（生产环境推荐）

```go
// 1. 创建管理器
manager := configx.NewManager(AppConfig{})

// 2. 设置配置文件
opts := configx.NewOption()
opts.Filename.Set("config.yaml")
manager.SetOption(opts)

// 3. 绑定敏感配置到环境变量
manager.BindEnv("database.password", "DB_PASSWORD")
manager.BindEnv("api.key", "API_KEY")

// 4. 加载配置（环境变量会自动覆盖）
manager.LoadConfig()
```

### 方式 2: 自动环境变量（开发环境）

```go
opts := configx.NewOption()
opts.Filename.Set("config.yaml")
opts.SetEnvPrefix("MYAPP")
opts.EnableAutomaticEnv(true)
opts.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
manager.SetOption(opts)
```

---

## 热重载

```go
manager.Init(func(ctx *configx.Context) {
    log.Println("配置已更新")
    
    // 获取最新配置
    config, _ := manager.GetConfig()
    // 使用新配置...
})
```

---

## 钩子系统

```go
manager.SetHook(configx.Info, func(ctx configx.HookContext) {
    log.Println("Info:", ctx.Message)
})

manager.SetHook(configx.Error, func(ctx configx.HookContext) {
    log.Println("Error:", ctx.Message)
})
```

---

## 完整示例

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/kawaiirei0/configx/v2"
)

type Config struct {
    Database struct {
        Host     string `mapstructure:"host"`
        Port     int    `mapstructure:"port"`
        Password string `mapstructure:"password"`
    } `mapstructure:"database"`
    
    API struct {
        Key string `mapstructure:"key"`
    } `mapstructure:"api"`
}

func main() {
    // 设置环境变量（生产环境中由部署系统设置）
    os.Setenv("DB_PASSWORD", "secure_password_123")
    os.Setenv("API_KEY", "sk_live_xyz")
    
    // 创建管理器
    manager := configx.NewManager(Config{})
    
    // 配置文件
    opts := configx.NewOption()
    opts.Filename.Set("config.yaml")
    manager.SetOption(opts)
    
    // 绑定敏感信息到环境变量
    manager.BindEnv("database.password", "DB_PASSWORD")
    manager.BindEnv("api.key", "API_KEY")
    
    // 加载配置
    if err := manager.LoadConfig(); err != nil {
        log.Fatal(err)
    }
    
    // 获取配置
    config, err := manager.GetConfig()
    if err != nil {
        log.Fatal(err)
    }
    
    // 使用配置
    fmt.Printf("Database: %s:%d\n", 
        config.Database.Host, 
        config.Database.Port)
    fmt.Printf("Password: %s\n", 
        maskSecret(config.Database.Password))
    fmt.Printf("API Key: %s\n", 
        maskSecret(config.API.Key))
}

func maskSecret(s string) string {
    if len(s) <= 8 {
        return "********"
    }
    return s[:4] + "****" + s[len(s)-4:]
}
```

---

## 配置文件示例

### YAML (推荐)

```yaml
database:
  host: localhost
  port: 5432
  password: placeholder  # 将被环境变量覆盖

api:
  key: placeholder  # 将被环境变量覆盖
```

### JSON

```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "password": "placeholder"
  },
  "api": {
    "key": "placeholder"
  }
}
```

### TOML

```toml
[database]
host = "localhost"
port = 5432
password = "placeholder"

[api]
key = "placeholder"
```

---

## 更多示例

查看 `example/` 目录获取更多示例：

- `example/basic/` - 基础使用
- `example/env-bind/` - 环境变量绑定
- `example/multi-format/` - 多格式支持
- `example/hotreload/` - 热重载
- `example/hooks/` - 钩子系统
- `example/complex/` - 复杂配置

---

## 常见问题

### Q: 支持哪些配置文件格式？
A: YAML, JSON, TOML, HCL, INI, Properties

### Q: 如何保护敏感信息？
A: 使用环境变量覆盖，不要将密码写入配置文件

### Q: 是否线程安全？
A: 是的，完全线程安全，可以在多个 goroutine 中使用

### Q: 如何在 Docker 中使用？
A: 通过 `-e` 参数传递环境变量：
```bash
docker run -e DB_PASSWORD=secret123 myapp
```

### Q: 如何在 Kubernetes 中使用？
A: 使用 Secret 和 ConfigMap：
```yaml
env:
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: myapp-secrets
      key: db-password
```

---

## 下一步

- 阅读 [NEW_FEATURES.md](NEW_FEATURES.md) 了解新功能
- 阅读 [CONCURRENCY_SAFETY.md](CONCURRENCY_SAFETY.md) 了解并发安全
- 查看 [example/README.md](example/README.md) 获取更多示例

---

## 获取帮助

- 查看示例代码: `example/`
- 查看文档: `*.md` 文件
- 运行测试: `go test -v ./...`
