package configs

import (
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestManager_parseEnvs(t *testing.T) {

	envKeys := map[string]string{
		"STRING_FIELD": "true",
	}

	mock := gomonkey.ApplyFunc(os.Getenv, func(key string) string {
		if val, exist := envKeys[key]; exist {
			return val
		}
		return ""
	})
	defer mock.Reset()

	manager, err := NewManager[flagTestConf](nil)
	assert.Nil(t, err)
	config := manager.Vars()
	assert.Equal(t, config.StringField, "true")
}
