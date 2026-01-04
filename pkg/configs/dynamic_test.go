package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager_ParseMap(t *testing.T) {
	// 创建测试配置管理器
	manager := NewManager[testConf](&testConf{
		EnableService: false,
		ServiceName:   "default-service",
		ServicePort:   8080,
		ClusterNodes:  5,
		RolloverFactor: 0.5,
	})

	// 测试用例1: 正常情况 - 更新所有字段
	t.Run("UpdateAllFields", func(t *testing.T) {
		manager.ParseMap(map[string]string{
			"EnableService":  "true",
			"ServiceName":    "test-service",
			"ServicePort":    "9090",
			"ClusterNodes":   "10",
			"RolloverFactor": "0.75",
		})

		config := manager.Vars()
		assert.Equal(t, true, config.EnableService)
		assert.Equal(t, "test-service", config.ServiceName)
		assert.Equal(t, 9090, config.ServicePort)
		assert.Equal(t, int64(10), config.ClusterNodes)
		assert.Equal(t, 0.75, config.RolloverFactor)
	})

	// 测试用例2: nil map
	t.Run("NilMapShouldNotChangeAnything", func(t *testing.T) {
		// 先设置一些值
		manager.ParseMap(map[string]string{
			"ServiceName": "before-nil",
		})
		config1 := manager.Vars()
		assert.Equal(t, "before-nil", config1.ServiceName)

		// 传入 nil map
		manager.ParseMap(nil)

		// 值应该保持不变
		config2 := manager.Vars()
		assert.Equal(t, "before-nil", config2.ServiceName)
	})

	// 测试用例3: 部分字段更新
	t.Run("UpdatePartialFields", func(t *testing.T) {
		// 重置为初始值
		manager.ParseMap(map[string]string{
			"EnableService": "false",
			"ServiceName":   "initial",
			"ServicePort":   "8080",
		})

		// 只更新一个字段
		manager.ParseMap(map[string]string{
			"ServicePort": "3000",
		})

		config := manager.Vars()
		assert.Equal(t, false, config.EnableService)
		assert.Equal(t, "initial", config.ServiceName)
		assert.Equal(t, 3000, config.ServicePort)
	})

	// 测试用例4: 不存在的字段名（应该被忽略）
	t.Run("NonExistentFieldShouldBeIgnored", func(t *testing.T) {
		originalConfig := manager.Vars()

		// 尝试更新不存在的字段
		manager.ParseMap(map[string]string{
			"NonExistentField": "some-value",
			"ServicePort":      "7777",
		})

		config := manager.Vars()
		// ServicePort 应该被更新
		assert.Equal(t, 7777, config.ServicePort)
		// 其他字段应该保持不变
		assert.Equal(t, originalConfig.EnableService, config.EnableService)
		assert.Equal(t, originalConfig.ServiceName, config.ServiceName)
	})

	// 测试用例5: 无效的值（类型转换失败，应该被忽略）
	t.Run("InvalidValueShouldBeIgnored", func(t *testing.T) {
		originalConfig := manager.Vars()

		// 尝试设置无效的值
		manager.ParseMap(map[string]string{
			"ServicePort":   "not-a-number",
			"EnableService": "not-a-bool",
			"ServiceName":   "valid-string",
		})

		config := manager.Vars()
		// ServiceName 应该被更新（字符串类型）
		assert.Equal(t, "valid-string", config.ServiceName)
		// ServicePort 和 EnableService 应该保持不变（转换失败）
		assert.Equal(t, originalConfig.ServicePort, config.ServicePort)
		assert.Equal(t, originalConfig.EnableService, config.EnableService)
	})

	// 测试用例6: 空 map
	t.Run("EmptyMapShouldNotChangeAnything", func(t *testing.T) {
		originalConfig := manager.Vars()

		manager.ParseMap(map[string]string{})

		config := manager.Vars()
		assert.Equal(t, originalConfig, config)
	})
}
