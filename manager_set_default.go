package configx

func (m *Manager[T]) SetDefault(key string, value any) {
	m.vp.SetDefault(key, value)
}
