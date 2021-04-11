package gofsm

// StatedValue represents an object that can switch states, and it should be embedded in a struct in order to make it tansitionalble by FSM.
type StatedValue struct {
    State string
}

// SetState sets the object state to name.
func (sv *StatedValue) SetState(name string) {
    sv.State = name
}

// GetState returns the object atate.
func (sv StatedValue) GetState() string {
    return sv.State
}

// Stater represents a stated type, and should be used only as a function argument.
type Stater interface {
    SetState(name string)
    GetState() string
}

