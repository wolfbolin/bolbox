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
		"MAP_FIELD":    `{"key1":"value1","key2":"value2"}`,
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
	assert.Equal(t, len(config.MapField), 2)
}
