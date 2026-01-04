package configs

import (
	"os"
)

// parseEnvs 解析环境变量并更新配置
// 该函数会遍历配置结构体的所有字段，查找带有 "env" 标签的字段，
// 并从对应的环境变量中读取值来更新配置。
// 支持的字段类型包括：string、bool、int、int32、int64、float32、float64
// 如果环境变量不存在或解析失败，会记录警告日志，但不会中断程序执行
func (m *Manager[T]) parseEnvs() {
	t := m.confElem.Type()

	for i := 0; i < t.NumField(); i++ {
		fieldKey := t.Field(i)
		envName := fieldKey.Tag.Get("env")
		if envName == "" {
			continue
		}
		envValue := os.Getenv(envName)
		if envValue == "" {
			return
		}
		conf, err := m.Conf(fieldKey.Name)
		if err != nil {
			continue
		}
		err = conf.SetByString(envValue)
		if err != nil {
			continue
		}
	}
}
