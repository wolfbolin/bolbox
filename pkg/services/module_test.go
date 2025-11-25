package services

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type fakeModule struct {
	name     string
	status   *ModuleStatus
	requires []string
	run      func(ctx context.Context)
}

func (f *fakeModule) Name() string {
	return f.name
}

func (f *fakeModule) Status() *ModuleStatus {
	return f.status
}

func (f *fakeModule) Run(ctx context.Context) {
	if f.run != nil {
		f.run(ctx)
	}
}

func (f *fakeModule) Requires() []string {
	return f.requires
}

func TestManager(t *testing.T) {
	ctx, stop := context.WithCancel(context.TODO())
	defer stop()

	mgr := NewManager()

	go mgr.StartAndServe(ctx)
	res := <-mgr.Done(stop)
	assert.NotNil(t, res)
}

func TestManager_checkAndSort_success(t *testing.T) {
	moduleMap := map[string]Module{
		"A": &fakeModule{
			name:     "A",
			requires: []string{"B", "C", "D"},
		},
		"B": &fakeModule{
			name:     "B",
			requires: []string{"E"},
		},
		"C": &fakeModule{
			name:     "C",
			requires: []string{"D"},
		},
		"D": &fakeModule{
			name:     "D",
			requires: nil,
		},
		"E": &fakeModule{
			name:     "E",
			requires: []string{"C"},
		},
	}
	manager := &Manager{
		moduleMap: moduleMap,
	}
	order, err := manager.checkAndSort()
	assert.Nil(t, err)
	fmt.Printf("Start order %v", order)
	assert.Equal(t, []string{"D", "C", "E", "B", "A"}, order)
}

func TestManager_checkAndSort_cyclic(t *testing.T) {
	moduleMap := map[string]Module{
		"A": &fakeModule{
			name:     "A",
			requires: []string{"B"},
		},
		"B": &fakeModule{
			name:     "B",
			requires: []string{"A", "C"},
		},
		"C": &fakeModule{
			name:     "C",
			requires: []string{},
		},
	}
	manager := &Manager{
		moduleMap: moduleMap,
	}
	_, err := manager.checkAndSort()
	assert.NotNil(t, err)

	moduleMap = map[string]Module{
		"A": &fakeModule{
			name:     "A",
			requires: []string{"B"},
		},
		"B": &fakeModule{
			name:     "B",
			requires: []string{"C"},
		},
		"C": &fakeModule{
			name:     "C",
			requires: []string{"A"},
		},
	}
	manager.moduleMap = moduleMap
	_, err = manager.checkAndSort()
	assert.NotNil(t, err)
}
