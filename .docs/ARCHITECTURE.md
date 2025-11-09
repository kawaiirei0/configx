# 架构设计文档

## 项目概述

ConfigX 是一个基于 Go 泛型的轻量级配置管理库，基于 Viper 实现。通过泛型设计，ConfigX 允许开发者在自己的项目中定义配置结构体，而不是使用库内部硬编码的配置模型。这使得 ConfigX 成为一个真正通用的配置管理工具库。

核心特性：
- **泛型设计** - 使用 Go 1.18+ 泛型特性，支持任意配置结构体
- **类型安全** - 编译时类型检查，避免运行时类型断言错误
- **热重载** - 自动监控配置文件变更并重新加载
- **线程安全** - 使用读写锁保证并发访问安全
- **防抖机制** - 避免频繁重载配置文件
- **钩子系统** - 支持多级别事件钩子

## 架构设计

### 核心架构

```
┌──────────────────────────────────────────────────────────────────┐
│                      ConfigX 泛型架构                            │
├──────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌────────────────────────────────────────────────────────┐     │
│  │                    用户层                               │     │
│  │  ┌──────────────────────────────────────────────┐      │     │
│  │  │  用户定义的配置结构体 (T)                     │      │     │
│  │  │  - AppConfig                                 │      │     │
│  │  │  - ServerConfig                              │      │     │
│  │  │  - DatabaseConfig                            │      │     │
│  │  │  - 任意自定义结构体                           │      │     │
│  │  └──────────────────────────────────────────────┘      │     │
│  └────────────────────────────────────────────────────────┘     │
│                            ↓                                     │
│  ┌────────────────────────────────────────────────────────┐     │
│  │                  泛型管理器层                           │     │
│  │  ┌──────────────────────────────────────────────┐      │     │
│  │  │  Manager[T any]                              │      │     │
│  │  │  - config: *T (泛型配置对象)                  │      │     │
│  │  │  - GetConfig() (T, error)                    │      │     │
│  │  │  - LoadConfig() error                        │      │     │
│  │  │  - Init(callback) error                      │      │     │
│  │  │  - SetHook() *Manager[T]                     │      │     │
│  │  └──────────────────────────────────────────────┘      │     │
│  └────────────────────────────────────────────────────────┘     │
│                            ↓                                     │
│  ┌────────────────────────────────────────────────────────┐     │
│  │                  核心功能层                             │     │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐             │     │
│  │  │ 热重载    │  │ 深拷贝    │  │ 钩子系统  │             │     │
│  │  │ - 监控    │  │ - Clone  │  │ - Debug  │             │     │
│  │  │ - 防抖    │  │ - JSON   │  │ - Info   │             │     │
│  │  │ - 回调    │  │          │  │ - Warn   │             │     │
│  │  │          │  │          │  │ - Error  │             │     │
│  │  └──────────┘  └──────────┘  └──────────┘             │     │
│  └────────────────────────────────────────────────────────┘     │
│                            ↓                                     │
│  ┌────────────────────────────────────────────────────────┐     │
│  │                  基础设施层                             │     │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐             │     │
│  │  │  Viper   │  │ RWMutex  │  │ fsnotify │             │     │
│  │  │ 配置解析  │  │ 并发控制  │  │ 文件监控  │             │     │
│  │  └──────────┘  └──────────┘  └──────────┘             │     │
│  └────────────────────────────────────────────────────────┘     │
│                                                                  │
└──────────────────────────────────────────────────────────────────┘
```

### 模块结构

```
configx/
├── manager.go                    # 泛型管理器核心实现
├── init_manager.go               # 管理器初始化和热重载
├── monitor_config_changes.go    # 配置监控和防抖
├── option.go                     # 配置选项定义
├── const.go                      # 默认常量
├── hook.go                       # 钩子机制
├── context.go                    # 上下文管理
├── cloneable.go                  # Cloneable 接口定义
├── errors.go                     # 错误类型定义
├── utils.go                      # 工具函数
├── utils_manager.go              # 管理器工具
├── manager_file.go               # 文件操作
├── manager_update_field.go       # 字段更新
└── utils/
    ├── singleton.go              # 单例工具
    └── user_config_path.go       # 配置路径工具
```

