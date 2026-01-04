package configs

import "github.com/wolfbolin/bolbox/pkg/errors"

var (
	ConfNotExistError = errors.New("Config key is not exist.")
	ConfValueSetError = errors.New("Config value can not be set.")
)
