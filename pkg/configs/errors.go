package configs

import "github.com/wolfbolin/bolbox/pkg/errors"

var (
	ErrConfNotExist = errors.New("Config key is not exist.")
	ErrConfValueSet = errors.New("Config value can not be set.")
	ErrParseFlags   = errors.New("Parse command flags error.")
	ErrPrintUsage   = errors.New("User request to print usage.")
)