注意：v2.x 版本移除了以下文件：
- `config.go` - 硬编码的配置结构体（已删除）
- `configure/` - 硬编码的配置模块目录（已删除）
- `get_config.go` - 全局单例函数（已删除）

## 核心组件

### 1. 泛型配置管理器 (Manager[T])

**职责：**
- 管理任意类型配置的生命周期
- 提供类型安全的配置访问
- 处理配置文件加载和更新
- 实现防抖机制和热重载

**关键特性：**
- 泛型设计支持任意配置结构体
- 编译时类型安全，无需类型断言
- 读写锁保证并发安全
- 防抖机制避免频繁重载
- 支持自定义 Clone 方法优化性能

```go
type Manager[T any] struct {
    config              *T            // 泛型配置对象
    vp                  *viper.Viper  // Viper 实例
    rwMutex             sync.RWMutex  // 读写锁
    lastChange          time.Time     // 上次触发时间（用于防抖）
    debounceDur         time.Duration // 防抖间隔
    hooks               *Hook         // 钩子系统
    pathName            string        // 配置文件路径
    opts                *Option       // 配置选项
    optsInit            bool          // 选项初始化标志
    validateConfigValue bool          // 验证标志
    defaultConfig       any           // 默认配置
}
```

**泛型设计原理：**

使用 Go 1.18+ 的泛型特性，Manager 接受类型参数 `T any`，这使得：
1. 配置类型在编译时确定，提供类型安全
2. 无需在库内部定义配置结构，由用户自定义
3. `GetConfig()` 返回具体类型 `T`，无需类型断言
4. 支持任意复杂的配置结构，包括嵌套结构

### 2. 用户自定义配置结构

**v2.x 重大变更：**

ConfigX 不再在库内部定义配置结构体。配置结构完全由用户在自己的项目中定义。

**设计原则：**
- 用户在项目中定义配置结构体
- 使用 `mapstructure` 标签映射 YAML 字段
- 支持任意复杂的嵌套结构
- 可选实现 `Cloneable[T]` 接口优化性能

**示例：**

```go
// 用户项目中定义
type AppConfig struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
}

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

// 创建管理器
manager := configx.NewManager(AppConfig{})
```

### 3. 配置选项 (Option)

**职责：**
- 提供灵活的配置选项
- 支持默认值设置
- 实现类型安全的选项管理

**特性：**
- 类型化选项定义
- 默认值自动初始化
- 支持链式调用

```go
type Option struct {
    Filename    OptionString       // 配置文件名
    Filepath    OptionString       // 配置文件路径
    DebounceDur OptionTimeDuration // 防抖间隔
}

// 使用示例
opts := configx.NewOption()
opts.Filename.Set("config.yaml")
opts.Filepath.Set("./configs")
opts.DebounceDur.Set(1000 * configx.OptionDateMillisecond)
```

### 4. 深拷贝机制

**职责：**
- 提供配置的深拷贝，避免并发修改问题
- 支持自定义高效克隆方法

**实现策略：**

1. **优先使用自定义 Clone 方法**：
   ```go
   type Cloneable[T any] interface {
       Clone() T
   }
   
   // 如果配置类型实现了 Cloneable 接口
   if cloneable, ok := any(*m.config).(Cloneable[T]); ok {
       return cloneable.Clone(), nil
   }
   ```

2. **回退到 JSON 序列化**：
   ```go
   // 默认使用 JSON 序列化/反序列化实现深拷贝
   data, _ := json.Marshal(*m.config)
   var copy T
   json.Unmarshal(data, &copy)
   return copy, nil
   ```

**性能考虑：**
- JSON 序列化简单通用，但性能较低
- 对于复杂配置结构，建议实现自定义 Clone 方法
- 自定义 Clone 可以避免序列化开销，提升 10-100 倍性能

### 5. 文件监控与热重载

**职责：**
- 监听配置文件变更
- 触发配置重载
- 实现防抖机制
- 执行用户回调

**实现机制：**
- 基于 fsnotify 的文件系统监控
- 防抖算法避免频繁触发（默认 800ms）
- 异步处理避免阻塞主流程
- 重载失败时保持原有配置

