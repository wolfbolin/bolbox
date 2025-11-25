package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	modStatus := NewModuleStatus()
	assert.Equal(t, StoppedStatus, modStatus.Get())

	modStatus.Set(RunningStatus)
	assert.Equal(t, RunningStatus, modStatus.Get())

	go func() {
		time.After(100 * time.Millisecond)
		modStatus.Set(StoppedStatus)
	}()

	status := <-modStatus.Watch()
	assert.Equal(t, StoppedStatus, status)
}
