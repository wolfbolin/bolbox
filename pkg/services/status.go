package services

import "sync"

// Status 定义状态类型
type Status string

const (
	// StatusRunning 模块正在运行
	StatusRunning Status = "running"
	// StatusStopped 模块停止运行
	StatusStopped Status = "stopped"
)

// ModuleStatus 模块状态类
type ModuleStatus struct {
	lock     sync.Mutex
	status   Status
	syncChan chan Status
}

// NewModuleStatus 新建一个 ModuleStatus 实例
func NewModuleStatus() *ModuleStatus {
	return &ModuleStatus{
		status:   StatusStopped,
		syncChan: make(chan Status),
	}
}

// Set 设置模块状态
func (s *ModuleStatus) Set(status Status) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.status = status
	for {
		select {
		case s.syncChan <- s.status:
		default:
			return // channel is full
		}
	}
}

// Get 获取模块状态
func (s *ModuleStatus) Get() Status {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.status
}

// Watch 监听状态变化
func (s *ModuleStatus) Watch() <-chan Status {
	return s.syncChan
}
