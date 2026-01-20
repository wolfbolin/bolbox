package configs

// ParseMap 从map中解析并更新配置值，键为字段名（区分大小写）
func (m *Manager[T]) ParseMap(data map[string]string) {
	if data == nil {
		return
	}
	m.confLock.Lock()
	defer m.confLock.Unlock()

	for mapName, mapValue := range data {
		conf, err := m.Conf(mapName)
		if err != nil {
			continue
		}
		err = conf.SetByString(mapValue)
		if err != nil {
			continue
		}
	}
}
