package zap

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewZapLogger 根据给定的日志选项创建日志记录器
func NewZapLogger(option *Option) *Logger {
	var writerSyncer []zapcore.WriteSyncer
	if option.FilePath != "" {
		writerSyncer = append(writerSyncer, zapcore.AddSync(newLumberjackSyncer(option)))
	}

	if option.ConsoleLogger {
		writerSyncer = append(writerSyncer, zapcore.AddSync(os.Stdout))
	}

	core := zapcore.NewCore(
		newConsoleEncoder(option),
		zapcore.NewMultiWriteSyncer(writerSyncer...),
		option.LogLevel,
	)

	zapOptions := []zap.Option{zap.AddCaller(), zap.AddCallerSkip(option.CallSkip)}
	if option.Stack {
		zapOptions = append(zapOptions, zap.AddStacktrace(zap.ErrorLevel))
	}

	return NewLogger(zap.New(core, zapOptions...))
}

func newLumberjackSyncer(option *Option) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   option.FilePath,
		MaxSize:    option.MaxSize,
		MaxBackups: option.MaxBackups,
		MaxAge:     option.MaxAge,
		Compress:   true,
		LocalTime:  true,
	}
}

func newConsoleEncoder(option *Option) zapcore.Encoder {
	encoderCfg := zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		MessageKey:       "msg",
		CallerKey:        "line",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		EncodeDuration:   zapcore.SecondsDurationEncoder,
		EncodeName:       zapcore.FullNameEncoder,
		ConsoleSeparator: option.ConsoleSeparator,
	}

	if option.PathEncoderType == ShortPathEncoderType {
		encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	} else if option.PathEncoderType == FullPathEncoderType {
		encoderCfg.EncodeCaller = zapcore.FullCallerEncoder
	}

	return zapcore.NewConsoleEncoder(encoderCfg)
}
