package config

type HookType int

const (
	BeforeChange HookType = iota
	AfterChange
)

func (m *Manager) RegisterHook(t HookType, fn HandlerFunc) {
	
}
