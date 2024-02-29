package pkg

import (
	"fmt"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type Manager interface {
	Go(id string, f func()) error
}

type GoRoutineManager struct {
	maxGoRoutines int
	nameToRoutine cmap.ConcurrentMap[string, []routine]
}

type routine struct {
	id          string
	routineFunc func()
}

func NewGoRoutineManager(maxGoRoutines int) Manager {
	return &GoRoutineManager{
		maxGoRoutines: maxGoRoutines,
		nameToRoutine: cmap.New[[]routine](),
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
	go func() {
		fmt.Println("starting to run go routine id '", id, "'")
		newRoutine.routineFunc()
		fmt.Println("completed run of go routine id '", id, "'")
	}()
	return nil
}
