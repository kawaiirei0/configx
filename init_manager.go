package config

import (
	"fmt"
	"path/filepath"
)

func (m *Manager) Init(managerLogger *Logger, opt *Option, handles ...HandlerFunc) error {
	// 初始化设置选项默认值
	if opt != nil {
		opt.defaultValueInit()
	} else {
		opt = NewOption().defaultValueInit()
	}

	// 拿到设置选项
	m.opt = opt

	// 管理器拿到日志实例
	m.logger = managerLogger

	// hook init
	m.logger.Debug.Exec("[config] 开始初始化")

	// 设置防抖
	m.debounceDur = m.opt.DebounceDur.ToValue()

	absPath, _ := filepath.Abs(m.opt.Path.ToValue())

	m.vp.SetConfigType(m.opt.FileType.ToValue()) // 设置文件类型
	// 根据环境加载不同配置文件 设置文件名
	if m.opt.Env != "" {
		// 设置文件名.环境名
		m.vp.SetConfigName(fmt.Sprintf("%s.%s", m.opt.Filename, m.opt.Env))
	} else {
		// 设置文件名
		m.vp.SetConfigName(m.opt.Filename.ToValue())
	}
	m.vp.AddConfigPath(absPath) // 设置文件路径

	// 如果文件不存在，则创建默认配置文件
	if err := m.ensureConfigFile(); err != nil {
		m.logger.Error.Exec(fmt.Sprintf("[config] 创建默认配置文件失败: %v", err))
		return err
	}

	// 读取配置文件
	if err := m.vp.ReadInConfig(); err != nil {
		m.logger.Error.Exec(fmt.Sprintf("[config] 加载配置失败: %v", err))
		return err
	}

	m.logger.Debug.Exec(fmt.Sprintf("[config] 已加载配置文件: %s", m.vp.ConfigFileUsed()))

	// 解析配置到结构体
	if err := m.Unmarshal(); err != nil {
		m.logger.Error.Exec(fmt.Sprintf("[config] 解析配置到结构体Error: %s", err.Error()))
		return err
	}

	// 监听配置变更
	m.monitorConfigChanges(handles)

	// 验证配置通过
	m.validateConfig(true)

	return nil
}