**防抖逻辑：**
```go
func (m *Manager[T]) monitorConfigChanges(callback func(*Context)) {
    m.vp.OnConfigChange(func(e fsnotify.Event) {
        now := time.Now()
        
        // 防抖检查
        if now.Sub(m.lastChange) < m.debounceDur {
            return
        }
        m.lastChange = now
        
        // 重新加载配置
        if err := m.LoadConfig(); err != nil {
            // 触发错误钩子
            return
        }
        
        // 执行回调
        if callback != nil {
            callback(&Context{Manager: m})
        }
    })
}
```

### 6. 钩子系统

**职责：**
- 提供事件通知机制
- 支持多级别日志记录
- 便于集成现有日志系统

**钩子级别：**
```go
const (
    InitHook HookPattern = iota  // 初始化事件
    Debug                        // 调试信息
    Info                         // 一般信息
    Warn                         // 警告信息
    Error                        // 错误信息
)
```

**使用示例：**
```go
manager.SetHook(configx.Info, func(ctx configx.HookContext) {
    log.Printf("[INFO] %s", ctx.Message)
}).SetHook(configx.Error, func(ctx configx.HookContext) {
    log.Printf("[ERROR] %s", ctx.Message)
})
```

## 数据流设计

### 配置加载流程

```
用户调用 manager.LoadConfig()
        ↓
获取写锁 (rwMutex.Lock)
        ↓
配置 Viper 实例
        ↓
读取配置文件 (vp.ReadInConfig)
        ↓
创建泛型配置实例 (var newConfig T)
        ↓
Viper 解析 YAML 到泛型类型 (vp.Unmarshal(&newConfig))
        ↓
更新配置指针 (m.config = &newConfig)
        ↓
释放写锁
        ↓
返回结果
```

### 配置获取流程

```
用户调用 manager.GetConfig()
        ↓
获取读锁 (rwMutex.RLock)
        ↓
检查配置是否已初始化
        ↓
检查是否实现 Cloneable[T] 接口
        ↓
是：调用自定义 Clone() 方法
否：使用 JSON 序列化深拷贝
        ↓
释放读锁
        ↓
返回配置副本 (T, error)
```

### 热更新流程

```
用户调用 manager.Init(callback)
        ↓
加载初始配置 (LoadConfig)
        ↓
启动文件监控 (monitorConfigChanges)
        ↓
fsnotify 监听文件变更
        ↓
文件变更事件触发
        ↓
防抖检查 (now.Sub(lastChange) < debounceDur)
        ↓
触发 Info 钩子（开始重载）
        ↓
重新加载配置 (LoadConfig)
        ↓
成功：触发 Info 钩子 + 执行用户回调
失败：触发 Error 钩子 + 保持原配置
        ↓
更新 lastChange 时间戳
```

### 泛型类型流转

```
用户定义配置类型 T
        ↓
创建 Manager[T]
        ↓
YAML 文件 → Viper → Unmarshal → T
        ↓
存储为 *T
        ↓
GetConfig() → Clone/JSON → T (副本)
        ↓
用户使用类型安全的 T
```

## 并发设计

### 线程安全机制

ConfigX 使用读写锁 (sync.RWMutex) 保证并发安全：

1. **读操作 (GetConfig)**：
   - 使用读锁 `rwMutex.RLock()`
   - 允许多个 goroutine 同时读取
   - 不阻塞其他读操作

2. **写操作 (LoadConfig)**：
   - 使用写锁 `rwMutex.Lock()`
   - 独占访问，阻塞所有读写操作
   - 确保配置更新的原子性

3. **防抖机制**：
   - 使用 `lastChange` 时间戳避免频繁重载
   - 在热重载场景下避免并发重复加载

### 性能优化

1. **配置深拷贝**：
   - `GetConfig()` 返回配置副本，避免外部修改
   - 支持自定义 Clone 方法，避免 JSON 序列化开销
   - 对于简单结构，JSON 序列化性能可接受
   - 对于复杂结构，自定义 Clone 可提升 10-100 倍性能

2. **读写分离**：
   - 读操作不阻塞其他读操作
   - 写操作较少（仅在加载/重载时）
   - 适合读多写少的场景

