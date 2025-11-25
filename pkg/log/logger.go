package log

// Logger is an abstract definition of the functionality of a logging component
type Logger interface {
	Log(level Level, msg string, keyvalues ...interface{})
	Enabled(level Level) bool
}

// Level is the definition of the log level in the log component
type Level int

const (
	// DebugLevel means Debug and above log level
	DebugLevel Level = iota
	// InfoLevel means Info and above log level
	InfoLevel
	// WarnLevel means Warn and above log level
	WarnLevel
	// ErrorLevel means Error and above log level
	ErrorLevel
	// PanicLevel means Panic and above log level
	PanicLevel
	// FatalLevel means Fatal and above log level
	FatalLevel
)

// String provides a method for converting log level files to strings
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "DEBUG"
	case InfoLevel:
		return "INFO"
	case WarnLevel:
		return "WARN"
	case ErrorLevel:
		return "ERROR"
	case PanicLevel:
		return "PANIC"
	case FatalLevel:
		return "FATAL"
	default:
		return "DEBUG"
	}
}

// ParseLevel provides a method for creating log level objects through strings
func ParseLevel(l string) Level {
	switch l {
	case "DEBUG":
		return DebugLevel
	case "INFO":
		return InfoLevel
	case "WARN":
		return WarnLevel
	case "ERROR":
		return ErrorLevel
	case "PANIC":
		return PanicLevel
	case "FATAL":
		return FatalLevel
	default:
		return DebugLevel
	}
}
