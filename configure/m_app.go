package configure

// App 应用基础配置
type App struct {
	Name        string `yaml:"name" mapstructure:"name"`
	Version     string `yaml:"version" mapstructure:"version"`
	Description string `yaml:"description" mapstructure:"description"`
}