3. **防抖优化**：
   - 避免短时间内多次重载
   - 减少文件 I/O 和解析开销
   - 可自定义防抖间隔

## 泛型设计的技术决策

### 决策 1: 使用泛型而非接口

**原因：**
- 提供编译时类型安全
- 避免运行时类型断言
- 更好的 IDE 支持和代码补全
- 性能更优（无接口调用开销）

**权衡：**
- 需要 Go 1.18+
- 泛型语法略复杂
- 但带来的类型安全和开发体验提升值得

### 决策 2: 不提供全局单例

**原因：**
- Go 泛型不支持泛型全局变量
- 强行实现会增加复杂度和运行时开销
- 鼓励依赖注入等更好的设计模式

**影响：**
- 用户需要自行管理 Manager 实例
- 可以在应用层实现单例模式

**示例：**
```go
// 用户项目中实现单例
var configManager = configx.NewManager(AppConfig{})

func GetConfigManager() *configx.Manager[AppConfig] {
    return configManager
}
```

### 决策 3: 使用 JSON 序列化作为默认深拷贝

**原因：**
- 简单通用，适用于所有可序列化类型
- 无需用户额外实现
- 对于大多数配置场景性能可接受

**优化：**
- 提供 `Cloneable[T]` 接口
- 用户可实现自定义高效克隆
- 库自动检测并使用自定义方法

### 决策 4: 配置结构由用户定义

**原因：**
- 真正的通用性，不限制配置结构
- 避免库内部硬编码配置模型
- 用户可以根据项目需求自由设计

**优势：**
- 灵活性最大化
- 类型安全
- 易于维护和扩展

### 决策 5: 主版本号升级到 v2

**原因：**
- API 有破坏性变更
- 遵循语义化版本规范
- 清晰区分新旧版本

**迁移支持：**
- 提供详细迁移指南
- 示例代码展示新旧对比
- 文档说明所有变更

## 错误处理

### 错误类型

- **文件错误：** 配置文件不存在或格式错误
- **解析错误：** YAML 解析失败或结构体映射错误
- **权限错误：** 文件读取权限不足
- **并发错误：** 读写锁获取失败

### 错误处理策略

- **快速失败：** 配置加载失败时立即返回错误
- **优雅降级：** 热更新失败时保持原有配置
- **日志记录：** 通过钩子机制记录错误信息

## 安全设计

### 数据安全

- **配置副本：** 避免直接暴露内部配置对象
- **只读访问：** 外部只能通过副本访问配置
- **类型安全：** 结构体定义确保类型正确性

### 文件安全

- **路径验证：** 配置文件路径安全检查
- **权限控制：** 文件读取权限验证
- **格式验证：** YAML 格式和内容验证

## 部署考虑

### 环境配置

支持不同环境的配置管理：

- **开发环境：** dev 配置，调试信息完整
- **测试环境：** test 配置，模拟生产环境
- **生产环境：** prod 配置，性能优化

### 监控建议

- **配置变更监控：** 记录配置变更历史
- **性能监控：** 监控配置加载时间和频率
- **错误监控：** 监控配置加载错误率

## 扩展性设计

### 当前扩展点

1. **自定义配置结构**：
   - 用户可定义任意复杂的配置结构
   - 支持嵌套结构、数组、映射等

2. **自定义 Clone 方法**：
   - 实现 `Cloneable[T]` 接口
   - 优化深拷贝性能

3. **钩子系统**：
   - 支持多级别钩子
   - 便于集成日志系统

4. **配置选项**：
   - 自定义文件路径
   - 自定义防抖间隔

### 未来扩展方向

1. **配置验证**：
   ```go
   type Validatable interface {
       Validate() error
   }
   
   // 在 LoadConfig 后自动验证
   if v, ok := any(newConfig).(Validatable); ok {
       if err := v.Validate(); err != nil {
           return err
       }
   }
   ```

2. **多格式支持**：
   - JSON、TOML、INI 等格式
   - Viper 已支持，只需暴露接口

3. **配置加密**：
   - 敏感字段加密存储
   - 自动解密机制

4. **远程配置**：
   - 支持从配置中心拉取
   - etcd、Consul 等集成

5. **配置合并**：
   - 多个配置源合并
   - 优先级控制