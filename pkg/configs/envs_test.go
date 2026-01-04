package configs

import (
	"os"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestManager_parseEnvs(t *testing.T) {

	envKeys := map[string]string{
		"ENABLE_SERVICE":  "true",
		"SERVICE_NAME":    "test_service",
		"SERVICE_PORT":    "3306",
		"CLUSTER_NODES":   "12345678910",
		"ROLLOVER_FACTOR": "0.6",
	}

	mock := gomonkey.ApplyFunc(os.Getenv, func(key string) string {
		if val, exist := envKeys[key]; exist {
			return val
		}
		return ""
	})
	defer mock.Reset()

	manager := NewManager[testConf](nil)
	config := manager.Vars()
	assert.Equal(t, config.EnableService, true)
	assert.Equal(t, config.ServiceName, "test_service")
	assert.Equal(t, config.ServicePort, 3306)
	assert.Equal(t, config.ClusterNodes, int64(12345678910))
	assert.Equal(t, config.RolloverFactor, 0.6)
}
