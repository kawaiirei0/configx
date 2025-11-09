package configx

import (
	"fmt"
)

// SetOption 设置配置选项
// 参数：
//   opts: 配置选项，如果为 nil 则使用默认选项
// 返回值：
//   *Manager[T]: 返回管理器实例以支持链式调用
func (m *Manager[T]) SetOption(opts *Option) *Manager[T] {
	if !m.optsInit {
		// 标记已初始化option
		m.optsInit = true
		// 初始化设置选项默认值
		if opts != nil {
			opts.setDefaultValue()
		} else {
			opts = NewOption()
		}
		m.opts = opts
	}
	return m
}

// Init 初始化配置管理器并启动热重载监控
// 参数：
//   handles: 配置变更时的回调函数列表
// 返回值：
//   error: 初始化失败时返回错误
func (m *Manager[T]) Init(handles ...HandlerFunc) error {
	// 如果option不存在配置则设置默认选项
	m.SetOption(nil)

	// hook init
	m.hooks.Handles[InitHook].Exec(HookContext{
		Message: "开始初始化",
	})

	// setting debouncedur
	m.debounceDur = m.opts.DebounceDur.ToValue()

	inFile := m.opts.File()

	m.vp.SetConfigFile(inFile)

	// 如果文件不存在，则创建默认配置文件
	if err := m.ensureConfigFile(m.opts); err != nil {
		m.hooks.Handles[Error].Exec(HookContext{
			Message: fmt.Sprintf("[config] 创建默认配置文件失败: %v", err),
		})
		return err
	}

	// 读取配置文件
	if err := m.vp.ReadInConfig(); err != nil {
		m.hooks.Handles[Error].Exec(HookContext{
			Message: fmt.Sprintf("[config] 加载配置失败: %v", err),
		})
		return err
	}

	m.hooks.Handles[Info].Exec(HookContext{
		Message: fmt.Sprintf("[config] 已加载配置文件: %s", m.vp.ConfigFileUsed()),
	})

	// 解析配置到结构体
	if err := m.Unmarshal(); err != nil {
		m.hooks.Handles[Error].Exec(HookContext{
			Message: fmt.Sprintf("[config] 解析配置到结构体失败 Error: %s", err.Error()),
		})
		return err
	}

	// 监听配置变更，传递回调函数
	m.monitorConfigChanges(handles)

	// 验证配置通过
	m.validateConfig(true)

	return nil
}
