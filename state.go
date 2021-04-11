package gofsm

// State represents a state name, and possible callbacks.
type State struct {
    Name    string
    enters  []func(value interface{}) error
    exits   []func(value interface{}) error
}

// Enter registers a possible callback to run when entering a state.
func (s *State) Enter(fn func(value interface{}) error) *State {
    s.enters = append(s.enters, fn)
    return s
}

// Exit registers a possible callback to run when exiting a state.
func (s *State) Exit(fn func(value interface{}) error) *State {
    s.exits = append(s.exits, fn)
    return s
}

