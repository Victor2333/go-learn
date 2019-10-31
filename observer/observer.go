package observer

import "sync"

type observer struct {
	stateUpdate func(oldState interface{}, newState interface{})
	updateMutex sync.Mutex
}

func NewObserver(stateUpdate func(oldState interface{}, newState interface{})) *observer {
	return &observer{stateUpdate: stateUpdate}
}

func (o *observer) update(oldState interface{}, newState interface{}) {
	o.updateMutex.Lock()
	defer o.updateMutex.Unlock()
	o.stateUpdate(oldState, newState)
}

func (o *observer) SetStateUpdate(stateUpdate func(oldState interface{}, newState interface{})) {
	o.updateMutex.Lock()
	defer o.updateMutex.Unlock()
	o.stateUpdate = stateUpdate
}
