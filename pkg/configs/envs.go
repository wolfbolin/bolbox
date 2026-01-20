package configs

import (
	"os"
)

// parseEnvs 解析环境变量并更新配置
func (m *Manager[T]) parseEnvs() error {
	t := m.confElem.Type()

	for i := 0; i < t.NumField(); i++ {
		fieldKey := t.Field(i)
		envName := fieldKey.Tag.Get("env")
		if envName == "" {
			continue
		}
		envValue := os.Getenv(envName)
		if envValue == "" {
			continue
		}
		conf, _ := m.Conf(fieldKey.Name)
		err := conf.SetByString(envValue)
		if err != nil {
			return err
		}
	}
	return nil
}
