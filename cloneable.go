package configx

// Cloneable 可克隆接口
// 配置类型可以实现此接口以提供自定义的深拷贝逻辑
// 这比默认的 JSON 序列化方式更高效，并且可以处理特殊字段
//
// 性能优势：
//   - 避免 JSON 序列化/反序列化的开销
//   - 可以精确控制哪些字段需要深拷贝
//   - 可以处理 JSON 不支持的类型（如 channels, functions）
//
// 示例实现：
//
//	type MyConfig struct {
//	    Name    string
//	    Version string
//	    Data    map[string]interface{}
//	}
//
//	// 实现 Cloneable 接口
//	func (c MyConfig) Clone() MyConfig {
//	    // 实现自定义克隆逻辑
//	    clone := MyConfig{
//	        Name:    c.Name,
//	        Version: c.Version,
//	        Data:    make(map[string]interface{}),
//	    }
//	    // 深拷贝 map
//	    for k, v := range c.Data {
//	        clone.Data[k] = v
//	    }
//	    return clone
//	}
//
// 使用示例：
//
//	// 创建管理器
//	manager := configx.NewManager(MyConfig{})
//	manager.LoadConfig()
//
//	// GetConfig 会自动检测并使用 Clone() 方法
//	config, err := manager.GetConfig()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// 注意事项：
//   - Clone() 方法必须返回值类型（不是指针）
//   - 确保所有需要深拷贝的字段都被正确处理
//   - 如果不实现此接口，GetConfig() 会自动使用 JSON 深拷贝
type Cloneable[T any] interface {
	Clone() T
}
