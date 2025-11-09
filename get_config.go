package configx

import "fmt"

// GetConfig 获取配置
// 返回值：
//
//	*Config: 配置副本
//	error: 获取过程中的错误
func GetConfig() (*Config, error) {
	m := Default()
	m.rwMutex.RLock()
	defer m.rwMutex.RUnlock()

	if m.config == nil {
		return nil, fmt.Errorf("config not initialized")
	}

	configCopy := *m.config
	return &configCopy, nil
}
