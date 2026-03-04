package configs

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestCheckAndShowHelp(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		shouldExit bool
	}{
		{
			name:       "WithHelpFlag",
			args:       []string{"program", "--help"},
			shouldExit: true,
		},
		{
			name:       "WithHFlag",
			args:       []string{"program", "-h"},
			shouldExit: true,
		},
		{
			name:       "WithHelpInMiddle",
			args:       []string{"program", "arg1", "--help", "arg2"},
			shouldExit: true,
		},
		{
			name:       "WithoutHelpFlag",
			args:       []string{"program", "arg1", "arg2"},
			shouldExit: false,
		},
		{
			name:       "WithOtherFlags",
			args:       []string{"program", "--other", "value"},
			shouldExit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置测试用的 os.Args
			os.Args = tt.args

			// 创建 flagSet
			flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
			usageCalled := false
			exitCalled := false
			exitCode := -1

			// Mock flagSet.Usage
			flagSet.Usage = func() {
				usageCalled = true
			}

			// Mock os.Exit，只设置标志变量，不真正退出也不 panic
			patches := gomonkey.ApplyFunc(os.Exit, func(code int) {
				exitCalled = true
				exitCode = code
			})
			defer patches.Reset()

			// 调用函数
			checkAndShowHelp(flagSet)

			if tt.shouldExit {
				// 验证 Usage 被调用
				assert.True(t, usageCalled, "Usage should be called when help flag is present")
				// 验证 Exit 被调用
				assert.True(t, exitCalled, "Exit should be called when help flag is present")
				// 验证退出码为 0
				assert.Equal(t, 0, exitCode, "Exit code should be 0")
			} else {
				// 验证 Usage 没有被调用
				assert.False(t, usageCalled, "Usage should not be called when help flag is absent")
				// 验证 Exit 没有被调用
				assert.False(t, exitCalled, "Exit should not be called when help flag is absent")
			}
		})
	}
}

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
			manager, err := NewManager[flagTestConf](&tt.defaultConfig)
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

