package configs

import (
	"os"
	"reflect"
	"strconv"

	"github.com/wolfbolin/bolbox/pkg/log"
)

// parseEnvs 解析环境变量并更新配置
// 该函数会遍历配置结构体的所有字段，查找带有 "env" 标签的字段，
// 并从对应的环境变量中读取值来更新配置。
// 支持的字段类型包括：string、bool、int、int32、int64、float32、float64
// 如果环境变量不存在或解析失败，会记录警告日志，但不会中断程序执行
func (m *Manager[T]) parseEnvs() {
	v := reflect.ValueOf(m.userConf).Elem()
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		fieldKey := t.Field(i)
		fieldVal := v.Field(i)
		if !fieldVal.CanSet() {
			continue
		}
		parseEnv(fieldKey, fieldVal)
	}
}

// parseEnv 从环境变量中解析单个配置字段的值
// 参数说明：
//   - fieldKey: 结构体字段的元信息，包含 "env" 标签，用于指定环境变量名
//   - fieldVal: 字段的值，用于接收解析后的环境变量值
//
// 该函数会从字段的 "env" 标签中获取环境变量名，然后从系统环境变量中读取值。
// 如果环境变量不存在或为空，函数会直接返回，不会更新配置值。
// 对于数值类型（int、float），会进行字符串到数值的转换，转换失败时会记录警告日志。
// 支持的字段类型包括：string、bool、int、int32、int64、float32、float64
func parseEnv(fieldKey reflect.StructField, fieldVal reflect.Value) {
	envName := fieldKey.Tag.Get("env")
	if envName == "" {
		return
	}

	envValue := os.Getenv(envName)
	if envValue == "" {
		return
	}

	switch fieldVal.Kind() {
	case reflect.String:
		fieldVal.SetString(envValue)
	case reflect.Bool:
		boolVal, convErr := strconv.ParseBool(envValue)
		if convErr == nil {
			fieldVal.SetBool(boolVal)
		} else {
			log.Warnf("Parse env[%s] to conf[%s](Bool) failed. %s", envName, fieldKey.Name, convErr.Error())
		}
	case reflect.Int, reflect.Int32, reflect.Int64:
		intVal, convErr := strconv.ParseInt(envValue, 10, 64)
		if convErr == nil {
			fieldVal.SetInt(intVal)
		} else {
			log.Warnf("Parse env[%s] to conf[%s](Int) failed. %s", envName, fieldKey.Name, convErr.Error())
		}
	case reflect.Float32, reflect.Float64:
		floatVal, convErr := strconv.ParseFloat(envValue, 64)
		if convErr == nil {
			fieldVal.SetFloat(floatVal)
		} else {
			log.Warnf("Parse env[%s] to conf[%s](Float) failed. %s", envName, fieldKey.Name, convErr.Error())
		}
	default:
		log.Warnf("Not suppose parse process env[%s] for var[%s]", envName, fieldKey.Name)
	}
}
