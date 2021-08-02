package sqlm

import "fmt"

type Error struct {
	full  string
	wrong string
	desc  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("full: %s, wrong: %s, desc: %s", e.full, e.wrong, e.desc)
}

type MethodError struct {
	Method string
	Err    error
}

func (e *MethodError) Error() string {
	return fmt.Sprintf("parse method %q failed: %v",
		e.Method, e.Err)
}
