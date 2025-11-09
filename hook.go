package configx

// HookPattern 钩子级别类型
type HookPattern int

const (
	// InitHook 初始化钩子
	InitHook HookPattern = iota
	// Debug 调试级别钩子
	Debug
	// Info 信息级别钩子
	Info
	// Warn 警告级别钩子
	Warn
	// Error 错误级别钩子
	Error
	// HookIndex 钩子数组大小
	HookIndex
)

// HookContext 钩子上下文
// 包含钩子触发时的消息和级别信息
type HookContext struct {
	Message string
	Pattern HookPattern
}

// HookHandlerFunc 钩子处理函数类型
type HookHandlerFunc func(ctx HookContext)

// Exec 执行钩子处理函数
// 如果处理函数为 nil，则不执行
func (h HookHandlerFunc) Exec(ctx HookContext) {
	if h == nil {
		return
	}
	h(ctx)
}

// Hook 钩子管理器
// 管理不同级别的钩子处理函数
// 与泛型 Manager[T] 完全兼容
type Hook struct {
	Handles [HookIndex]HookHandlerFunc
}

// NewHook 创建新的钩子管理器
func NewHook() *Hook {
	return &Hook{}
}

// SetHook 设置指定级别的钩子处理函数
// 参数：
//   index: 钩子级别（InitHook, Debug, Info, Warn, Error）
//   h: 钩子处理函数
// 返回值：
//   *Hook: 返回钩子管理器实例以支持链式调用
func (hooks *Hook) SetHook(index HookPattern, h HookHandlerFunc) *Hook {
	hooks.Handles[index] = h
	return hooks
}

// SetHook is deprecated - use Manager.hooks.SetHook instead
// Global singleton removed due to Go generics limitations
// func SetHook(index HookPattern, h HookHandlerFunc) *Hook {
// 	hooks := Default().hooks
// 	hooks.Handles[index] = h
// 	return hooks
// }
