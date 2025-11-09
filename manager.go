package configx

import (
	"github.com/kawaiirei0/configx/utils"
	"sync"
	"time"

	"github.com/spf13/viper"
)

type Manager struct {
	config              *Config       // 全局配置对象
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

var globalManager = utils.NewSingleton[Manager]()

// Default 获取配置管理器
// 返回值：
//
//	*Manager: 管理器
func Default() *Manager {
	v, _ := globalManager.Get(func() (*Manager, error) {
		return NewManager(), nil
	})
	return v
}

func NewManager() *Manager {
	return &Manager{
		config:     NewConfig(),
		vp:         viper.New(),
		lastChange: time.Time{},
		// opts:       NewOption(),
		hooks: NewHook(),
	}
}
