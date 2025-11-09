package configx

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Manager[T any] struct {
	config              *T            // 泛型配置对象
	vp                  *viper.Viper  // Viper 实例
	rwMutex             sync.RWMutex  // 读写锁
	lastChange          time.Time     // 上次触发时间（用于防抖）
	debounceDur         time.Duration // 防抖间隔
	hooks               *Hook         // hook
	pathName            string        // 配置文件
	opts                *Option       // 设置选项
	optsInit            bool          // 初始化选项
	validateConfigValue bool          // 验证
	defaultConfig       any           // default config
}

// Note: Global singleton removed due to Go generics limitations
// Users should manage Manager instances in their own code

// NewManager 创建泛型配置管理器
// 参数 defaultConfig: 配置结构体的零值或默认值
// 返回值：
//
//	*Manager[T]: 泛型管理器实例
func NewManager[T any](defaultConfig T) *Manager[T] {
	return &Manager[T]{
		config:        nil, // 配置将在 LoadConfig 时初始化
		vp:            viper.New(),
		lastChange:    time.Time{},
		hooks:         NewHook(),
		defaultConfig: defaultConfig,
	}
}

// GetConfig 获取配置副本（类型安全）
// 使用读锁保护并发访问，返回配置的深拷贝
// 如果配置类型实现了 Cloneable[T] 接口，将使用自定义的 Clone() 方法
// 否则使用 JSON 序列化/反序列化实现深拷贝
// 返回值：
//
//	T: 配置副本
//	error: 如果配置未初始化则返回错误
func (m *Manager[T]) GetConfig() (T, error) {
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()

	var zero T
	if m.config == nil {
		return zero, ErrConfigNotInitialized
	}

	// 如果配置类型实现了 Cloneable 接口，使用自定义克隆
	if cloneable, ok := any(*m.config).(Cloneable[T]); ok {
		return cloneable.Clone(), nil
	}

	// 否则使用 JSON 序列化深拷贝
	return m.jsonDeepCopy()
}

// LoadConfig 加载配置文件
// 使用写锁保护配置更新，读取并解析配置文件到泛型类型
// 返回值：
//
//	error: 如果读取或解析失败则返回详细错误信息
func (m *Manager[T]) LoadConfig() error {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()

	// 配置 Viper
	if err := m.setupViper(); err != nil {
		return fmt.Errorf("配置 Viper 失败: %w", err)
	}

	// 读取配置文件
	if err := m.vp.ReadInConfig(); err != nil {
		return fmt.Errorf("%w: %s, 错误: %v", ErrConfigFileNotFound, m.vp.ConfigFileUsed(), err)
	}

	// 创建新的泛型配置实例
	var newConfig T

	// 解析配置到泛型类型
	if err := m.vp.Unmarshal(&newConfig); err != nil {
		return fmt.Errorf("%w: 文件 %s, 错误: %v", ErrConfigParseFailed, m.vp.ConfigFileUsed(), err)
	}

	// 更新配置
	m.config = &newConfig

	return nil
}

// jsonDeepCopy 使用 JSON 序列化/反序列化实现深拷贝
// 这是默认的深拷贝方法，当配置类型未实现 Cloneable 接口时使用
// 返回值：
//
//	T: 深拷贝后的配置对象
//	error: 序列化或反序列化失败时返回错误
func (m *Manager[T]) jsonDeepCopy() (T, error) {
	var zero T
	
	// 序列化
	data, err := json.Marshal(*m.config)
	if err != nil {
		return zero, fmt.Errorf("序列化配置失败: %w", err)
	}
	
	// 反序列化
	var copy T
	if err := json.Unmarshal(data, &copy); err != nil {
		return zero, fmt.Errorf("反序列化配置失败: %w", err)
	}
	
	return copy, nil
}

// setupViper 配置 Viper 实例
// 设置配置文件路径、文件名和文件类型
// 此方法与泛型 Manager[T] 完全兼容，保持文件路径、文件名、文件类型的配置逻辑
// 返回值：
//
//	error: 如果配置失败则返回错误
func (m *Manager[T]) setupViper() error {
	// 确保选项已初始化
	if !m.optsInit {
		m.SetOption(nil)
	}

	// 设置配置文件路径
	inFile := m.opts.File()
	m.vp.SetConfigFile(inFile)

	return nil
}

// SetHook 设置钩子处理函数
// 参数：
//   pattern: 钩子级别（Debug, Info, Warn, Error）
//   handler: 钩子处理函数
// 返回值：
//   *Manager[T]: 返回管理器实例以支持链式调用
// 功能：
//   - 支持设置不同级别的钩子
//   - 保持现有的钩子级别（Debug, Info, Warn, Error）
//   - 与泛型 Manager[T] 完全兼容
func (m *Manager[T]) SetHook(pattern HookPattern, handler HookHandlerFunc) *Manager[T] {
	m.hooks.SetHook(pattern, handler)
	return m
}
