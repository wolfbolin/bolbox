package log

import (
	"log"
)

// DefaultLogger provides the default implementation of the logging component
type DefaultLogger struct {
}

// Log implements the log output of the default log component
func (d *DefaultLogger) Log(level Level, msg string, keyvals ...interface{}) {
	if len(keyvals) == 0 {
		log.Println(level.String(), "|", msg, " ")
	} else {
		log.Println(level.String(), "|", msg, " ", keyvals)
	}
}

// Enabled implements log level queries for default log components
func (d *DefaultLogger) Enabled(level Level) bool {
	return true
}
