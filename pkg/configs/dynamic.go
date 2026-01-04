package configs

import (
	"reflect"
	"strconv"

	"github.com/wolfbolin/bolbox/pkg/log"
)

// ParseMap 从 map 中解析配置值并更新配置（线程安全）
// 参数说明：
//   - data: 键值对映射，键为配置结构体的字段名（区分大小写），值为字符串格式的配置值
//
// 该函数会遍历 map 中的每个键值对，尝试将值解析并设置到对应的配置字段中。
// map 的键必须与配置结构体的字段名完全匹配（区分大小写）。
// 支持的字段类型包括：string、bool、int、int32、int64、float32、float64
//
// 行为说明：
//   - 如果 map 为 nil，函数直接返回，不进行任何操作
//   - 如果 map 中的键在配置结构体中不存在，会记录警告日志并跳过该键
//   - 如果字段不可设置（如未导出的字段），会记录警告日志并跳过
//   - 对于数值类型，会进行字符串到数值的转换，转换失败时会记录警告日志但不中断执行
//
// 使用示例：
//   mgr.ParseMap(map[string]string{
//       "Port":    "9090",
//       "Debug":   "true",
//       "LogFile": "/var/log/app.log",
//   })
func (m *Manager[T]) ParseMap(data map[string]string) {
	if data == nil {
		return
	}
	m.confLock.Lock()
	defer m.confLock.Unlock()

	v := reflect.ValueOf(m.userConf).Elem()
	for mapName, mapValue := range data {
		fieldVal := v.FieldByName(mapName)
		if !fieldVal.IsValid() {
			log.Warnf("Can not found map key[%s] in config. Skip it.", mapName)
			continue
		}
		if !fieldVal.CanSet() {
			log.Warnf("Field[%s] is not settable. Skip it.", mapName)
			continue
		}

		switch fieldVal.Kind() {
		case reflect.String:
			fieldVal.SetString(mapValue)
		case reflect.Bool:
			boolVal, convErr := strconv.ParseBool(mapValue)
			if convErr == nil {
				fieldVal.SetBool(boolVal)
			} else {
				log.Warnf("Parse map[%s] to conf[%s](Bool) failed. %s", mapName, fieldVal.Type().Name(), convErr.Error())
			}
		case reflect.Int, reflect.Int32, reflect.Int64:
			intVal, convErr := strconv.ParseInt(mapValue, 10, 64)
			if convErr == nil {
				fieldVal.SetInt(intVal)
			} else {
				log.Warnf("Parse map[%s] to conf[%s](Int) failed. %s", mapName, fieldVal.Type().Name(), convErr.Error())
			}
		case reflect.Float32, reflect.Float64:
			floatVal, convErr := strconv.ParseFloat(mapValue, 64)
			if convErr == nil {
				fieldVal.SetFloat(floatVal)
			} else {
				log.Warnf("Parse map[%s] to conf[%s](Float) failed. %s", mapName, fieldVal.Type().Name(), convErr.Error())
			}
		default:
			log.Warnf("Not suppose parse process map[%s](%s) for config", mapName, fieldVal.Type().Name())
		}
	}
}
