package configs

import (
	"reflect"
	"sync"
)

// Manager 配置管理器，支持从环境变量、命令行参数和动态映射中加载配置
// 优先级：默认值 < 环境变量 < 命令行参数 < 动态映射
type Manager[T any] struct {
	userConf *T
	valueMap map[string]*Config
	confElem reflect.Value
	confLock sync.RWMutex
}

// NewManager 创建配置管理器，自动解析环境变量和命令行参数
func NewManager[T any](def *T) (*Manager[T], error) {
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
			cbs: make([]func(any), 0),
		}
	}

	err := mgr.parseEnvs()
	if err != nil {
		return nil, err
	}
	err = mgr.parseFlags()
	if err != nil {
		return nil, err
	}
	return mgr, nil
}

// Vars 返回配置的只读副本
func (m *Manager[T]) Vars() T {
	m.confLock.RLock()
	defer m.confLock.RUnlock()
	return *m.userConf
}

// Raws 返回配置的原始指针（线程不安全）
func (m *Manager[T]) Raws() *T {
	return m.userConf
}
