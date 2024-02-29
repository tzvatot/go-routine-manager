package pkg

import (
	"context"
	"fmt"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type Manager interface {
	Go(id string, f func()) error
}

type GoRoutineManager struct {
	ctx           context.Context
	maxGoRoutines int
	nameToRoutine cmap.ConcurrentMap[string, []routine]
	stopCh        <-chan struct{}
}

type routine struct {
	id          string
	routineFunc func()
}

func NewGoRoutineManager(ctx context.Context, maxGoRoutines int, stopCh <-chan struct{}) Manager {
	return &GoRoutineManager{
		ctx:           ctx,
		maxGoRoutines: maxGoRoutines,
		nameToRoutine: cmap.New[[]routine](),
		stopCh:        stopCh,
	}
}

func (m *GoRoutineManager) Go(id string, routineFunc func()) error {
	newRoutine := routine{
		id:          id,
		routineFunc: routineFunc,
	}
	routines, ok := m.nameToRoutine.Get(id)
	if !ok {
		routines = append(make([]routine, 0), newRoutine)
	} else {
		routines = append(routines, newRoutine)
	}
	m.nameToRoutine.Set(id, routines)
	go m.runRoutine(newRoutine)
	return nil
}

func (m *GoRoutineManager) runRoutine(routine routine) {
	fmt.Println("starting to run go routine id '", routine.id, "'")
	routine.routineFunc()
	fmt.Println("completed run of go routine id '", routine.id, "'")
}
