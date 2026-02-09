package configs

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// flagTestConf 用于测试 configs 功能的配置结构体，包含所有支持的数据类型
type flagTestConf struct {
	// 字符串类型
	StringField string `flag:"string-field" env:"STRING_FIELD" desc:"字符串字段"`

	// 布尔类型
	BoolField bool `flag:"bool-field" desc:"布尔字段"`

	// 整数类型
	IntField   int   `flag:"int-field" desc:"int字段"`
	Int32Field int32 `flag:"int32-field" desc:"int32字段"`
	Int64Field int64 `flag:"int64-field" desc:"int64字段"`

	// 浮点数类型
	Float32Field float32 `flag:"float32-field" desc:"float32字段"`
	Float64Field float64 `flag:"float64-field" desc:"float64字段"`

	MapField map[string]string `flag:"map-field" env:"MAP_FIELD" desc:"map字段"`

	// 没有 flag 标签的字段（不应该被解析）
	NoFlagField string `desc:"没有flag标签的字段"`
}

func TestMain(m *testing.M) {
	originalArgs := os.Args
	os.Args = []string{""}
	defer func() {
		os.Args = originalArgs
	}()
	code := m.Run()
	os.Exit(code)
}

func TestNewManager(t *testing.T) {
	manager, err := NewManager[flagTestConf](nil)
	assert.Nil(t, err)
	assert.NotNil(t, manager, manager.userConf, manager.valueMap, manager.confElem)

	confTable := flagTestConf{}
	manager, err = NewManager[flagTestConf](&confTable)
	assert.Nil(t, err)
	assert.NotNil(t, manager, manager.userConf, manager.valueMap, manager.confElem)
}

func TestManager_Vars(t *testing.T) {
	manager, err := NewManager[flagTestConf](nil)
	assert.Nil(t, err)
	assert.NotNil(t, manager, manager.userConf, manager.valueMap, manager.confElem)

	confTable := manager.Vars()
	assert.False(t, &confTable == manager.userConf)
}

func TestManager_Raws(t *testing.T) {
	manager, err := NewManager[flagTestConf](nil)
	assert.Nil(t, err)
	assert.NotNil(t, manager, manager.userConf, manager.valueMap, manager.confElem)

	confTable := manager.Raws()
	assert.True(t, confTable == manager.userConf)
}
