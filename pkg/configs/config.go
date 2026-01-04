package configs

import (
	"sync"
)

// Manager 是一个泛型配置管理器，支持从环境变量、命令行参数和动态映射中加载配置
// 配置加载的优先级顺序为：默认值 < 环境变量 < 命令行参数 < 动态映射（ParseMap）
// 该管理器是线程安全的，可以在并发环境中安全使用
type Manager[T any] struct {
	userConf *T
	confLock sync.RWMutex
}

// NewManager 创建一个新的配置管理器实例
// 参数说明：
//   - def: 默认配置结构体指针，如果为 nil 则使用类型 T 的零值作为默认配置
//
// 创建管理器时会自动执行以下操作：
//   1. 使用默认配置初始化配置对象
//   2. 解析环境变量并更新配置（通过字段的 "env" 标签）
//   3. 解析命令行参数并更新配置（通过字段的 "flag" 标签）
//
// 使用示例：
//   type Config struct {
//       Port    int    `env:"PORT" flag:"port" desc:"服务端口"`
//       Debug   bool   `env:"DEBUG" flag:"debug" desc:"调试模式"`
//       LogFile string `env:"LOG_FILE" flag:"log-file" desc:"日志文件路径"`
//   }
//   mgr := NewManager(&Config{Port: 8080, Debug: false})
func NewManager[T any](def *T) *Manager[T] {
	var userConfig T
	if def != nil {
		userConfig = *def
	}
	mgr := &Manager[T]{
		userConf: &userConfig,
	}
	mgr.parseEnvs()
	mgr.parseFlags()
	return mgr
}

// Vars 返回当前配置的只读副本（线程安全）
// 返回值是指向配置结构体的指针，调用者不应修改返回的配置对象
// 如果需要更新配置，请使用 ParseMap 方法
func (m *Manager[T]) Vars() *T {
	m.confLock.RLock()
	defer m.confLock.RUnlock()
	return m.userConf
}
