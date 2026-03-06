package configs

import (
	"fmt"
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/wolfbolin/bolbox/pkg/errors"
)

func TestManager_parseFlags(t *testing.T) {
	// Mock 环境变量，确保环境变量不会影响测试
	mockEnv := gomonkey.ApplyFunc(os.Getenv, func(key string) string {
		return ""
	})
	defer mockEnv.Reset()

	tests := []struct {
		name           string
		args           []string
		defaultConfig  flagTestConf
		expectedConfig flagTestConf
		description    string
	}{
		{
			name: "解析所有数据类型",
			args: []string{
				"program",
				"--string-field", "test-string",
				"--bool-field",
				"--int-field", "42",
				"--int32-field", "32",
				"--int64-field", "64",
				"--float32-field", "3.14",
				"--float64-field", "2.718",
				"--unknown-field", "unknown",
			},
			defaultConfig: flagTestConf{
				StringField:  "default",
				BoolField:    false,
				IntField:     0,
				Int32Field:   0,
				Int64Field:   0,
				Float32Field: 0.0,
				Float64Field: 0.0,
				NoFlagField:  "no-flag",
			},
			expectedConfig: flagTestConf{
				StringField:  "test-string",
				BoolField:    true,
				IntField:     42,
				Int32Field:   32,
				Int64Field:   64,
				Float32Field: 3.14,
				Float64Field: 2.718,
				NoFlagField:  "no-flag", // 没有 flag 标签，应该保持默认值
			},
			description: "验证所有支持的数据类型都能正确解析",
		},
		{
			name: "部分字段解析",
			args: []string{
				"program",
				"--string-field", "partial",
				"--int-field", "100",
			},
			defaultConfig: flagTestConf{
				StringField:  "default",
				BoolField:    false,
				IntField:     0,
				Int32Field:   0,
				Int64Field:   0,
				Float32Field: 0.0,
				Float64Field: 0.0,
			},
			expectedConfig: flagTestConf{
				StringField:  "partial",
				BoolField:    false, // 未提供，保持默认值
				IntField:     100,
				Int32Field:   0, // 未提供，保持默认值
				Int64Field:   0, // 未提供，保持默认值
				Float32Field: 0.0,
				Float64Field: 0.0,
			},
			description: "验证只提供部分字段时，其他字段保持默认值",
		},
		{
			name: "布尔类型false值",
			args: []string{
				"program",
				"--bool-field=false",
			},
			defaultConfig: flagTestConf{
				BoolField: true,
			},
			expectedConfig: flagTestConf{
				BoolField: false,
			},
			description: "验证布尔类型可以设置为false",
		},
		{
			name: "整数类型边界值",
			args: []string{
				"program",
				"--int-field", "-2147483648",
				"--int32-field", "2147483647",
				"--int64-field", "9223372036854775807",
			},
			defaultConfig: flagTestConf{},
			expectedConfig: flagTestConf{
				IntField:   -2147483648,
				Int32Field: 2147483647,
				Int64Field: 9223372036854775807,
			},
			description: "验证整数类型的边界值",
		},
		{
			name: "浮点数类型",
			args: []string{
				"program",
				"--float32-field", "123.456",
				"--float64-field", "789.0123456789",
			},
			defaultConfig: flagTestConf{},
			expectedConfig: flagTestConf{
				Float32Field: 123.456,
				Float64Field: 789.0123456789,
			},
			description: "验证浮点数类型的解析",
		},
		{
			name: "无命令行参数时使用默认值",
			args: []string{
				"program",
			},
			defaultConfig: flagTestConf{
				StringField:  "default-string",
				BoolField:    true,
				IntField:     100,
				Int32Field:   200,
				Int64Field:   300,
				Float32Field: 1.5,
				Float64Field: 2.5,
			},
			expectedConfig: flagTestConf{
				StringField:  "default-string",
				BoolField:    true,
				IntField:     100,
				Int32Field:   200,
				Int64Field:   300,
				Float32Field: 1.5,
				Float64Field: 2.5,
			},
			description: "验证没有命令行参数时，使用默认值",
		},
		{
			name: "空字符串值",
			args: []string{
				"program",
				"--string-field", "",
			},
			defaultConfig: flagTestConf{
				StringField: "non-empty",
			},
			expectedConfig: flagTestConf{
				StringField: "",
			},
			description: "验证字符串类型可以设置为空字符串",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置测试用的 os.Args
			os.Args = tt.args

			// 创建管理器，这会自动调用 parseFlags
			manager, err := NewManager[flagTestConf](&tt.defaultConfig).Parse()
			if err != nil {
				fmt.Printf("Test failed stack: %+v\n", err)
			}
			assert.Nil(t, err)

			// 获取配置并验证
			config := manager.Vars()

			// 验证所有字段
			assert.Equal(t, tt.expectedConfig.StringField, config.StringField, "StringField should match")
			assert.Equal(t, tt.expectedConfig.BoolField, config.BoolField, "BoolField should match")
			assert.Equal(t, tt.expectedConfig.IntField, config.IntField, "IntField should match")
			assert.Equal(t, tt.expectedConfig.Int32Field, config.Int32Field, "Int32Field should match")
			assert.Equal(t, tt.expectedConfig.Int64Field, config.Int64Field, "Int64Field should match")
			assert.Equal(t, tt.expectedConfig.Float32Field, config.Float32Field, "Float32Field should match")
			assert.Equal(t, tt.expectedConfig.Float64Field, config.Float64Field, "Float64Field should match")
			assert.Equal(t, tt.expectedConfig.NoFlagField, config.NoFlagField, "NoFlagField should match (no flag tag)")
		})
	}
}

func TestManager_showHelp(t *testing.T) {
	os.Args = []string{"program", "--help"}

	defaultConfig := flagTestConf{
		StringField: "non-empty",
	}

	manager := NewManager[flagTestConf](&defaultConfig)
	manager.Options.ExitOnHelp = false
	_, err := manager.Parse()
	assert.True(t, errors.Is(err, ErrPrintUsage))
}
