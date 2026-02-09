package services

import (
	"context"
	"sync"
	"time"

	"github.com/wolfbolin/bolbox/pkg/errors"
	"github.com/wolfbolin/bolbox/pkg/log"
)

// Module 接口定义模块结构，包括获取名称、状态、运行模块和列出依赖的方法。
type Module interface {
	Name() string
	Status() *ModuleStatus
	Run(ctx context.Context)
	Requires() []string
}

// Manager 结构体管理模块的生命周期，包括上下文和取消功能。提供添加、删除和管理模块的方法。
type Manager struct {
	ctx        context.Context
	role       string
	mapLock    sync.RWMutex
	moduleMap  map[string]Module
	contextMap map[string]context.Context
	cancelMap  map[string]context.CancelFunc
}

// NewManager 创建一个新的Manager实例，初始化模块、上下文和取消功能的映射。
func NewManager() *Manager {
	return &Manager{
		role:       "",
		moduleMap:  make(map[string]Module),
		contextMap: make(map[string]context.Context),
		cancelMap:  make(map[string]context.CancelFunc),
	}
}

// StartAndServe 初始化管理器的上下文，锁定模块映射，检查和排序模块启动顺序，并按顺序启动每个模块。处理启动错误并监控状态变化。
func (m *Manager) StartAndServe(ctx context.Context) {
	m.ctx = ctx
	m.mapLock.Lock()
	defer m.mapLock.Unlock()

	order, err := m.checkAndSort()
	if err != nil {
		log.Fatalf("Check for module startup sequence errors. %+v", err)
	}
	log.Infof("Module manager will sequentially start the following modules: %v", order)

	for _, modName := range order {
		log.Infof("Start function module[%s] by order", modName)
		if cancel, ok := m.cancelMap[modName]; ok {
			cancel() // 避免异常退出的模块的协程泄露
		}

		m.contextMap[modName], m.cancelMap[modName] = context.WithCancel(m.ctx)
		modStatus := m.moduleMap[modName].Status()
		if modStatus == nil {
			log.Errorf("Unable to obtain module[%s] status.", modName)
			continue
		}
		go startAndServe(m.contextMap[modName], m.moduleMap[modName])
		select {
		case <-time.After(time.Second):
			log.Fatalf("Module[%s] startup time exceeds expectations", modName)
		case status, ok := <-modStatus.Watch():
			if !ok {
				log.Fatalf("Module[%s] status has not been properly initialized", modName)
			}
			log.Infof("Module[%s] has been switched to status[%s]", modName, status)
		}
	}

	select {
	case <-m.ctx.Done():
		log.Infof("Module manager exit by context")
		return
	}
}

func startAndServe(ctx context.Context, module Module) {
	defer func() {
		if err := recover(); err != nil {
			log.Fatalf("Module[%s] throws a panic during running. %+v", module.Name(), err)
		}
	}()
	time.Sleep(time.Millisecond)
	module.Run(ctx)
}

// AddModule 将新模块添加到管理器的模块映射中，使用锁确保线程安全。
func (m *Manager) AddModule(name string, module Module) {
	m.mapLock.Lock()
	defer m.mapLock.Unlock()
	m.moduleMap[name] = module
}

// DelModule 从管理器的模块映射中删除模块，使用锁确保线程安全。
func (m *Manager) DelModule(name string) {
	m.mapLock.Lock()
	defer m.mapLock.Unlock()
	delete(m.moduleMap, name)
}

// Done 等待所有运行中的模块优雅地退出，使用等待组同步它们的完成。它返回一个信号所有模块退出的通道。
func (m *Manager) Done(stop context.CancelFunc) <-chan struct{} {
	stopCount := 0
	doneChan := make(chan struct{})
	wg := sync.WaitGroup{}
	for _, mod := range m.moduleMap {
		modStatus := mod.Status()
		if modStatus == nil {
			log.Errorf("Unable to obtain module[%s] status.", mod.Name())
			continue
		}
		if modStatus.Get() == StatusRunning {
			log.Infof("Module[%s] is currently running. Notify it to gracefully exit", mod.Name())
			stopCount += 1
		} else {
			continue
		}
		wg.Add(1)
		go func(modName string) {
			<-modStatus.Watch()
			log.Warnf("Module[%s] has gracefully exited", modName)
			wg.Done()
		}(mod.Name())
	}

	log.Infof("Waiting for a total of %d modules to gracefully exit", stopCount)
	go func() {
		wg.Wait()
		doneChan <- struct{}{}
	}()
	stop()
	return doneChan
}

// checkAndSort 检查模块的依赖关系并按照启动顺序进行排序。如果存在循环依赖，则返回错误。
func (m *Manager) checkAndSort() ([]string, error) {
	queue := make(chan string, len(m.moduleMap))
	depNum := make(map[string]int)      // 节点对外依赖的数量
	depMap := make(map[string][]string) // 节点被外部依赖的列表
	for name, module := range m.moduleMap {
		depNum[name] = len(module.Requires())
		if depNum[name] == 0 {
			queue <- name
		}
		for _, depMod := range module.Requires() {
			depMap[depMod] = append(depMap[depMod], name)
		}
	}

	order := make([]string, 0) // 节点启动顺序
	for len(queue) != 0 {
		name := <-queue
		order = append(order, name)

		for _, mod := range depMap[name] {
			depNum[mod] -= 1
			if depNum[mod] <= 0 {
				queue <- mod
			}
		}
	}

	abnormal := make([]string, 0)
	for name, count := range depNum {
		if count > 0 {
			log.Warnf("Module[%s] dependency tree not cleared to zero", name)
			abnormal = append(abnormal, name)
		}
	}
	if len(abnormal) != 0 {
		return nil, errors.Errorf("Module may have cyclic dependencies in [%v]", abnormal)
	}
	return order, nil
}
