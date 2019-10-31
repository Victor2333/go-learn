package observer

import "sync"

type subject struct {
	observers []*observer
	setMutex  sync.Mutex
	state     interface{}
	oldState  interface{}
}

func NewSubject(state interface{}) *subject {
	return &subject{observers: []*observer{}, state: state, oldState: nil}
}

func (s *subject) GetState() interface{} {
	return s.state
}

func (s *subject) SetState(state interface{}) {
	s.setMutex.Lock()
	defer s.setMutex.Unlock()
	s.oldState = s.state
	s.state = state
	s.notify()
}

func (s *subject) Attach(observer ...*observer) {
	s.observers = append(s.observers, observer...)
}

func (s *subject) notify() {
	for _, observer := range s.observers {
		observer.update(s.oldState, s.state)
	}
}
