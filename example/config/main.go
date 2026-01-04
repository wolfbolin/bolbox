package main

import (
	"fmt"

	"github.com/wolfbolin/bolbox/example/config/cfg"
)

func main() {
	fmt.Println("=== Bolbox Config 示例 ===")
	fmt.Println()

	// 1. 读取配置（配置已自动从环境变量和命令行参数加载）
	printConfig("初始配置（默认值 + 环境变量 + 命令行参数）", cfg.Vars())
	fmt.Println()

	// 2. 演示动态配置更新（优先级最高）
	fmt.Println("正在通过 ParseMap 动态更新配置...")
	cfg.ParseMap(map[string]string{
		"ServiceName": "updated-service",
		"ServicePort": "9090",
		"Debug":       "true",
		"Timeout":     "60.5",
	})
	printConfig("动态更新后的配置", cfg.Vars())
	fmt.Println()
}

// printConfig 打印配置信息
func printConfig(title string, config *cfg.Config) {
	fmt.Printf("【%s】\n", title)
	fmt.Printf("  服务名称:     %s\n", config.ServiceName)
	fmt.Printf("  服务端口:     %d\n", config.ServicePort)
	fmt.Printf("  调试模式:     %v\n", config.Debug)
	fmt.Printf("  日志文件:     %s\n", config.LogFile)
	fmt.Printf("  超时时间:     %.2f 秒\n", config.Timeout)
	fmt.Printf("  数据库主机:   %s\n", config.DBHost)
	fmt.Printf("  数据库端口:   %d\n", config.DBPort)
	fmt.Printf("  版本号:       %s\n", config.Version)
}
