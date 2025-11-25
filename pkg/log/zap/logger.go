package zap

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/wolfbolin/bolbox/pkg/log"
)

var _ log.Logger = (*Logger)(nil)

// Logger 定义了一个 zap 日志记录器
type Logger struct {
	logger *zap.Logger
}

// NewLogger 创建并返回一个 Logger 实例
func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

// Enabled 检查给定的日志级别是否启用
func (l *Logger) Enabled(level log.Level) bool {
	return l.logger.Core().Enabled(zapcore.Level(level - 1))
}

// Log 用于记录用户日志
func (l *Logger) Log(level log.Level, msg string, keyvals ...interface{}) {
	if len(keyvals)%2 != 0 {
		l.logger.Warn("keyvals must appear int pairs, len: %d", zap.Int("len", len(keyvals)))
		return
	}

	logData := make([]zap.Field, 0, len(keyvals)/2)
	for i := 0; i < len(keyvals); i += 2 {
		logData = append(logData, zap.Any(fmt.Sprint(keyvals[i]), keyvals[i+1]))
	}

	switch level {
	case log.DebugLevel:
		l.logger.Debug(msg, logData...)
	case log.InfoLevel:
		l.logger.Info(msg, logData...)
	case log.WarnLevel:
		l.logger.Warn(msg, logData...)
	case log.ErrorLevel:
		l.logger.Error(msg, logData...)
	case log.PanicLevel:
		l.logger.Panic(msg, logData...)
	case log.FatalLevel:
		l.logger.Fatal(msg, logData...)
	default:
		l.logger.Debug(msg, logData...)
	}
}

// Sync 用于确保日志被写入
func (l *Logger) Sync() error {
	return l.logger.Sync()
}
