package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// 在 Manager.Init() 中增加判断
func (m *Manager) ensureConfigFile() error {
	absPath, _ := filepath.Abs(m.opt.Path.ToValue())
	if err := os.MkdirAll(absPath, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	cfgFile := filepath.Join(
		absPath,
		fmt.Sprintf("%s.%s.%s", m.opt.Filename.ToValue(), m.opt.Env.ToValue(), m.opt.FileType.ToValue()),
	)

	_, err := os.Stat(cfgFile)
	if os.IsNotExist(err) {
		// 文件不存在，写入默认配置
		defaultCfg := NewConfig() // 自定义函数，返回 *Config 带默认值

		data, err := yaml.Marshal(defaultCfg)
		if err != nil {
			return fmt.Errorf("failed to marshal default config: %w", err)
		}

		if err := os.WriteFile(cfgFile, data, 0644); err != nil {
			return fmt.Errorf("failed to write default config file: %w", err)
		}

		m.logger.Info.Exec(fmt.Sprintf("[config] 默认配置文件已生成: %s", cfgFile))
	}

	m.vp.SetConfigFile(cfgFile)
	return nil
}
