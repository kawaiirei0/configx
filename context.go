package configx

import "github.com/fsnotify/fsnotify"

// HandlerFunc 配置变更回调函数类型
type HandlerFunc func(ctx *Context)

// Context 配置变更回调上下文
// 提供配置变更事件信息和访问管理器的能力
type Context struct {
	// FSEvent 文件系统变更事件
	FSEvent fsnotify.Event
	// manager 存储管理器实例（类型为 interface{} 以支持泛型）
	// 使用 GetManager[T]() 方法获取类型安全的管理器实例
	manager interface{}
}

// GetManager 获取管理器实例（需要手动类型断言）
// 返回值：
//   interface{}: 管理器实例，需要调用者进行类型断言
// 示例：
//   if mgr, ok := ctx.GetManager().(*configx.Manager[MyConfig]); ok {
//       config, _ := mgr.GetConfig()
//       // 使用 config
//   }
func (c *Context) GetManager() interface{} {
	return c.manager
}