func TestFilterKnownArgs(t *testing.T) {
	// 构造一个注册了 known-flag 和 bool-flag 的 flagSet
	newFlagSet := func() *flag.FlagSet {
		fs := flag.NewFlagSet("test", flag.ContinueOnError)
		fs.String("known-flag", "", "")
		fs.Bool("bool-flag", false, "")
		return fs
	}

	tests := []struct {
		name     string
		args     []string
		expected []string
	}{
		{
			name:     "全部是已知flag",
			args:     []string{"--known-flag", "value", "--bool-flag"},
			expected: []string{"--known-flag", "value", "--bool-flag"},
		},
		{
			name:     "全部是未知flag，应全部过滤",
			args:     []string{"--unknown-flag", "value", "--another-unknown"},
			expected: []string{},
		},
		{
			name:     "已知与未知混合，保留已知",
			args:     []string{"--unknown-flag", "value", "--known-flag", "hello"},
			expected: []string{"--known-flag", "hello"},
		},
		{
			name:     "使用等号形式的已知flag",
			args:     []string{"--known-flag=hello"},
			expected: []string{"--known-flag=hello"},
		},
		{
			name:     "使用等号形式的未知flag，应过滤",
			args:     []string{"--unknown-flag=world"},
			expected: []string{},
		},
		{
			name:     "单横线前缀的已知flag",
			args:     []string{"-known-flag", "value"},
			expected: []string{"-known-flag", "value"},
		},
		{
			name:     "单横线前缀的未知flag，应过滤",
			args:     []string{"-unknown", "value"},
			expected: []string{},
		},
		{
			name:     "未知flag后跟另一个flag（非位置参数），只过滤未知flag本身",
			args:     []string{"--unknown-flag", "--known-flag", "hello"},
			expected: []string{"--known-flag", "hello"},
		},
		{
			name:     "位置参数（非flag）直接保留",
			args:     []string{"positional", "--known-flag", "value"},
			expected: []string{"positional", "--known-flag", "value"},
		},
		{
			name:     "空参数列表",
			args:     []string{},
			expected: []string{},
		},
		{
			name:     "已知bool flag无值形式",
			args:     []string{"--bool-flag"},
			expected: []string{"--bool-flag"},
		},
		{
			name:     "未知bool flag无值，紧跟已知flag，不误吞已知flag的value",
			args:     []string{"--unknown-bool", "--known-flag", "myval"},
			expected: []string{"--known-flag", "myval"},
		},
		{
			name:     "未知flag夹在两个已知flag之间",
			args:     []string{"--known-flag", "v1", "--unknown", "x", "--bool-flag"},
			expected: []string{"--known-flag", "v1", "--bool-flag"},
		},
		{
			name:     "未知flag后跟负数值，应将负数值一并跳过",
			args:     []string{"--unknown", "-42", "--known-flag", "kept"},
			expected: []string{"--known-flag", "kept"},
		},
		{
			name:     "已知flag的value为负数，负数不被looksLikeFlag识别为flag，完整保留",
			args:     []string{"--known-flag", "-42"},
			expected: []string{"--known-flag", "-42"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := newFlagSet()
			got := filterKnownArgs(fs, tt.args)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestParseFlags_WithUnknownFlags(t *testing.T) {
	// 验证 parseFlags 在遇到未知 flag 时不会报错
	mockEnv := func(key string) string { return "" }
	_ = mockEnv

	tests := []struct {
		name           string
		args           []string
		defaultConfig  flagTestConf
		expectedConfig flagTestConf
	}{
		{
			name: "夹杂未知flag，已知flag仍正常解析",
			args: []string{
				"program",
				"--unknown-flag", "ignored",
				"--string-field", "parsed",
				"--another-unknown=also-ignored",
			},
			defaultConfig:  flagTestConf{StringField: "default"},
			expectedConfig: flagTestConf{StringField: "parsed"},
		},
		{
			name: "全部为未知flag，保持默认值",
			args: []string{
				"program",
				"--foo", "bar",
				"--baz=qux",
			},
			defaultConfig:  flagTestConf{StringField: "default", IntField: 99},
			expectedConfig: flagTestConf{StringField: "default", IntField: 99},
		},
		{
			name: "未知flag在已知flag之前",
			args: []string{
				"program",
				"--unknown", "value",
				"--int-field", "42",
			},
			defaultConfig:  flagTestConf{IntField: 0},
			expectedConfig: flagTestConf{IntField: 42},
		},
		{
			name: "未知flag在已知flag之后",
			args: []string{
				"program",
				"--string-field", "hello",
				"--unknown", "world",
			},
			defaultConfig:  flagTestConf{StringField: "default"},
			expectedConfig: flagTestConf{StringField: "hello"},
		},
		{
			name: "等号形式未知flag不影响已知flag解析",
			args: []string{
				"program",
				"--unknown=ignored",
				"--bool-field",
			},
			defaultConfig:  flagTestConf{BoolField: false},
			expectedConfig: flagTestConf{BoolField: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = tt.args
			manager, err := NewManager[flagTestConf](&tt.defaultConfig)
			assert.Nil(t, err, "parseFlags should not error on unknown flags")

			config := manager.Vars()
			assert.Equal(t, tt.expectedConfig.StringField, config.StringField, "StringField mismatch")
			assert.Equal(t, tt.expectedConfig.BoolField, config.BoolField, "BoolField mismatch")
			assert.Equal(t, tt.expectedConfig.IntField, config.IntField, "IntField mismatch")
		})
	}
}

func TestParseFlagName(t *testing.T) {
	tests := []struct {
		arg              string
		expectedName     string
		expectedHasValue bool
	}{
		// 双横线形式
		{"--foo", "foo", false},
		{"--foo=bar", "foo", true},
		{"--foo=", "foo", true},
		// 单横线形式
		{"-foo", "foo", false},
		{"-foo=bar", "foo", true},
		// 非 flag（位置参数）
		{"foo", "", false},
		{"", "", false},
		// 仅横线
		{"-", "", false}, // "-" 去掉前缀后为空串，视为非 flag
		// 等号出现在 value 里时只取第一个等号前的内容
		{"--key=val=extra", "key", true},
	}

	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			name, hasValue := parseFlagName(tt.arg)
			assert.Equal(t, tt.expectedName, name, "name mismatch for arg %q", tt.arg)
			assert.Equal(t, tt.expectedHasValue, hasValue, "hasValue mismatch for arg %q", tt.arg)
		})
	}
}

func TestLooksLikeFlag(t *testing.T) {
	tests := []struct {
		arg      string
		expected bool
	}{
		// 明确是 flag
		{"--foo", true},
		{"-foo", true},
		{"-f", true},
		// 负数，不应被视为 flag
		{"-1", false},
		{"-42", false},
		{"-0", false},
		// 边界：空串或单横线
		{"", false},
		{"-", false},
		// 无横线前缀
		{"foo", false},
		{"123", false},
	}

	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			got := looksLikeFlag(tt.arg)
			assert.Equal(t, tt.expected, got, "looksLikeFlag(%q)", tt.arg)
		})
	}
}
