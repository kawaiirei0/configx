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

## 通用说明

### 配置文件格式

所有示例都使用 YAML 格式的配置文件，configx 也支持：
- JSON
- TOML
- HCL
- INI
- Properties

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

### 更多信息

- 查看主 README: `../README.md`
- 查看架构文档: `../.docs/ARCHITECTURE.md`
- 查看 API 文档: 运行 `go doc github.com/kawaiirei0/configx`
