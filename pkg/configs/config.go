package configs

import (
	"reflect"
	"strconv"
	"sync"

	"github.com/wolfbolin/bolbox/pkg/errors"
)

type Config struct {
	key  string
	val  reflect.Value
	cbs  []func(any)
	lock sync.RWMutex
}

// Conf 根据配置键获取配置对象
func (m *Manager[T]) Conf(confKey string) (*Config, error) {
	if value, exist := m.valueMap[confKey]; exist {
		return value, nil
	}
	return nil, errors.Wrapf(ConfNotExistError, "Conf key[%s] is not exist", confKey)
}

// SetByValue 直接设置配置值
func (c *Config) SetByValue(value any) error {
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
		return errors.Wrapf(ConfValueSetError, "Not suppose set value for conf[%s](%s) ", c.key, c.val.Kind().String())
	}
	c.notify(value)
	return nil
}

// SetByString 从字符串解析并设置配置值
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
	}
	c.notify(c.val.Interface())
	return nil
}

// OnChange 注册配置变更回调函数
func (c *Config) OnChange(callback func(any)) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cbs = append(c.cbs, callback)
}

func (c *Config) notify(val any) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	for _, f := range c.cbs {
		go f(val)
	}
}
