package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultOptions(t *testing.T) {
	// 测试默认选项
	defaultOpts := DefaultOptions()
	assert.True(t, defaultOpts.ExitOnHelp)
	assert.Equal(t, []Flow{FlowEnv, FlowFlag}, defaultOpts.ParseFlows)
}

func TestParseFlowsOrder(t *testing.T) {
	// 测试不同的 ParseFlows 顺序
	testCases := []struct {
		name      string
		parseFlow []Flow
		envValue  string
		flagValue string
		expected  string
	}{
		{
			name:      "Env then Flag",
			parseFlow: []Flow{FlowEnv, FlowFlag},
			envValue:  "env-value",
			flagValue: "flag-value",
			expected:  "flag-value", // Flag 应该覆盖 Env
		},
		{
			name:      "Flag then Env",
			parseFlow: []Flow{FlowFlag, FlowEnv},
			envValue:  "env-value",
			flagValue: "flag-value",
			expected:  "env-value", // Env 应该覆盖 Flag
		},
		{
			name:      "Only Env",
			parseFlow: []Flow{FlowEnv},
			envValue:  "env-value",
			flagValue: "flag-value",
			expected:  "env-value", // 只使用 Env
		},
		{
			name:      "Only Flag",
			parseFlow: []Flow{FlowFlag},
			envValue:  "env-value",
			flagValue: "flag-value",
			expected:  "flag-value", // 只使用 Flag
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 设置环境变量
			os.Setenv("STRING_FIELD", tc.envValue)
			// 设置命令行参数
			os.Args = []string{"test", "--string-field", tc.flagValue}

			// 创建配置管理器
			mgr := NewManager(&flagTestConf{})
			// 设置自定义的 ParseFlows 顺序
			mgr.Options.ParseFlows = tc.parseFlow

			// 解析配置
			_, err := mgr.Parse()
			assert.NoError(t, err)

			// 获取配置值并验证
			config := mgr.Vars()
			assert.Equal(t, tc.expected, config.StringField)
		})
	}
}
