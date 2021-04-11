// Package gofsm provides simple functionality for FSM
package gofsm

import (
    "fmt"
    "strings"
    "time"
)

// TransitionRecord represents an info on a single FSM transtion.
type TransitionRecord struct {
    When    time.Time
    From    string
    To      string
    Note    string
}


// New initialize a new FSM.
func NewFSM() *FSM {
    return &FSM{
        states: map[string]*State{},
        events: map[string]*Event{},
    }
}

// A FSM represents an FSM with its initial state, states, and events.
type FSM struct {
    initial     string
    states      map[string]*State
    events      map[string]*Event
}

// Initial sets the FSM initial state.
func (sm *FSM) Initial(name string) *FSM {
    sm.initial = name
    return sm
}

// State adds a state name to the FSM.
func (sm *FSM) State(name string) *State {
    state := &State{Name: name}
    sm.states[name] = state
    return state
}

// Event adds an event name to the FSM.
func (sm *FSM) Event(name string) *Event {
    event := &Event{Name: name}
    sm.events[name] = event
    return event
}

// Trigger performs a transition according to event name on a stated type, and returns the transition info.
func (sm *FSM) Trigger(name string, value Stater, desc ...string) (*TransitionRecord, error) {

    current := value.GetState()
    if current == "" {
        current = sm.initial
        value.SetState(sm.initial)
    }

    event := sm.events[name]
    if event == nil {
        return nil, fmt.Errorf("failed to perform event %s from state %s, no such event", name, current)
    }

    var matches []*Transition
    for _, transition := range event.transitions {
        valid := len(transition.froms) == 0
        for _, from := range transition.froms {
            if from == current {
                valid = true
            }
        }

        if valid {
            matches = append(matches, transition)
        }
    }

    if len(matches) != 1 {
        return nil, fmt.Errorf("failed to perform event %s from state %s, invalid number of transitions (%d) in event", name, current, len(matches))
    }

    transition := matches[0]

    if err := sm.exitFromCurrentState(value, current); err != nil {
        return nil, err
    }

    if err := sm.performActionsBeforeTransition(value, transition); err != nil {
        return nil, err
    }

    value.SetState(transition.to)
    previous := current

    if err := sm.performActionsAfterTransition(value, previous, transition); err != nil {
        return nil, err
    }

    if err := sm.enterNextState(value, previous, transition.to); err != nil {
        return nil, err
    }

    return &TransitionRecord{When: time.Now().UTC(), From: previous, To: transition.to, Note: strings.Join(desc, "")}, nil
}

func (sm *FSM) exitFromCurrentState(value Stater, current string) error {
    state, ok := sm.states[current]
    if !ok {
        return nil
    }
    for _, exit := range state.exits {
        if err := exit(value); err != nil {
            return err
        }
    }
    return nil
}

func (sm *FSM) performActionsBeforeTransition(value Stater, transition *Transition) error {
    for _, before := range transition.befores {
        if err := before(value); err != nil {
            return err
        }
    }
    return nil
}

func (sm *FSM) enterNextState(value Stater, previous, next string) error {
    state, ok := sm.states[next]
    if !ok {
        return nil
    }
    for _, enter := range state.enters {
        if err := enter(value); err != nil {
            value.SetState(previous)
            return err
        }
    }
    return nil
}

func (sm *FSM) performActionsAfterTransition(value Stater, previous string, transition *Transition) error {
    for _, after := range transition.afters {
        if err := after(value); err != nil {
            value.SetState(previous)
            return err
        }
    }
    return nil
}
