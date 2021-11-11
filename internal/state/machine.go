package state

import (
	"errors"
	"sync"
)

type State string

type StateMachine struct {
	states map[State]Event
	lock   *sync.Mutex
}

type Event func(state State) State

func NewStateMachine() *StateMachine {
	sm := make(map[State]Event)
	return &StateMachine{states: sm,
		lock: &sync.Mutex{},
	}
}

type StateEvent struct {
	State
	Event
}

func (s *StateMachine) RegisterState(stateEvents ...StateEvent) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, stateEvent := range stateEvents {
		if stateEvent.State != "" && stateEvent.Event != nil {
			s.states[stateEvent.State] = stateEvent.Event
		}
	}
}

func (s *StateMachine) ProcessState(state State) (State, error) {
	event, ok := s.states[state]
	if !ok {
		return "", errors.New("provided state not valid")
	}
	next := event(state)
	return next, nil
}
