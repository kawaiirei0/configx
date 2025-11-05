package config

import (
	"path/filepath"
)

// LoadConfigPath 返回配置文件路径
func LoadConfigPath(p ...string) string {
	path, _ := filepath.Abs(p[0])
	return path
}
