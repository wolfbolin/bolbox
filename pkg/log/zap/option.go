package zap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/wolfbolin/bolbox/pkg/log"
)

// PathEncoderType 表示路径编码器类型
type PathEncoderType int

const (
	// NoPathEncoderType 不打印文件路径
	NoPathEncoderType = iota
	// ShortPathEncoderType 打印段文件路径
	ShortPathEncoderType
	// FullPathEncoderType 打印完整文件路径
	FullPathEncoderType
)

// Option 包含 zap 日志库的可定制选项
type Option struct {
	FilePath      string // 如果为空则表示不向文件打印日志
	LogLevel      zap.AtomicLevel
	ConsoleLogger bool // 是否打印控制台日志
	CallSkip      int  // 跳过调用函数的数量
	Stack         bool // 是否打印堆栈

	PathEncoderType  PathEncoderType // 文件路径编码类型：不打印、段路径编码、全路径编码
	ConsoleSeparator string          // 日志字段的分隔符

	MaxSize    int // 每个日志文件保存的最大尺寸 单位：M
	MaxBackups int // 日志文件最多保存多少个备份
	MaxAge     int // 文件最多保存多少天
}

// NewDefaultOption 返回默认的 zap 日志选项
func NewDefaultOption(path string, level log.Level) *Option {
	return &Option{
		FilePath:         path,
		LogLevel:         zap.NewAtomicLevelAt(zapcore.Level(level - 1)),
		ConsoleLogger:    true,
		CallSkip:         2,
		Stack:            false,
		PathEncoderType:  ShortPathEncoderType,
		ConsoleSeparator: " | ",
		MaxSize:          20,
		MaxBackups:       50,
		MaxAge:           30,
	}
}

// SetLogLevel 动态设置 zap 日志的级别
func (o *Option) SetLogLevel(level log.Level) {
	o.LogLevel.SetLevel(zapcore.Level(level - 1))
}
