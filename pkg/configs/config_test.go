package configs

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wolfbolin/bolbox/pkg/errors"
)

type testConf struct {
	// 服务配置
	EnableService bool   `flag:"enable-service" env:"ENABLE_SERVICE" desc:"启用服务"`
	ServiceName   string `flag:"service-name" env:"SERVICE_NAME" desc:"服务名称"`
	ServicePort   int    `flag:"service-port" env:"SERVICE_PORT" desc:"服务端口号"`

	ClusterNodes   int64   `flag:"cluster-nodes" env:"CLUSTER_NODES" desc:"集群节点数"`
	RolloverFactor float64 `flag:"rollover-factor" env:"ROLLOVER_FACTOR" desc:"滚动系数"`
}

func TestNewManager(t *testing.T) {
	manager := NewManager[testConf](nil)
	assert.NotNil(t, manager, manager.userConf, manager.valueMap, manager.confElem)

	confTable := testConf{}
	manager = NewManager[testConf](&confTable)
	assert.NotNil(t, manager, manager.userConf, manager.valueMap, manager.confElem)
}

func TestManager_Vars(t *testing.T) {
	manager := NewManager[testConf](nil)
	assert.NotNil(t, manager, manager.userConf, manager.valueMap, manager.confElem)

	confTable := manager.Vars()
	assert.False(t, &confTable == manager.userConf)
}

func TestManager_Raws(t *testing.T) {
	manager := NewManager[testConf](nil)
	assert.NotNil(t, manager, manager.userConf, manager.valueMap, manager.confElem)

	confTable := manager.Raws()
	assert.True(t, confTable == manager.userConf)
}

func TestManager_Conf(t *testing.T) {
	confTable := testConf{
		ServicePort: 1,
	}
	manager := NewManager[testConf](&confTable)

	_, err := manager.Conf("NotExistKey")
	assert.True(t, errors.Is(err, ConfNotExistError))

	conf, err := manager.Conf("ServicePort")
	assert.Nil(t, err)
	err = conf.SetByString("null")
	assert.True(t, errors.Is(err, ConfValueSetError))
	err = conf.SetByString("2")
	assert.Nil(t, err)
	assert.Equal(t, 2, manager.Vars().ServicePort)
}
