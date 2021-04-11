package gofsm

import (
    "fmt"
    "bytes"
    "sort"
)

// Outputs a visualization of a FSM in Graphviz format.
func Visualize(sm *FSM) string {
	var buf bytes.Buffer

	writeHeader(&buf)
	writeTransitions(&buf, sm)
	writeStates(&buf, sm)
	writeFooter(&buf)

	return buf.String()
}

func writeHeader(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintf(`digraph fsm {`))
	buf.WriteString("\n")
}

func writeTransitions(buf *bytes.Buffer, sm *FSM) {
    keys := make([]string, 0, len(sm.events))
    for k := range sm.events {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    for _, k := range keys {
        e := sm.events[k]
        ts := e.transitions
        for _, t := range ts {
            buf.WriteString(fmt.Sprintf(`    "%s" -> "%s" [ label = "%s" ];`, t.froms, t.to, e.Name))
            buf.WriteString("\n")
        }
    }

	buf.WriteString("\n")
}

func writeStates(buf *bytes.Buffer, sm *FSM) {
    keys := make([]string, 0, len(sm.states))
    for k := range sm.states {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    buf.WriteString(fmt.Sprintf(`    "%s (initial);"`, sm.initial))
    buf.WriteString("\n")
	for _, k := range keys {
		buf.WriteString(fmt.Sprintf(`    "%s";`, k))
		buf.WriteString("\n")
	}
}

func writeFooter(buf *bytes.Buffer) {
	buf.WriteString(fmt.Sprintln("}"))
}
