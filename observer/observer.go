package observer

type observer struct {
	stateUpdate func(s interface{})
}

func NewObserver(stateUpdate func(s interface{})) *observer {
	return &observer{stateUpdate: stateUpdate}
}

func (o *observer) update(s interface{}) {
	o.stateUpdate(s)
}

func (o *observer) SetStateUpdate(stateUpdate func(s interface{})) {
	o.stateUpdate = stateUpdate
}
