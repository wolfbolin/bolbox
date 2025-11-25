package log

import (
	"fmt"
	"sync"
)

var (
	globalLogger Logger
	mutex        sync.Mutex
)

func init() {
	SetLogger(&DefaultLogger{})
}

// SetLogger is configured here to use the global logging component
func SetLogger(logger Logger) {
	mutex.Lock()
	defer mutex.Unlock()
	globalLogger = logger
}

// Enabled is used to confirm the log level of the current configuration
func Enabled(level Level) bool {
	return globalLogger.Enabled(level)
}

// Debug is used to print debug logs
func Debug(msg string) {
	globalLogger.Log(DebugLevel, msg)
}

// Debugf is used to print formatted debug level logs
func Debugf(format string, args ...interface{}) {
	globalLogger.Log(DebugLevel, fmt.Sprintf(format, args...))
}

// Debugw is used to print debug level logs containing additional kv information
func Debugw(msg string, keyvals ...interface{}) {
	globalLogger.Log(DebugLevel, msg, keyvals...)
}

// Info is used to print info level logs
func Info(msg string) {
	globalLogger.Log(InfoLevel, msg)
}

// Infof is used to print formatted info level logs
func Infof(format string, args ...interface{}) {
	globalLogger.Log(InfoLevel, fmt.Sprintf(format, args...))
}

// Infow is used to print info level logs containing additional kv information
func Infow(msg string, keyvals ...interface{}) {
	globalLogger.Log(InfoLevel, msg, keyvals...)
}

// Warn is used to print warning level logs
func Warn(msg string) {
	globalLogger.Log(WarnLevel, msg)
}

// Warnf is used to print formatted warning level logs
func Warnf(format string, args ...interface{}) {
	globalLogger.Log(WarnLevel, fmt.Sprintf(format, args...))
}

// Warnw is used to print warning level logs containing additional kv information
func Warnw(msg string, keyvals ...interface{}) {
	globalLogger.Log(WarnLevel, msg, keyvals...)
}

// Error is used to print error level logs
func Error(msg string) {
	globalLogger.Log(ErrorLevel, msg)
}

// Errorf is used to print formatted error level logs
func Errorf(format string, args ...interface{}) {
	globalLogger.Log(ErrorLevel, fmt.Sprintf(format, args...))
}

// Errorw is used to print error level logs containing additional kv information
func Errorw(msg string, keyvals ...interface{}) {
	globalLogger.Log(ErrorLevel, msg, keyvals...)
}

// Panic is used to print panic level logs
// This function will call the panic(err) after the log is printed
func Panic(msg string) {
	globalLogger.Log(PanicLevel, msg)
}

// Panicf is used to print formatted panic level logs
// This function will call the panic(err) after the log is printed
func Panicf(format string, args ...interface{}) {
	globalLogger.Log(PanicLevel, fmt.Sprintf(format, args...))
}

// Panicw is used to print painc level logs containing additional kv information
// This function will call the panic(err) after the log is printed
func Panicw(msg string, keyvals ...interface{}) {
	globalLogger.Log(PanicLevel, msg, keyvals...)
}

// Fatal is used to print fatal level logs
// This function will call the os.Exit(1) to exit the process after the log is printed
func Fatal(msg string) {
	globalLogger.Log(FatalLevel, msg)
}

// Fatalf is used to print formatted fatal level logs
// This function will call the os.Exit(1) to exit the process after the log is printed
func Fatalf(format string, args ...interface{}) {
	globalLogger.Log(FatalLevel, fmt.Sprintf(format, args...))
}

// Fatalw is used to print fatal level logs containing additional kv information
// This function will call the os.Exit(1) to exit the process after the log is printed
func Fatalw(msg string, keyvals ...interface{}) {
	globalLogger.Log(FatalLevel, msg, keyvals...)
}
