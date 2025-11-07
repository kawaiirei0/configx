# Go 配置管理器

一个基于 Viper 的轻量级配置管理库，支持 YAML 文件加载、热更新、防抖处理等功能。

## 功能特性

- 支持 YAML 格式配置文件
- 配置热更新（文件变更监听）
- 防抖机制，避免频繁重载
- 线程安全的配置访问
- 结构化配置定义
- 日志钩子支持
- 单例模式管理器

## 快速开始

### 安装依赖

```bash
go mod tidy
```

### 基本使用

具体请查看示例代码

```go
package main

import (
    "config"
    "fmt"
)

func main() {
    // 初始化配置管理器
    manager := config.Default()
    
    // 初始化加载配置文件
    err := manager.Init()
    if err != nil {
        panic(err)
    }
    
    // 获取配置
    cfg, err := config.GetConfig()
    if err != nil {
        panic(err)
    }
    
    // 使用配置
    fmt.Printf("App Name: %s\n", cfg.App.Name)
    fmt.Printf("Version: %s\n", cfg.App.Version)
}
```

### 配置文件格式

创建 `configs/config.yaml` 文件：

```yaml
app:
  name: "MyApplication"
  version: "1.0.0"
  description: "A sample application"
```

### 配置选项

可以通过选项自定义配置行为：

```go
option := config.NewOption()
option.Filename.Set("myconfig.yaml")     // 配置文件名（默认：config.dev.yaml）
option.Filepath.Set("./config")         // 配置路径（默认：./configs）
option.DebounceDur.Set(1000)        // 防抖间隔，毫秒（默认：800ms）

manager := config.Default()
manager.SetOption(option)
```

### 热更新

配置管理器支持文件变更监听和热更新：

```go
manager := config.Default()
manager.Init()

// 初始化会默认开启监听
// 配置文件变更时会自动重新加载
// 防抖机制会避免频繁重载
```

### 日志钩子

可以设置日志钩子记录配置变更：

```go
logger := config.NewLogger()
logger.SetHook(func(msg string) {
    // 自定义日志处理
    fmt.Printf("[Config] %s\n", msg)
})

manager := config.Default()
manager.SetLogger(logger)
```

## API 参考

### 核心函数

- `config.Default() *Manager` - 获取默认配置管理器实例
- `config.GetConfig() (*Config, error)` - 获取当前配置副本
- `manager.LoadConfig() error` - 加载配置文件
- `manager.StartMonitor() error` - 启动文件监控
- `manager.StopMonitor() error` - 停止文件监控

### 配置结构

```go
type Config struct {
    App configure.App `mapstructure:"app"`
}

type App struct {
    Name        string `yaml:"name" mapstructure:"name"`
    Version     string `yaml:"version" mapstructure:"version"`
    Description string `yaml:"description" mapstructure:"description"`
}
```

## 默认配置

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| 文件名 | config | 配置文件名 |
| 文件类型 | yaml | 配置文件格式 |
| 配置路径 | ./configs | 配置文件目录 |
| 环境标识 | dev | 运行环境 |
| 防抖间隔 | 800ms | 文件变更防抖时间 |

## 注意事项

1. 配置文件必须放在指定的配置目录下
2. 配置访问是线程安全的，但配置更新时需要重新获取
3. 热更新功能需要显式启动监控
4. 建议使用结构体标签确保配置正确映射

## 依赖库

- [spf13/viper](https://github.com/spf13/viper) - 配置解析
- [fsnotify/fsnotify](https://github.com/fsnotify/fsnotify) - 文件监控
- [go-viper/mapstructure](https://github.com/go-viper/mapstructure) - 结构体映射