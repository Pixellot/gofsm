package gofsm

// Event represents an event name and possible tranitions.
type Event struct {
    Name        string
    transitions []*Transition
}

// To sets an event target state.
func (e *Event) To(name string) *Transition {
    transition := &Transition{to: name}
    e.transitions = append(e.transitions, transition)
    return transition
}

// Transition represents a tarnsition source, targets, and related callbacks.
type Transition struct {
    to      string
    froms   []string
    befores []func(value interface{}) error
    afters  []func(value interface{}) error
}

// From sets an event possible source states
func (t *Transition) From(states ...string) *Transition {
    t.froms = states
    return t
}

// Before registers a possible callback to run before the tarnsition.
func (t *Transition) Before(fn func(value interface{}) error) *Transition {
    t.befores = append(t.befores, fn)
    return t
}

// After is a possible callback to run after the transition.
func (t *Transition) After(fn func(value interface{}) error) *Transition {
    t.afters = append(t.afters, fn)
    return t
}

