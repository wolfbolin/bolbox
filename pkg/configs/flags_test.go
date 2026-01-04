package configs

import (
	"flag"
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
			// 保存原始 os.Args
			originalArgs := os.Args
			defer func() {
				os.Args = originalArgs
			}()

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
