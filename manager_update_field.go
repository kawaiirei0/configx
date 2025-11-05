package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// UpdateField 更新配置字段
// 参数：
//
//	updateFunc: 更新函数
//
// 返回值：
//
//	error: 更新过程中的错误
func UpdateField(updateFunc func(*Config)) error {
	m := Default()
	return m.UpdateField(updateFunc)
}

// UpdateField 更新配置字段
// 参数：
//
//	updateFunc: 更新函数
//
// 返回值：
//
//	error: 更新过程中的错误
func (m *Manager) UpdateField(updateFunc func(*Config)) error {
	m.rwMutex.Lock()
	defer m.rwMutex.Unlock()

	oldConfig := *m.config
	updateFunc(m.config)

	configFile := m.vp.ConfigFileUsed()
	content, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	newContent := string(content)

	var updateContent func(reflect.Value, reflect.Value, reflect.Type)
	updateContent = func(oldVal, newVal reflect.Value, t reflect.Type) {
		for i := 0; i < oldVal.NumField(); i++ {
			oldField, newField := oldVal.Field(i), newVal.Field(i)
			if tag := t.Field(i).Tag.Get("mapstructure"); tag != "" {
				if oldField.Kind() == reflect.Struct {
					updateContent(oldField, newField, oldField.Type())
				} else if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
					var old, new string
					if oldField.Kind() == reflect.Slice || oldField.Kind() == reflect.Array {
						// 数组类型
						var oldElems, newElems []string
						for i := 0; i < oldField.Len(); i++ {
							elem := oldField.Index(i)
							if elem.Kind() == reflect.String {
								oldElems = append(oldElems, fmt.Sprintf(`"%s"`, elem.String()))
							} else {
								oldElems = append(oldElems, fmt.Sprintf("%v", elem.Interface()))
							}
						}
						for i := 0; i < newField.Len(); i++ {
							elem := newField.Index(i)
							if elem.Kind() == reflect.String {
								newElems = append(newElems, fmt.Sprintf(`"%s"`, elem.String()))
							} else {
								newElems = append(newElems, fmt.Sprintf("%v", elem.Interface()))
							}
						}
						old, new = fmt.Sprintf("[%s]", strings.Join(oldElems, ", ")), fmt.Sprintf("[%s]", strings.Join(newElems, ", "))

						for _, pattern := range []string{fmt.Sprintf(`%s: %s`, tag, old), fmt.Sprintf(`%s: []`, tag)} {
							if strings.Contains(newContent, pattern) {
								newContent = strings.ReplaceAll(newContent, pattern, fmt.Sprintf(`%s: %s`, tag, new))
								break
							}
						}
					} else {
						// 非数组类型
						old, new = fmt.Sprintf("%v", oldField.Interface()), fmt.Sprintf("%v", newField.Interface())
						for _, pattern := range []string{
							fmt.Sprintf(`%s: "%s"`, tag, old),
							fmt.Sprintf(`%s: %s`, tag, old),
							fmt.Sprintf(`%s: ""`, tag),
						} {
							if strings.Contains(newContent, pattern) {
								newContent = strings.ReplaceAll(newContent, pattern, fmt.Sprintf(`%s: "%s"`, tag, new))
								break
							}
						}
					}
				}
			}
		}
	}

	updateContent(reflect.ValueOf(oldConfig), reflect.ValueOf(*m.config), reflect.TypeOf(oldConfig))

	if newContent != string(content) {
		return os.WriteFile(configFile, []byte(newContent), 0644)
	}

	return nil
}
