# 新功能说明

## 版本: v2.0.0+
## 更新日期: 2025-11-10

---

## 🎉 新增功能

### 1. 多种配置文件格式支持

ConfigX 现在支持多种主流配置文件格式，自动根据文件扩展名识别：

#### 支持的格式

| 格式 | 扩展名 | 说明 |
|------|--------|------|
| **YAML** | `.yaml`, `.yml` | 推荐格式，可读性好 |
| **JSON** | `.json` | 标准格式，易于生成 |
| **TOML** | `.toml` | 配置文件专用格式 |
| **HCL** | `.hcl` | HashiCorp 配置语言 |
| **INI** | `.ini` | 传统配置格式 |
| **Properties** | `.properties`, `.props`, `.prop` | Java 风格 |

#### 使用示例

```go
// YAML 格式
opts.Filename.Set("config.yaml")

// JSON 格式
opts.Filename.Set("config.json")

// TOML 格式
opts.Filename.Set("config.toml")

// 无需额外配置，自动识别格式
manager.LoadConfig()
```

#### 示例程序

查看 `example/multi-format/` 目录获取完整示例。

---

### 2. 环境变量覆盖支持

ConfigX 现在支持使用环境变量覆盖配置文件中的值，这对于以下场景非常有用：
- 保护敏感信息（API 密钥、密码）
- 不同环境使用不同配置（开发、测试、生产）
- 容器化部署（Docker, Kubernetes）
- 符合 12-Factor App 原则

#### 方式 1: 自动环境变量（推荐用于开发）

```go
opts := configx.NewOption()
opts.SetEnvPrefix("MYAPP")                    // 设置前缀
opts.EnableAutomaticEnv(true)                 // 启用自动读取
opts.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
manager.SetOption(opts)
```

配置项 `database.password` 会自动从环境变量 `MYAPP_DATABASE_PASSWORD` 读取。

#### 方式 2: 精确绑定（推荐用于生产）

```go
// 只绑定需要的敏感配置项
manager.BindEnv("database.password", "DB_PASSWORD")
manager.BindEnv("api.key", "API_KEY")
manager.BindEnv("jwt.secret", "JWT_SECRET")
```

这种方式更安全可控，只有明确绑定的配置项才会从环境变量读取。

#### 环境变量命名规则

1. **使用前缀**: `PREFIX_SECTION_KEY`
   - 配置项: `database.password`
   - 环境变量: `MYAPP_DATABASE_PASSWORD`

2. **嵌套结构**: 使用下划线分隔
   - 配置项: `server.database.host`
   - 环境变量: `MYAPP_SERVER_DATABASE_HOST`

3. **大小写**: 环境变量通常使用大写

#### 优先级

环境变量的优先级高于配置文件：

```
环境变量 > 配置文件 > 默认值
```

#### 示例程序

- `example/env-bind/` - 精确绑定示例（推荐）
- `example/env-override/` - 自动环境变量示例

---

## 🔒 安全最佳实践

### 敏感信息管理

#### ❌ 不推荐：将敏感信息写入配置文件

```yaml
# config.yaml - 不安全
database:
  password: "my_real_password_123"
api:
  key: "sk_live_real_api_key"
```

#### ✅ 推荐：使用环境变量

```yaml
# config.yaml - 安全
database:
  password: "placeholder"  # 将被环境变量覆盖
api:
  key: "placeholder"       # 将被环境变量覆盖
```

```go
// 代码中绑定环境变量
manager.BindEnv("database.password", "DB_PASSWORD")
manager.BindEnv("api.key", "API_KEY")
```

```bash
# 在部署环境中设置
export DB_PASSWORD="my_real_password_123"
export API_KEY="sk_live_real_api_key"
```

### 生产环境配置

```go
// 生产环境推荐配置
manager := configx.NewManager(AppConfig{})

// 1. 加载基础配置文件
opts := configx.NewOption()
opts.Filename.Set("config.yaml")
manager.SetOption(opts)

// 2. 绑定敏感配置到环境变量
manager.BindEnv("database.password", "DB_PASSWORD")
manager.BindEnv("redis.password", "REDIS_PASSWORD")
manager.BindEnv("jwt.secret", "JWT_SECRET")
manager.BindEnv("aws.access_key", "AWS_ACCESS_KEY_ID")
manager.BindEnv("aws.secret_key", "AWS_SECRET_ACCESS_KEY")

// 3. 加载配置（环境变量会自动覆盖）
if err := manager.LoadConfig(); err != nil {
    log.Fatal(err)
}
```

---

## 📦 容器化部署

### Docker 示例

```dockerfile
# Dockerfile
FROM golang:1.21-alpine

WORKDIR /app
COPY . .
RUN go build -o myapp

# 配置文件
COPY config.yaml /app/config.yaml

# 运行时通过环境变量传递敏感信息
CMD ["./myapp"]
```

```bash
# 运行容器时传递环境变量
docker run -e DB_PASSWORD=secret123 \
           -e API_KEY=sk_live_xyz \
           myapp
```

### Kubernetes 示例

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
spec:
  template:
    spec:
      containers:
      - name: myapp
        image: myapp:latest
        env:
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: myapp-secrets
              key: db-password
        - name: API_KEY
          valueFrom:
            secretKeyRef:
              name: myapp-secrets
              key: api-key
```

---

## 🔄 迁移指南

### 从 v1 迁移到 v2

如果你正在使用 ConfigX v1，以下是主要变更：

#### 1. 模块路径更新

```go
// v1
import "github.com/kawaiirei0/configx"

// v2
import "github.com/kawaiirei0/configx/v2"
```

#### 2. 新增环境变量支持

```go
// v2 新增功能
opts.SetEnvPrefix("MYAPP")
opts.EnableAutomaticEnv(true)
manager.BindEnv("api.key", "API_KEY")
```

#### 3. 多格式支持

```go
// v2 自动识别格式
opts.Filename.Set("config.json")  // JSON
opts.Filename.Set("config.toml")  // TOML
opts.Filename.Set("config.yaml")  // YAML
```

---

## 📊 性能说明

### 环境变量性能

- 环境变量读取是在配置加载时进行的
- 不会影响 `GetConfig()` 的性能
- 环境变量值会被缓存在 Viper 中

### 多格式支持性能

- 格式识别基于文件扩展名，无性能开销
- 不同格式的解析性能差异：
  - JSON: 最快
  - YAML: 中等
  - TOML: 中等
  - INI: 快

---

## 🧪 测试

### 运行示例

```bash
# 多格式支持
go run example/multi-format/main.go

# 环境变量绑定
go run example/env-bind/main.go

# 环境变量覆盖
go run example/env-override/main.go
```

### 单元测试

```bash
# 运行所有测试
go test -v ./...

# 运行并发测试
go test -v -run Concurrent ./...
```

---

## 📚 相关文档

- `README.md` - 主文档
- `CONCURRENCY_SAFETY.md` - 并发安全说明
- `example/README.md` - 示例说明
- `VERIFICATION_REPORT.md` - 验证报告

---

## 🎯 总结

ConfigX v2 新增的功能使其成为一个功能完整、安全可靠的配置管理库：

✅ **多格式支持** - 支持 6+ 种配置文件格式  
✅ **环境变量** - 灵活的环境变量覆盖机制  
✅ **安全性** - 保护敏感信息的最佳实践  
✅ **容器友好** - 完美支持 Docker/Kubernetes  
✅ **向后兼容** - 平滑升级，无破坏性变更  
✅ **线程安全** - 完全的并发安全保证  

开始使用这些新功能，让你的应用配置管理更加灵活和安全！
