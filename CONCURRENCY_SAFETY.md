# 并发安全说明

## 概述

ConfigX v2 已经过完整的并发安全审查和修复，可以安全地在多 goroutine 环境中使用。

## 线程安全保证

### ✅ 完全线程安全的操作

以下操作可以在多个 goroutine 中并发调用：

```go
manager := configx.NewManager(MyConfig{})

// 并发安全
go func() { config, _ := manager.GetConfig() }()
go func() { config, _ := manager.GetConfig() }()

// 并发安全
go func() { manager.SetHook(configx.Info, handler1) }()
go func() { manager.SetHook(configx.Error, handler2) }()

// 并发安全
go func() { manager.SetOption(opts1) }()
go func() { manager.SetOption(opts2) }()

// 并发安全
go func() { manager.LoadConfig() }()
go func() { config, _ := manager.GetConfig() }()
```

### 🔒 同步机制

1. **配置读写**: 使用 `sync.RWMutex` 保护
   - 多个 goroutine 可以同时读取配置
   - 写入时会阻塞所有读写操作

2. **钩子系统**: 使用 `sync.RWMutex` 保护
   - 多个 goroutine 可以同时执行钩子
   - 设置钩子时会阻塞所有钩子操作

3. **选项初始化**: 使用 `sync.Mutex` 保护
   - 确保选项只初始化一次
   - 防止并发初始化冲突

4. **防抖机制**: 使用 `atomic.Int64`
   - 无锁高性能时间戳访问
   - 完全线程安全

## 使用示例

### 并发读取配置

```go
manager := configx.NewManager(AppConfig{})
manager.LoadConfig()

// 在多个 goroutine 中安全读取
for i := 0; i < 100; i++ {
    go func() {
        config, err := manager.GetConfig()
        if err != nil {
            log.Println(err)
            return
        }
        // 使用 config
        fmt.Println(config.AppName)
    }()
}
```

### 热重载 + 并发访问

```go
manager := configx.NewManager(AppConfig{})

// 启动热重载
manager.Init(func(ctx *configx.Context) {
    log.Println("配置已更新")
})

// 在其他 goroutine 中持续读取
go func() {
    ticker := time.NewTicker(1 * time.Second)
    for range ticker.C {
        config, _ := manager.GetConfig()
        // 总是能获取到最新的配置
        fmt.Println(config.Version)
    }
}()
```

### 动态设置钩子

```go
manager := configx.NewManager(AppConfig{})

// 可以在运行时动态修改钩子
go func() {
    manager.SetHook(configx.Info, func(ctx configx.HookContext) {
        log.Println("Info:", ctx.Message)
    })
}()

go func() {
    manager.SetHook(configx.Error, func(ctx configx.HookContext) {
        log.Println("Error:", ctx.Message)
    })
}()
```

## 性能特性

### 高性能场景

1. **并发读取**: 使用读写锁，多个 goroutine 可以同时读取配置
2. **防抖检查**: 使用原子操作，无锁访问
3. **钩子执行**: 在锁外执行，不阻塞其他操作

### 性能建议

1. **频繁读取**: 使用 `GetConfig()` 获取配置副本，避免重复加载
2. **自定义克隆**: 实现 `Cloneable` 接口可以提升 `GetConfig()` 性能
3. **钩子函数**: 保持钩子函数简短，避免长时间阻塞

## 测试验证

### 运行并发测试

```bash
# 运行所有测试
go test -v ./...

# 只运行并发测试
go test -v -run Concurrent ./...

# 使用 race detector（需要 GCC）
go test -race ./...
```

### 测试覆盖

- ✅ 并发读取配置
- ✅ 并发设置钩子
- ✅ 并发设置选项
- ✅ 防抖机制竞态
- ✅ 并发执行和修改钩子
- ✅ 并发配置重载

## 技术细节

### 锁的层次结构

```
optsMutex (选项锁)
    ↓
hookMutex (钩子锁)
    ↓
rwMutex (配置锁)
```

**规则**: 始终按照上述顺序获取锁，避免死锁

### 原子操作

```go
// 防抖时间戳使用 atomic.Int64
m.lastChangeNano.Store(now.UnixNano())  // 写入
lastNano := m.lastChangeNano.Load()     // 读取
```

### 钩子执行

```go
// 在锁外执行钩子函数，避免死锁
m.hookMutex.RLock()
handler := m.hooks.Handles[pattern]
m.hookMutex.RUnlock()

if handler != nil {
    handler(ctx)  // 在锁外执行
}
```

## 相关文档

- `CONCURRENCY_ANALYSIS.md` - 详细的并发安全分析
- `CONCURRENCY_FIX_SUMMARY.md` - 修复总结
- `VERIFICATION_REPORT.md` - 验证报告
- `concurrency_test.go` - 并发测试用例

## 保证

✅ **无数据竞态**: 所有共享数据都有适当的同步保护  
✅ **无死锁**: 锁的获取顺序一致，不在锁内调用外部代码  
✅ **无活锁**: 没有重试循环或自旋等待  
✅ **无饥饿**: 使用标准库的公平锁  

## 版本信息

- **修复版本**: v2.0.0+
- **修复日期**: 2025-11-10
- **测试状态**: ✅ 所有测试通过

---

**结论**: ConfigX v2 可以安全地在高并发环境中使用，所有操作都是线程安全的。
