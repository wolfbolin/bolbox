package services

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	modStatus := NewModuleStatus()
	assert.Equal(t, StatusStopped, modStatus.Get())

	modStatus.Set(StatusRunning)
	assert.Equal(t, StatusRunning, modStatus.Get())

	go func() {
		time.After(100 * time.Millisecond)
		modStatus.Set(StatusStopped)
	}()

	status := <-modStatus.Watch()
	assert.Equal(t, StatusStopped, status)
}
