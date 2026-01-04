package configs

import (
	"reflect"
	"strconv"
	"sync"

	"github.com/wolfbolin/bolbox/pkg/errors"
	"github.com/wolfbolin/bolbox/pkg/log"
)

// Manager 是一个泛型配置管理器，支持从环境变量、命令行参数和动态映射中加载配置
// 配置加载的优先级顺序为：默认值 < 环境变量 < 命令行参数 < 动态映射（ParseMap）
// 该管理器是线程安全的，可以在并发环境中安全使用
type Manager[T any] struct {
	userConf *T
	valueMap map[string]*Config
	confElem reflect.Value
	confLock sync.RWMutex
}

// NewManager 创建一个新的配置管理器实例
// 参数说明：
//   - def: 默认配置结构体指针，如果为 nil 则使用类型 T 的零值作为默认配置
//
// 创建管理器时会自动执行以下操作：
//  1. 使用默认配置初始化配置对象
//  2. 解析环境变量并更新配置（通过字段的 "env" 标签）
//  3. 解析命令行参数并更新配置（通过字段的 "flag" 标签）
//
// 使用示例：
//
//	type Config struct {
//	    Port    int    `env:"PORT" flag:"port" desc:"服务端口"`
//	    Debug   bool   `env:"DEBUG" flag:"debug" desc:"调试模式"`
//	    LogFile string `env:"LOG_FILE" flag:"log-file" desc:"日志文件路径"`
//	}
//	mgr := NewManager(&Config{Port: 8080, Debug: false})
func NewManager[T any](def *T) *Manager[T] {
	var userConfig *T
	if def != nil {
		userConfig = def
	} else {
		userConfig = new(T)
	}
	mgr := &Manager[T]{
		userConf: userConfig,
		valueMap: make(map[string]*Config),
		confElem: reflect.ValueOf(userConfig).Elem(),
	}

	// 建立反射对象索引
	t := mgr.confElem.Type()
	for i := range t.NumField() {
		key := t.Field(i)
		mgr.valueMap[key.Name] = &Config{
			key: key.Name,
			val: mgr.confElem.Field(i),
		}
	}

	mgr.parseEnvs()
	mgr.parseFlags()
	return mgr
}

// Vars 返回当前配置的只读副本（线程安全）
// 调用者的任何修改对全局配置中的数据不生效
// 如果需要更新配置，请使用 Config.Set, ParseMap 等方法修改
func (m *Manager[T]) Vars() T {
	m.confLock.RLock()
	defer m.confLock.RUnlock()
	return *m.userConf
}

// Raws 返回当前配置的读写副本（线程不安全）
// 返回值是指向配置结构体的指针，调用者不应修改返回的配置对象
// 如果需要更新配置，请使用 Config.Set, ParseMap 等方法修改
func (m *Manager[T]) Raws() *T {
	return m.userConf
}

type Config struct {
	key string
	val reflect.Value
}

func (m *Manager[T]) Conf(confKey string) (*Config, error) {
	if value, exist := m.valueMap[confKey]; exist {
		return value, nil
	}
	return nil, errors.Wrapf(ConfNotExistError, "Conf key[%s] is not exist", confKey)
}

func (c *Config) SetByString(value string) error {
	var convErr error
	switch c.val.Kind() {
	case reflect.String:
		c.val.SetString(value)
	case reflect.Bool:
		var boolVal bool
		boolVal, convErr = strconv.ParseBool(value)
		if convErr == nil {
			c.val.SetBool(boolVal)
		}
	case reflect.Int, reflect.Int32, reflect.Int64:
		var intVal int64
		intVal, convErr = strconv.ParseInt(value, 10, 64)
		if convErr == nil {
			c.val.SetInt(intVal)
		}
	case reflect.Float32, reflect.Float64:
		var floatVal float64
		floatVal, convErr = strconv.ParseFloat(value, 64)
		if convErr == nil {
			c.val.SetFloat(floatVal)
		}
	default:
		return errors.Wrapf(ConfValueSetError, "Not suppose set value for conf[%s](%s) ", c.key, c.val.Kind().String())
	}
	if convErr != nil {
		return errors.Wrapf(ConfValueSetError, "Parse value to set conf[%s](%s) failed. %s", c.key, c.val.Kind().String(), convErr.Error())
		//return errors.Newf("Parse value to set conf[%s](%s) failed. %s", c.key, c.val.Kind().String(), convErr.Error())
	}
	return nil
}

func (c *Config) Set(value any) {
	switch c.val.Kind() {
	case reflect.String:
		c.val.SetString(value.(string))
	case reflect.Bool:
		c.val.SetBool(value.(bool))
	case reflect.Int, reflect.Int32, reflect.Int64:
		c.val.SetInt(value.(int64))
	case reflect.Float32, reflect.Float64:
		c.val.SetFloat(value.(float64))
	default:
		log.Warnf("Not suppose set value for conf[%s](%s) ", c.key, c.val.Kind().String())
	}
}

func (c *Config) OnChange(callback func(any)) {

}
