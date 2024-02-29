package pkg

import (
	"context"
	"fmt"
	"sync/atomic"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type Manager interface {
	Go(id string, f func()) error
}

type GoRoutineManager struct {
	ctx                context.Context
	maxGoRoutines      int64
	nameToRoutineCount cmap.ConcurrentMap[string, int64]
	totalCount         int64
	stopCh             <-chan struct{}
}

type routine struct {
	id          string
	routineFunc func()
}

func NewGoRoutineManager(ctx context.Context, maxGoRoutines int64, stopCh <-chan struct{}) Manager {
	return &GoRoutineManager{
		ctx:                ctx,
		maxGoRoutines:      maxGoRoutines,
		nameToRoutineCount: cmap.New[int64](),
		stopCh:             stopCh,
	}
}

func (m *GoRoutineManager) Go(id string, routineFunc func()) error {
	if m.maxGoRoutines <= m.totalCount {
		return fmt.Errorf("max amount of go routine exeeded")
	}
	newRoutine := routine{
		id:          id,
		routineFunc: routineFunc,
	}
	count, ok := m.nameToRoutineCount.Get(id)
	if !ok {
		added := m.nameToRoutineCount.SetIfAbsent(id, 1)
		if !added { // some other thread got to add it first
			count, _ = m.nameToRoutineCount.Get(id) // grab the up-to-date count
		}
	}
	atomic.AddInt64(&count, 1)
	atomic.AddInt64(&m.totalCount, 1)
	go m.runRoutine(newRoutine)
	return nil
}

func (m *GoRoutineManager) runRoutine(routine routine) {
	fmt.Println("starting to run go routine id '", routine.id, "'")
	routine.routineFunc()
	fmt.Println("completed run of go routine id '", routine.id, "'")
	count, _ := m.nameToRoutineCount.Get(routine.id)
	atomic.AddInt64(&count, -1)
	atomic.AddInt64(&m.totalCount, -1)
}
