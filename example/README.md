# ConfigX 示例

本目录包含多个示例，展示如何使用 configx 泛型配置管理库的各种功能。

## 示例列表

### 1. 基础示例 (basic)

**位置**: `example/basic/`

**演示内容**:
- 定义简单的自定义配置结构体
- 使用 `NewManager` 创建泛型配置管理器
- 使用 `LoadConfig` 加载配置文件
- 使用 `GetConfig` 获取类型安全的配置

**运行方式**:
```bash
cd example/basic
go run main.go
```

**适用场景**: 快速入门，了解基本 API 使用

---

### 2. 热重载示例 (hotreload)

**位置**: `example/hotreload/`

**演示内容**:
- 使用 `Init` 方法启动配置热重载
- 配置文件变更时的自动重载
- 配置变更回调函数的使用
- 防抖机制的效果（500ms）

**运行方式**:
```bash
cd example/hotreload
go run main.go
# 程序运行后，修改 config.yaml 文件观察热重载效果
```

**适用场景**: 需要在运行时动态更新配置的应用

---

### 3. 复杂配置示例 (complex)

**位置**: `example/complex/`

**演示内容**:
- 定义包含嵌套结构的复杂配置（Server, Database, Redis, Logging）
- 复杂配置的加载和访问
- 实现自定义 `Clone()` 方法（实现 `Cloneable` 接口）
- 深拷贝验证

**运行方式**:
```bash
cd example/complex
go run main.go
```

**适用场景**: 
- 大型应用的多模块配置管理
- 需要优化配置拷贝性能的场景

---

### 4. 钩子示例 (hooks)

**位置**: `example/hooks/`

**演示内容**:
- 设置不同级别的日志钩子（InitHook, Debug, Info, Warn, Error）
- 钩子在配置生命周期中的触发时机
- 配置变更时的钩子触发
- 集成自定义日志系统

**运行方式**:
```bash
cd example/hooks
go run main.go
# 程序运行后，修改 config.yaml 文件观察钩子触发
```

**适用场景**: 
- 需要监控配置加载和变更事件
- 集成日志系统或监控系统

---

### 5. 环境变量覆盖示例 (env-override)

**位置**: `example/env-override/`

**演示内容**:
- 支持多种配置文件格式（YAML, JSON, TOML）
- 使用环境变量覆盖配置文件中的值
- 自动环境变量读取（AutomaticEnv）
- 环境变量前缀和键名转换
- 敏感信息（API密钥、密码）的安全处理

**运行方式**:
```bash
cd example/env-override
go run main.go
```

**适用场景**: 
- 需要在不同环境（开发、测试、生产）使用不同配置
- 需要保护敏感信息不写入配置文件
- 容器化部署（Docker, Kubernetes）

---

### 6. 特定环境变量绑定示例 (env-bind)

**位置**: `example/env-bind/`

**演示内容**:
- 精确绑定特定配置项到环境变量
- 数据库密码、Redis密码、JWT密钥等敏感信息管理
- AWS凭证等云服务配置
- 生产环境最佳实践

**运行方式**:
```bash
cd example/env-bind
go run main.go
```

**适用场景**: 
- 生产环境部署
- 需要精确控制哪些配置从环境变量读取
- 符合安全合规要求

---

### 7. 多格式配置文件示例 (multi-format)

**位置**: `example/multi-format/`

**演示内容**:
- 支持 YAML、JSON、TOML 等多种格式
- 自动格式识别（根据文件扩展名）
- 相同配置结构，不同文件格式
- 无需额外配置

**运行方式**:
```bash
cd example/multi-format
go run main.go
```

**适用场景**: 
- 需要支持多种配置文件格式
- 从其他系统迁移配置
- 团队偏好不同的配置格式

---

## 通用说明

### 配置文件格式

ConfigX 自动根据文件扩展名识别格式，支持：
- **YAML** (`.yaml`, `.yml`) - 推荐，可读性好
- **JSON** (`.json`) - 标准格式，易于生成
- **TOML** (`.toml`) - 配置文件专用格式
- **HCL** (`.hcl`) - HashiCorp 配置语言
- **INI** (`.ini`) - 传统配置格式
- **Properties** (`.properties`, `.props`, `.prop`) - Java 风格

**示例**:
```go
// YAML 格式
opts.Filename.Set("config.yaml")

// JSON 格式
opts.Filename.Set("config.json")

// TOML 格式
opts.Filename.Set("config.toml")
```

### 自定义配置结构

每个示例都定义了自己的配置结构体，展示了 configx 的泛型特性：

```go
// 定义配置结构
type MyConfig struct {
    Field1 string `mapstructure:"field1"`
    Field2 int    `mapstructure:"field2"`
}

// 创建管理器
manager := configx.NewManager(MyConfig{})
```

### 性能优化

如果配置结构较大或需要频繁调用 `GetConfig()`，建议实现 `Cloneable` 接口：

```go
func (c MyConfig) Clone() MyConfig {
    // 自定义克隆逻辑
    return MyConfig{
        Field1: c.Field1,
        Field2: c.Field2,
    }
}
```

这样可以避免 JSON 序列化的性能开销。

### 环境变量支持

ConfigX 提供两种方式使用环境变量：

#### 1. 自动环境变量（推荐用于开发）
```go
opts := configx.NewOption()
opts.SetEnvPrefix("MYAPP")           // 设置前缀
opts.EnableAutomaticEnv(true)        // 启用自动读取
opts.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
```

配置项 `database.password` 会自动从 `MYAPP_DATABASE_PASSWORD` 读取。

#### 2. 精确绑定（推荐用于生产）
```go
manager.BindEnv("database.password", "DB_PASSWORD")
manager.BindEnv("api.key", "API_KEY")
```

只绑定需要的敏感配置项，更安全可控。

### 更多信息

- 查看主 README: `../README.md`
- 查看并发安全说明: `../CONCURRENCY_SAFETY.md`
- 查看 API 文档: 运行 `go doc github.com/kawaiirei0/configx`
