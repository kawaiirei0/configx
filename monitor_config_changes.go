package configx

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Unmarshal 解析配置到结构体
func (m *Manager[T]) Unmarshal() error {
	var newConfig T
	if err := m.vp.Unmarshal(&newConfig); err != nil {
		m.hooks.Handles[Error].Exec(HookContext{
			Message: fmt.Sprintf("failed to unmarshal new config: %v", err),
		})
		return errors.New(fmt.Sprintf("failed to unmarshal new config: %v", err))
	}

	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()

	if m.config != nil {
		oldConfig := *m.config
		changes := make(map[string][2]any)

		if !compareStructs(oldConfig, newConfig, "", changes) {
			m.hooks.Handles[Error].Exec(HookContext{
				Message: "config type mismatch, changes blocked",
			})
			return errors.New(fmt.Sprintf("config type mismatch, changes blocked"))
		}
	}

	m.config = &newConfig
	return nil
}

// monitorConfigChanges 监听配置变更（带防抖与类型过滤）
// 参数：
//   handles: 配置变更时的回调函数列表
// 功能：
//   - 使用防抖机制避免频繁重载
//   - 在配置变更时调用泛型 LoadConfig 重新加载
//   - 触发钩子记录配置变更事件
//   - 执行开发者提供的回调函数
//   - 确保重载失败时保持原有配置不变
func (m *Manager[T]) monitorConfigChanges(handles []HandlerFunc) {
	m.vp.WatchConfig()
	m.vp.OnConfigChange(func(e fsnotify.Event) {
		// 仅响应写入事件，忽略 CHMOD/RENAME 等
		if e.Op != fsnotify.Write {
			return
		}

		// 防抖处理：忽略短时间内的重复变更
		now := time.Now()
		if now.Sub(m.lastChange) < m.debounceDur {
			return
		}
		m.lastChange = now

		// 触发钩子：检测到配置文件变更
		m.hooks.Handles[Info].Exec(HookContext{
			Message: fmt.Sprintf("[config] 检测到文件变更: %s", e.Name),
		})

		// 保存当前配置的副本，以便在重载失败时恢复
		m.rwMutex.RLock()
		oldConfig := m.config
		m.rwMutex.RUnlock()

		// 重新加载配置文件
		if err := m.vp.ReadInConfig(); err != nil {
			m.hooks.Handles[Error].Exec(HookContext{
				Message: fmt.Sprintf("[config] 重新加载配置文件失败: %v", err),
			})
			return
		}

		// 解析配置到结构体
		if err := m.Unmarshal(); err != nil {
			// 解析失败，恢复原有配置
			m.rwMutex.Lock()
			m.config = oldConfig
			m.rwMutex.Unlock()
			
			m.hooks.Handles[Error].Exec(HookContext{
				Message: fmt.Sprintf("[config] 解析配置失败，保持原有配置: %v", err),
			})
			return
		}

		// 触发钩子：配置重新加载成功
		m.hooks.Handles[Info].Exec(HookContext{
			Message: "[config] 配置重新加载成功",
		})

		// 创建回调上下文，包含管理器实例引用
		ctx := &Context{
			FSEvent: e,
			manager: m,
		}

		// 执行开发者提供的回调函数
		for _, handle := range handles {
			handle(ctx)
		}
	})
}

// compareStructs 比较结构体并收集变更
// 参数：
//
//	oldObj: 旧结构体
//	newObj: 新结构体
//	prefix: 字段路径前缀
//	changes: 记录变更的映射
//
// 返回值：
//
//	bool: 结构体类型是否一致
func compareStructs(oldObj, newObj any, prefix string, changes map[string][2]any) bool {
	oldVal := reflect.ValueOf(oldObj)
	newVal := reflect.ValueOf(newObj)

	if oldVal.Type() != newVal.Type() {
		return false
	}

	if oldVal.Kind() != reflect.Struct {
		return true
	}

	for i := 0; i < oldVal.NumField(); i++ {
		oldField := oldVal.Field(i)
		newField := newVal.Field(i)
		fieldName := oldVal.Type().Field(i).Name
		fullName := prefix + fieldName

		if oldField.Kind() == reflect.Struct {
			if !compareStructs(oldField.Interface(), newField.Interface(), fullName+".", changes) {
				return false
			}
			continue
		}

		if oldField.Kind() != newField.Kind() {
			return false
		}

		if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
			changes[fullName] = [2]any{oldField.Interface(), newField.Interface()}
		}
	}

	return true
}
