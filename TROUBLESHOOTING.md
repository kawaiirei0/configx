# 故障排除指南

本指南帮助你解决使用 ConfigX 时可能遇到的常见问题。

## 目录

- [配置加载问题](#配置加载问题)
- [配置文件格式问题](#配置文件格式问题)
- [泛型类型推断问题](#泛型类型推断问题)
- [热重载问题](#热重载问题)
- [深拷贝性能问题](#深拷贝性能问题)
- [并发访问问题](#并发访问问题)
- [钩子问题](#钩子问题)
- [其他常见问题](#其他常见问题)

---

## 配置加载问题

### 问题：配置文件不存在

**错误信息：**
```
配置文件不存在: ./configs/config.yaml, 错误: open ./configs/config.yaml: no such file or directory
```

**原因：**
- 配置文件路径错误
- 配置文件名错误
- 配置文件不存在

**解决方案：**

1. 检查文件路径是否正确：
   ```go
   opts := configx.NewOption()
   opts.Filename.Set("config.yaml")  // 确保文件名正确
   opts.Filepath.Set("./configs")    // 确保路径正确
   ```

2. 确认文件存在：
   ```bash
   ls -la ./configs/config.yaml
   ```

3. 使用绝对路径（不推荐，但可用于调试）：
   ```go
   opts.Filepath.Set("/absolute/path/to/configs")
   ```

4. 检查工作目录：
   ```go
   wd, _ := os.Getwd()
   fmt.Println("当前工作目录:", wd)
   ```

---

### 问题：配置未初始化

**错误信息：**
```
配置未初始化
```

**原因：**
- 在调用 `LoadConfig()` 或 `Init()` 之前调用了 `GetConfig()`

**解决方案：**

1. 确保先加载配置：
   ```go
   manager := configx.NewManager(AppConfig{})
   
   // 必须先加载配置
   if err := manager.LoadConfig(); err != nil {
       log.Fatal(err)
   }
   
   // 然后才能获取配置
   config, err := manager.GetConfig()
   ```

2. 或使用 `Init` 方法：
   ```go
   manager := configx.NewManager(AppConfig{})
   
   // Init 会自动加载配置
   if err := manager.Init(nil); err != nil {
       log.Fatal(err)
   }
   
   config, err := manager.GetConfig()
   ```

3. 添加错误检查：
   ```go
   config, err := manager.GetConfig()
   if errors.Is(err, configx.ErrConfigNotInitialized) {
       log.Println("配置未初始化，正在加载...")
       if err := manager.LoadConfig(); err != nil {
           log.Fatal(err)
       }
       config, err = manager.GetConfig()
   }
   ```

---

### 问题：配置解析失败

**错误信息：**
```
配置解析失败: 文件 ./configs/config.yaml, 错误: yaml: unmarshal errors:
  line 2: cannot unmarshal !!str `invalid` into int
```

**原因：**
- YAML 文件格式错误
- 配置结构与 YAML 不匹配
- 类型转换失败

**解决方案：**

1. 检查 YAML 格式：
   ```yaml
   # 错误：端口应该是数字
   port: "8080"
   
   # 正确
   port: 8080
   ```

2. 检查结构体标签：
   ```go
   type AppConfig struct {
       Port int `mapstructure:"port"`  // 确保标签与 YAML 字段匹配
   }
   ```

3. 使用 YAML 验证工具：
   ```bash
   # 使用 yamllint 验证
   yamllint config.yaml
   ```

4. 检查嵌套结构：
   ```yaml
   # YAML 文件
   server:
     host: "localhost"
     port: 8080
   ```
   
   ```go
   // 对应的结构体
   type AppConfig struct {
       Server ServerConfig `mapstructure:"server"`
   }
   
   type ServerConfig struct {
       Host string `mapstructure:"host"`
       Port int    `mapstructure:"port"`
   }
   ```

---

## 配置文件格式问题

### 问题：YAML 缩进错误

**错误信息：**
```
yaml: line 3: mapping values are not allowed in this context
```

**原因：**
- YAML 缩进不正确（必须使用空格，不能使用 Tab）
- 冒号后缺少空格

**解决方案：**

1. 使用空格而非 Tab：
   ```yaml
   # 错误：使用了 Tab
   server:
   	host: "localhost"
   
   # 正确：使用空格
   server:
     host: "localhost"
   ```

2. 冒号后添加空格：
   ```yaml
   # 错误：冒号后没有空格
   port:8080
   
   # 正确
   port: 8080
   ```

3. 检查缩进层级：
   ```yaml
   # 正确的缩进
   server:
     host: "localhost"
     port: 8080
   database:
     host: "localhost"
     port: 3306
   ```

---

### 问题：字段名不匹配

**症状：**
- 配置加载成功，但字段值为零值

**原因：**
- 结构体标签与 YAML 字段名不匹配
- 大小写不匹配

**解决方案：**

1. 确保标签匹配：
   ```yaml
   # YAML 文件
   app_name: "MyApp"
   ```
   
   ```go
   // 结构体
   type AppConfig struct {
       AppName string `mapstructure:"app_name"`  // 标签必须匹配
   }
   ```

2. 注意大小写：
   ```yaml
   # YAML 使用小写
   port: 8080
   ```
   
   ```go
   // 标签也使用小写
   type AppConfig struct {
       Port int `mapstructure:"port"`  // 不是 "Port"
   }
   ```

3. 使用调试输出：
   ```go
   config, _ := manager.GetConfig()
   fmt.Printf("配置内容: %+v\n", config)
   ```

---

## 泛型类型推断问题

### 问题：类型推断失败

**错误信息：**
```
cannot infer T
```

**原因：**
- 没有提供足够的类型信息

**解决方案：**

1. 显式指定类型参数：
   ```go
   // 如果编译器无法推断
   manager := configx.NewManager[AppConfig](AppConfig{})
   ```

2. 使用具体类型：
   ```go
   // 推荐：让编译器自动推断
   manager := configx.NewManager(AppConfig{})
   ```

---

### 问题：类型不匹配

**错误信息：**
```
cannot use manager (variable of type *Manager[AppConfig]) as *Manager[OtherConfig]
```

**原因：**
- 尝试将一个类型的 Manager 赋值给另一个类型

**解决方案：**

1. 确保类型一致：
   ```go
   // 错误
   var manager *configx.Manager[AppConfig]
   manager = configx.NewManager(OtherConfig{})  // 类型不匹配
   
   // 正确
   var manager *configx.Manager[AppConfig]
   manager = configx.NewManager(AppConfig{})
   ```

2. 使用接口（如果需要多态）：
   ```go
   type ConfigManager interface {
       LoadConfig() error
   }
   
   var manager ConfigManager
   manager = configx.NewManager(AppConfig{})
   ```

---

## 热重载问题

### 问题：配置变更不触发重载

**症状：**
- 修改配置文件后，应用没有重新加载配置

**原因：**
- 没有调用 `Init` 方法启动监控
- 防抖时间内多次修改
- 文件监控失败

**解决方案：**

1. 确保调用了 `Init`：
   ```go
   // 必须调用 Init 启动监控
   err := manager.Init(func(ctx *configx.Context) {
       fmt.Println("配置已更新")
   })
   ```

2. 检查防抖时间：
   ```go
   opts := configx.NewOption()
   // 减小防抖时间以便测试
   opts.DebounceDur.Set(200 * configx.OptionDateMillisecond)
   manager.SetOption(opts)
   ```

3. 检查文件监控是否工作：
   ```go
   manager.SetHook(configx.Info, func(ctx configx.HookContext) {
       fmt.Printf("[INFO] %s\n", ctx.Message)
   })
   ```

4. 确保文件保存成功：
   ```bash
   # 使用 touch 触发文件变更
   touch ./configs/config.yaml
   ```

---

### 问题：热重载频繁触发

**症状：**
- 配置重载过于频繁
- 日志显示多次重载

**原因：**
- 防抖时间设置过短
- 编辑器保存时产生多个文件事件

**解决方案：**

1. 增加防抖时间：
   ```go
   opts := configx.NewOption()
   // 增加防抖时间到 1 秒
   opts.DebounceDur.Set(1000 * configx.OptionDateMillisecond)
   manager.SetOption(opts)
   ```

2. 使用钩子监控重载频率：
   ```go
   var reloadCount int
   manager.SetHook(configx.Info, func(ctx configx.HookContext) {
       reloadCount++
       fmt.Printf("重载次数: %d, 消息: %s\n", reloadCount, ctx.Message)
   })
   ```

---

### 问题：热重载失败但应用继续运行

**症状：**
- 配置文件有错误，但应用没有崩溃
- 使用的是旧配置

**原因：**
- 这是设计行为：重载失败时保持原有配置

**解决方案：**

1. 使用错误钩子监控：
   ```go
   manager.SetHook(configx.Error, func(ctx configx.HookContext) {
       log.Printf("[ERROR] 配置重载失败: %s\n", ctx.Message)
       // 可以发送告警
   })
   ```

2. 在回调中检查配置有效性：
   ```go
   manager.Init(func(ctx *configx.Context) {
       config, err := manager.GetConfig()
       if err != nil {
           log.Printf("获取配置失败: %v", err)
           return
       }
       
       // 验证配置
       if config.Port < 1024 || config.Port > 65535 {
           log.Printf("警告：端口号无效: %d", config.Port)
       }
   })
   ```

---

## 深拷贝性能问题

### 问题：GetConfig 性能较差

**症状：**
- `GetConfig()` 调用耗时较长
- 高并发场景下性能瓶颈

**原因：**
- 默认使用 JSON 序列化进行深拷贝
- 配置结构复杂，序列化开销大

**解决方案：**

1. 实现 `Cloneable` 接口：
   ```go
   type AppConfig struct {
       Port int    `mapstructure:"port"`
       Host string `mapstructure:"host"`
   }
   
   // 实现自定义克隆方法
   func (c AppConfig) Clone() AppConfig {
       return AppConfig{
           Port: c.Port,
           Host: c.Host,
       }
   }
   ```

2. 对于复杂结构，手动深拷贝：
   ```go
   type ComplexConfig struct {
       Tags     []string          `mapstructure:"tags"`
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
           Tags:     tags,
           Metadata: metadata,
       }
   }
   ```

3. 性能测试：
   ```go
   func BenchmarkGetConfig(b *testing.B) {
       manager := configx.NewManager(AppConfig{})
       manager.LoadConfig()
       
       b.ResetTimer()
       for i := 0; i < b.N; i++ {
           _, _ = manager.GetConfig()
       }
   }
   ```

4. 性能对比：
   ```
   # 使用 JSON 序列化
   BenchmarkGetConfig-8    100000    10000 ns/op
   
   # 使用自定义 Clone
   BenchmarkGetConfig-8    10000000    100 ns/op
   ```

---

### 问题：深拷贝不完整

**症状：**
- 修改返回的配置影响了原配置
- 并发访问出现数据竞争

**原因：**
- 自定义 Clone 方法没有正确处理引用类型
- 浅拷贝而非深拷贝

**解决方案：**

1. 正确处理 slice：
   ```go
   // 错误：浅拷贝
   func (c Config) Clone() Config {
       return Config{
           Tags: c.Tags,  // 共享底层数组
       }
   }
   
   // 正确：深拷贝
   func (c Config) Clone() Config {
       tags := make([]string, len(c.Tags))
       copy(tags, c.Tags)
       return Config{
           Tags: tags,
       }
   }
   ```

2. 正确处理 map：
   ```go
   // 错误：浅拷贝
   func (c Config) Clone() Config {
       return Config{
           Metadata: c.Metadata,  // 共享底层数据
       }
   }
   
   // 正确：深拷贝
   func (c Config) Clone() Config {
       metadata := make(map[string]string, len(c.Metadata))
       for k, v := range c.Metadata {
           metadata[k] = v
       }
       return Config{
           Metadata: metadata,
       }
   }
   ```

3. 正确处理嵌套结构：
   ```go
   type Config struct {
       Server ServerConfig `mapstructure:"server"`
   }
   
   func (c Config) Clone() Config {
       return Config{
           Server: c.Server.Clone(),  // 递归克隆
       }
   }
   ```

4. 测试深拷贝：
   ```go
   func TestDeepCopy(t *testing.T) {
       manager := configx.NewManager(AppConfig{})
       manager.LoadConfig()
       
       config1, _ := manager.GetConfig()
       config2, _ := manager.GetConfig()
       
       // 修改 config1
       config1.Tags[0] = "modified"
       
       // config2 不应该受影响
       if config2.Tags[0] == "modified" {
           t.Error("深拷贝失败：配置被共享")
       }
   }
   ```

---

## 并发访问问题

### 问题：数据竞争

**错误信息：**
```
WARNING: DATA RACE
```

**原因：**
- 直接修改 `GetConfig()` 返回的配置
- 多个 goroutine 同时访问配置

**解决方案：**

1. 不要修改返回的配置：
   ```go
   // 错误：修改返回的配置
   config, _ := manager.GetConfig()
   config.Port = 9090  // 不要这样做
   
   // 正确：只读访问
   config, _ := manager.GetConfig()
   port := config.Port  // 只读取值
   ```

2. 每次都获取新副本：
   ```go
   // 每个 goroutine 获取自己的副本
   go func() {
       config, _ := manager.GetConfig()
       // 使用 config
   }()
   
   go func() {
       config, _ := manager.GetConfig()
       // 使用 config
   }()
   ```

3. 使用竞态检测：
   ```bash
   go run -race main.go
   ```

---

### 问题：死锁

**症状：**
- 程序挂起
- goroutine 阻塞

**原因：**
- 在回调函数中调用 `GetConfig()` 可能导致死锁

**解决方案：**

1. 避免在回调中调用 `GetConfig()`：
   ```go
   // 可能导致死锁
   manager.Init(func(ctx *configx.Context) {
       config, _ := manager.GetConfig()  // 可能死锁
   })
   
   // 推荐：在回调外获取配置
   manager.Init(func(ctx *configx.Context) {
       fmt.Println("配置已更新")
   })
   
   // 在回调外获取
   config, _ := manager.GetConfig()
   ```

2. 使用 goroutine：
   ```go
   manager.Init(func(ctx *configx.Context) {
       go func() {
           config, _ := manager.GetConfig()
           // 使用 config
       }()
   })
   ```

---

## 钩子问题

### 问题：钩子没有触发

**症状：**
- 设置了钩子，但没有输出

**原因：**
- 钩子级别不匹配
- 钩子设置在事件发生之后

**解决方案：**

1. 在初始化前设置钩子：
   ```go
   manager := configx.NewManager(AppConfig{})
   
   // 先设置钩子
   manager.SetHook(configx.Info, func(ctx configx.HookContext) {
       fmt.Println(ctx.Message)
   })
   
   // 再初始化
   manager.Init(nil)
   ```

2. 设置所有级别的钩子：
   ```go
   manager.SetHook(configx.Debug, handler)
   manager.SetHook(configx.Info, handler)
   manager.SetHook(configx.Warn, handler)
   manager.SetHook(configx.Error, handler)
   ```

3. 检查钩子函数：
   ```go
   manager.SetHook(configx.Info, func(ctx configx.HookContext) {
       // 确保这里有代码
       fmt.Printf("[INFO] %s\n", ctx.Message)
   })
   ```

---

### 问题：钩子执行顺序

**症状：**
- 钩子执行顺序不符合预期

**原因：**
- 钩子是按设置顺序执行的

**解决方案：**

1. 按期望顺序设置钩子：
   ```go
   // 先设置的先执行
   manager.SetHook(configx.Info, func(ctx configx.HookContext) {
       fmt.Println("第一个钩子")
   })
   
   manager.SetHook(configx.Info, func(ctx configx.HookContext) {
       fmt.Println("第二个钩子")
   })
   ```

2. 使用单个钩子处理多个任务：
   ```go
   manager.SetHook(configx.Info, func(ctx configx.HookContext) {
       // 任务 1
       log.Println(ctx.Message)
       
       // 任务 2
       metrics.RecordEvent("config_reload")
       
       // 任务 3
       notifyService(ctx.Message)
   })
   ```

---

## 其他常见问题

### 问题：配置选项不生效

**症状：**
- 设置了配置选项，但没有效果

**原因：**
- 在 `LoadConfig()` 或 `Init()` 之后设置选项

**解决方案：**

1. 在加载配置前设置选项：
   ```go
   manager := configx.NewManager(AppConfig{})
   
   // 先设置选项
   opts := configx.NewOption()
   opts.Filename.Set("config.yaml")
   manager.SetOption(opts)
   
   // 再加载配置
   manager.LoadConfig()
   ```

---

### 问题：内存泄漏

**症状：**
- 内存使用持续增长
- 程序运行一段时间后变慢

**原因：**
- 配置副本没有被释放
- 钩子函数持有大量引用

**解决方案：**

1. 及时释放配置副本：
   ```go
   func processConfig() {
       config, _ := manager.GetConfig()
       // 使用 config
       // 函数结束时 config 会被 GC 回收
   }
   ```

2. 避免在钩子中持有大量数据：
   ```go
   // 错误：持有大量数据
   var history []string
   manager.SetHook(configx.Info, func(ctx configx.HookContext) {
       history = append(history, ctx.Message)  // 无限增长
   })
   
   // 正确：限制大小
   var history []string
   const maxHistory = 100
   manager.SetHook(configx.Info, func(ctx configx.HookContext) {
       history = append(history, ctx.Message)
       if len(history) > maxHistory {
           history = history[1:]  // 保持固定大小
       }
   })
   ```

3. 使用内存分析工具：
   ```bash
   go test -memprofile=mem.prof
   go tool pprof mem.prof
   ```

---

### 问题：Go 版本不兼容

**错误信息：**
```
type parameter requires go1.18 or later
```

**原因：**
- ConfigX v2.x 需要 Go 1.18 或更高版本

**解决方案：**

1. 升级 Go 版本：
   ```bash
   # 检查当前版本
   go version
   
   # 升级到 Go 1.18+
   # 参考 https://go.dev/doc/install
   ```

2. 或使用 v1.x 版本：
   ```bash
   go get github.com/kawaiirei0/configx@v1
   ```

---

### 问题：导入路径错误

**错误信息：**
```
package github.com/kawaiirei0/configx: cannot find package
```

**原因：**
- 导入路径错误
- 依赖未下载

**解决方案：**

1. 检查导入路径：
   ```go
   import "github.com/kawaiirei0/configx"  // 正确
   ```

2. 下载依赖：
   ```bash
   go mod tidy
   go mod download
   ```

3. 清理缓存：
   ```bash
   go clean -modcache
   go mod download
   ```

---

## 调试技巧

### 1. 启用详细日志

```go
manager.SetHook(configx.Debug, func(ctx configx.HookContext) {
    log.Printf("[DEBUG] %s", ctx.Message)
}).SetHook(configx.Info, func(ctx configx.HookContext) {
    log.Printf("[INFO] %s", ctx.Message)
}).SetHook(configx.Warn, func(ctx configx.HookContext) {
    log.Printf("[WARN] %s", ctx.Message)
}).SetHook(configx.Error, func(ctx configx.HookContext) {
    log.Printf("[ERROR] %s", ctx.Message)
})
```

### 2. 打印配置内容

```go
config, err := manager.GetConfig()
if err != nil {
    log.Fatal(err)
}

// 打印完整配置
fmt.Printf("配置内容: %+v\n", config)

// 使用 JSON 格式化
data, _ := json.MarshalIndent(config, "", "  ")
fmt.Println(string(data))
```

### 3. 检查文件路径

```go
opts := configx.NewOption()
opts.Filename.Set("config.yaml")
opts.Filepath.Set("./configs")

// 打印完整路径
fullPath := opts.File()
fmt.Printf("配置文件路径: %s\n", fullPath)

// 检查文件是否存在
if _, err := os.Stat(fullPath); os.IsNotExist(err) {
    fmt.Println("配置文件不存在")
}
```

### 4. 使用竞态检测

```bash
# 运行时检测
go run -race main.go

# 测试时检测
go test -race ./...
```

### 5. 性能分析

```bash
# CPU 分析
go test -cpuprofile=cpu.prof -bench=.
go tool pprof cpu.prof

# 内存分析
go test -memprofile=mem.prof -bench=.
go tool pprof mem.prof
```

---

## 获取帮助

如果以上方法都无法解决你的问题：

1. **查看文档**：
   - [README.md](./README.md) - 快速开始
   - [API.md](./API.md) - API 参考
   - [ARCHITECTURE.md](./.docs/ARCHITECTURE.md) - 架构设计
   - [MIGRATION.md](./MIGRATION.md) - 迁移指南

2. **查看示例**：
   - [example/basic](./example/basic) - 基础用法
   - [example/complex](./example/complex) - 复杂配置
   - [example/hotreload](./example/hotreload) - 热重载
   - [example/hooks](./example/hooks) - 钩子系统

3. **提交 Issue**：
   - 在 GitHub 上提交 Issue
   - 提供完整的错误信息和复现步骤
   - 附上相关代码和配置文件

4. **社区支持**：
   - 查看已有的 Issues 和 Discussions
   - 参与社区讨论

---

## 常见错误速查表

| 错误信息 | 可能原因 | 解决方案 |
|---------|---------|---------|
| 配置文件不存在 | 文件路径错误 | 检查文件路径和文件名 |
| 配置未初始化 | 未调用 LoadConfig | 先调用 LoadConfig 或 Init |
| 配置解析失败 | YAML 格式错误 | 检查 YAML 格式和结构体标签 |
| 类型推断失败 | 缺少类型信息 | 显式指定类型参数 |
| 数据竞争 | 并发访问配置 | 使用 GetConfig 获取副本 |
| 钩子未触发 | 钩子设置时机错误 | 在初始化前设置钩子 |
| 热重载不工作 | 未调用 Init | 使用 Init 启动监控 |
| 性能问题 | 使用 JSON 深拷贝 | 实现 Cloneable 接口 |
| Go 版本错误 | Go < 1.18 | 升级到 Go 1.18+ |
| 导入路径错误 | 路径不正确 | 检查导入路径 |

---

希望本指南能帮助你解决问题。如果还有疑问，欢迎查看其他文档或提交 Issue！
