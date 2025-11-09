package configx

func (m *Manager[T]) validateConfig(ok ...bool) bool {
	if len(ok) != 0 {
		m.validateConfigValue = ok[0]
		return m.validateConfigValue
	}
	return m.validateConfigValue
}
