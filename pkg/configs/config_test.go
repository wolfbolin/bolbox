package configs

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wolfbolin/bolbox/pkg/errors"
)

func TestManager_Conf(t *testing.T) {
	confTable := flagTestConf{
		IntField: 1,
	}
	manager, err := NewManager[flagTestConf](&confTable)
	assert.Nil(t, err)

	_, err = manager.Conf("NotExistKey")
	assert.True(t, errors.Is(err, ConfNotExistError))

	conf, err := manager.Conf("IntField")
	assert.Nil(t, err)
	err = conf.SetByString("null")
	assert.True(t, errors.Is(err, ConfValueSetError))
	err = conf.SetByString("2")
	assert.Nil(t, err)
	assert.Equal(t, 2, manager.Vars().IntField)
}

func TestConfig_OnChange(t *testing.T) {
	// 测试单个回调函数
	t.Run("单个回调函数", func(t *testing.T) {
		confTable := flagTestConf{
			StringField: "initial",
		}
		manager, err := NewManager[flagTestConf](&confTable)
		assert.Nil(t, err)

		conf, err := manager.Conf("StringField")
		assert.Nil(t, err)

		// 使用 channel 来接收回调通知
		callbackCalled := make(chan any, 1)
		conf.OnChange(func(val any) {
			callbackCalled <- val
		})

		// 修改配置值
		err = conf.SetByValue("new value")
		assert.Nil(t, err)

		// 等待回调被调用
		select {
		case val := <-callbackCalled:
			assert.Equal(t, "new value", val)
		case <-time.After(1 * time.Second):
			t.Fatal("回调函数未被调用")
		}
	})

	// 测试多个回调函数
	t.Run("多个回调函数", func(t *testing.T) {
		confTable := flagTestConf{
			IntField: 10,
		}
		manager, err := NewManager[flagTestConf](&confTable)
		assert.Nil(t, err)

		conf, err := manager.Conf("IntField")
		assert.Nil(t, err)

		// 创建多个 channel 来接收回调通知
		callback1 := make(chan any, 1)
		callback2 := make(chan any, 1)
		callback3 := make(chan any, 1)

		conf.OnChange(func(val any) {
			callback1 <- val
		})
		conf.OnChange(func(val any) {
			callback2 <- val
		})
		conf.OnChange(func(val any) {
			callback3 <- val
		})

		// 修改配置值
		err = conf.SetByValue(int64(20))
		assert.Nil(t, err)

		// 等待所有回调被调用
		select {
		case val := <-callback1:
			assert.Equal(t, int64(20), val)
		case <-time.After(1 * time.Second):
			t.Fatal("回调函数1未被调用")
		}

		select {
		case val := <-callback2:
			assert.Equal(t, int64(20), val)
		case <-time.After(1 * time.Second):
			t.Fatal("回调函数2未被调用")
		}

		select {
		case val := <-callback3:
			assert.Equal(t, int64(20), val)
		case <-time.After(1 * time.Second):
			t.Fatal("回调函数3未被调用")
		}
	})

	// 测试回调函数接收正确的值
	t.Run("回调函数接收正确的值", func(t *testing.T) {
		confTable := flagTestConf{
			BoolField: false,
		}
		manager, err := NewManager[flagTestConf](&confTable)
		assert.Nil(t, err)

		conf, err := manager.Conf("BoolField")
		assert.Nil(t, err)

		receivedValues := make([]any, 0)
		var mu sync.Mutex

		conf.OnChange(func(val any) {
			mu.Lock()
			defer mu.Unlock()
			receivedValues = append(receivedValues, val)
		})

		// 多次修改配置值
		err = conf.SetByValue(true)
		assert.Nil(t, err)
		time.Sleep(10 * time.Millisecond) // 等待 goroutine 执行

		err = conf.SetByValue(false)
		assert.Nil(t, err)
		time.Sleep(10 * time.Millisecond) // 等待 goroutine 执行

		err = conf.SetByValue(true)
		assert.Nil(t, err)
		time.Sleep(10 * time.Millisecond) // 等待 goroutine 执行

		mu.Lock()
		assert.Len(t, receivedValues, 3)
		assert.Equal(t, true, receivedValues[0])
		assert.Equal(t, false, receivedValues[1])
		assert.Equal(t, true, receivedValues[2])
		mu.Unlock()
	})

	// 测试并发注册回调函数
	t.Run("并发注册回调函数", func(t *testing.T) {
		confTable := flagTestConf{
			Float64Field: 1.0,
		}
		manager, err := NewManager[flagTestConf](&confTable)
		assert.Nil(t, err)

		conf, err := manager.Conf("Float64Field")
		assert.Nil(t, err)

		// 并发注册多个回调函数
		var wg sync.WaitGroup
		callbackCount := 10
		callbacks := make([]chan any, callbackCount)

		for i := 0; i < callbackCount; i++ {
			callbacks[i] = make(chan any, 1)
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				conf.OnChange(func(val any) {
					callbacks[idx] <- val
				})
			}(i)
		}

		wg.Wait()

		// 修改配置值
		err = conf.SetByValue(2.5)
		assert.Nil(t, err)

		// 等待所有回调被调用
		for i := 0; i < callbackCount; i++ {
			select {
			case val := <-callbacks[i]:
				assert.Equal(t, 2.5, val)
			case <-time.After(1 * time.Second):
				t.Fatalf("回调函数%d未被调用", i)
			}
		}
	})

	// 测试 SetByString 触发回调
	t.Run("SetByString触发回调", func(t *testing.T) {
		confTable := flagTestConf{
			Int32Field: 100,
		}
		manager, err := NewManager[flagTestConf](&confTable)
		assert.Nil(t, err)

		conf, err := manager.Conf("Int32Field")
		assert.Nil(t, err)

		callbackCalled := make(chan any, 1)
		conf.OnChange(func(val any) {
			callbackCalled <- val
		})

		// 使用 SetByString 修改配置值
		err = conf.SetByString("200")
		assert.Nil(t, err)

		// 等待回调被调用
		select {
		case val := <-callbackCalled:
			// SetByString 会调用 c.val.Interface()，所以返回的是 int64 类型
			assert.Equal(t, int32(200), val)
		case <-time.After(1 * time.Second):
			t.Fatal("回调函数未被调用")
		}
	})
}
