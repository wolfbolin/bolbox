package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestManager_parseMap(t *testing.T) {
	// 创建测试配置管理器
	manager, err := NewManager[flagTestConf](&flagTestConf{
		BoolField:    false,
		StringField:  "default-service",
		IntField:     43,
		Int64Field:   5,
		Float64Field: 0.5,
	})
	assert.Nil(t, err)

	// 测试用例1: 正常情况 - 更新所有字段
	t.Run("UpdateAllFields", func(t *testing.T) {
		manager.ParseMap(map[string]string{
			"BoolField":    "true",
			"StringField":  "test-string",
			"IntField":     "72",
			"Int64Field":   "12345678910",
			"Float64Field": "0.75",
		})

		config := manager.Vars()
		assert.Equal(t, true, config.BoolField)
		assert.Equal(t, "test-string", config.StringField)
		assert.Equal(t, 72, config.IntField)
		assert.Equal(t, int64(12345678910), config.Int64Field)
		assert.Equal(t, 0.75, config.Float64Field)
	})

	// 测试用例2: nil map
	t.Run("NilMapShouldNotChangeAnything", func(t *testing.T) {
		// 先设置一些值
		manager.ParseMap(map[string]string{
			"StringField": "before-nil",
		})
		config1 := manager.Vars()
		assert.Equal(t, "before-nil", config1.StringField)

		// 传入 nil map
		manager.ParseMap(nil)

		// 值应该保持不变
		config2 := manager.Vars()
		assert.Equal(t, "before-nil", config2.StringField)
	})

	// 测试用例3: 部分字段更新
	t.Run("UpdatePartialFields", func(t *testing.T) {
		// 重置为初始值
		manager.ParseMap(map[string]string{
			"BoolField":   "false",
			"StringField": "initial",
			"IntField":    "8080",
		})

		// 只更新一个字段
		manager.ParseMap(map[string]string{
			"IntField": "3000",
		})

		config := manager.Vars()
		assert.Equal(t, false, config.BoolField)
		assert.Equal(t, "initial", config.StringField)
		assert.Equal(t, 3000, config.IntField)
	})

	// 测试用例4: 不存在的字段名（应该被忽略）
	t.Run("NonExistentFieldShouldBeIgnored", func(t *testing.T) {
		originalConfig := manager.Vars()

		// 尝试更新不存在的字段
		manager.ParseMap(map[string]string{
			"NonExistentField": "some-value",
			"IntField":         "7777",
		})

		config := manager.Vars()
		// IntField 应该被更新
		assert.Equal(t, 7777, config.IntField)
		// 其他字段应该保持不变
		assert.Equal(t, originalConfig.BoolField, config.BoolField)
		assert.Equal(t, originalConfig.StringField, config.StringField)
	})

	// 测试用例5: 无效的值（类型转换失败，应该被忽略）
	t.Run("InvalidValueShouldBeIgnored", func(t *testing.T) {
		originalConfig := manager.Vars()

		// 尝试设置无效的值
		manager.ParseMap(map[string]string{
			"IntField":    "not-a-number",
			"BoolField":   "not-a-bool",
			"StringField": "valid-string",
		})

		config := manager.Vars()
		// StringField 应该被更新（字符串类型）
		assert.Equal(t, "valid-string", config.StringField)
		// IntField 和 BoolField 应该保持不变（转换失败）
		assert.Equal(t, originalConfig.IntField, config.IntField)
		assert.Equal(t, originalConfig.BoolField, config.BoolField)
	})

	// 测试用例6: 空 map
	t.Run("EmptyMapShouldNotChangeAnything", func(t *testing.T) {
		originalConfig := manager.Vars()

		manager.ParseMap(map[string]string{})

		config := manager.Vars()
		assert.Equal(t, originalConfig, config)
	})
}
