package cfg

import (
	"github.com/wolfbolin/bolbox/pkg/configs"
	"github.com/wolfbolin/bolbox/pkg/log"
)

// Config 应用配置结构体
// 支持的字段类型：string、bool、int、int32、int64、float32、float64
// 标签说明：
//   - env: 环境变量名称（可选）
//   - flag: 命令行参数名称（可选）
//   - desc: 命令行参数的描述信息（可选，仅在使用 flag 标签时有效）
type Config struct {
	// 服务配置
	ServiceName string `flag:"service-name" env:"SERVICE_NAME" desc:"服务名称"`
	ServicePort int32  `flag:"service-port" env:"SERVICE_PORT" desc:"服务端口号"`

	// 应用配置
	Debug   bool    `flag:"debug" env:"DEBUG" desc:"是否启用调试模式"`
	LogFile string  `flag:"log-file" env:"LOG_FILE" desc:"日志文件路径"`
	Timeout float64 `flag:"timeout" env:"TIMEOUT" desc:"请求超时时间（秒）"`

	// 数据库配置（仅使用环境变量，不使用命令行参数）
	DBHost string `env:"DB_HOST" desc:"数据库主机地址"`
	DBPort int    `env:"DB_PORT" desc:"数据库端口号"`

	// 仅使用默认值，不使用环境变量和命令行参数
	Version string
}

var cm *configs.Manager[Config]

func init() {
	// 使用默认配置初始化管理器
	// 配置加载优先级：默认值 < 环境变量 < 命令行参数 < 动态映射（ParseMap）
	mgr, err := configs.NewManager[Config](&Config{
		ServiceName: "my-service",
		ServicePort: 8080,
		Debug:       false,
		LogFile:     "./app.log",
		Timeout:     30.0,
		DBHost:      "localhost",
		DBPort:      3306,
		Version:     "1.0.0",
	})
	if err != nil {
		log.Fatalf("Init config manager failed. %v", err)
	}
	cm = mgr

	// 其他初始化方式示例：
	// 1. 不使用默认值（所有字段使用零值）
	// cm = configs.NewManager[Config](nil)
	//
	// 2. 使用空结构体作为默认值
	// cm = configs.NewManager[Config](&Config{})
}

// Vars 返回当前配置的只读副本（线程安全）
func Vars() *Config {
	data := cm.Vars()
	return &data
}

// ParseMap 动态更新配置（线程安全）
// 优先级最高，会覆盖环境变量和命令行参数设置的值
func ParseMap(data map[string]string) {
	cm.ParseMap(data)
}
