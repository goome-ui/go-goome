package goome

/*
#include "native/goome.h"
*/
import "C"

import "fmt"

// Error is the error type goome calls return, holding the C error code and the message from the handle.
type Error struct {
	// Code is the goome error code, one of the Code constants.
	Code int
	// Msg is the detail message from goome. When it is empty, Error falls back to the text for Code.
	Msg string
}

// Error returns the message prefixed with goome.
func (e *Error) Error() string {
	if e.Msg == "" {
		return "goome: " + Strerror(e.Code)
	}
	return "goome: " + e.Msg
}

// Is reports whether target is a bare Err value with the same code, so errors.Is works with ErrDenied and friends.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && t.Msg == "" && t.Code == e.Code
}

var (
	// ErrNotSupported is the sentinel error for something this backend cannot do, for use with errors.Is.
	ErrNotSupported = &Error{Code: CodeNotSupported}
	// ErrDenied is the sentinel error for a call that was not allowed, for use with errors.Is.
	ErrDenied = &Error{Code: CodeDenied}
	// ErrNoDisplay is the sentinel error for no display to draw on, for use with errors.Is.
	ErrNoDisplay = &Error{Code: CodeNoDisplay}
	// ErrGone is the sentinel error for a monitor that is gone, for use with errors.Is.
	ErrGone = &Error{Code: CodeGone}
	// ErrAgain is the sentinel error for when nothing was ready yet, for use with errors.Is.
	ErrAgain = &Error{Code: CodeAgain}
	// ErrNoMem is the sentinel error for running out of memory, for use with errors.Is.
	ErrNoMem = &Error{Code: CodeNoMem}
	// ErrArg is the sentinel error for a bad argument, for use with errors.Is.
	ErrArg = &Error{Code: CodeArg}
)

// Strerror turns an error code into a short readable message.
func Strerror(code int) string {
	return C.GoString(C.goome_strerror(C.int(code)))
}

func newErr(code int, format string, a ...any) error {
	return &Error{Code: code, Msg: fmt.Sprintf(format, a...)}
}

func (g *Goome) check(r C.int) error {
	if r >= 0 {
		return nil
	}
	msg := ""
	if g != nil && g.g != nil {
		msg = C.GoString(C.goome_error(g.g))
	}
	if msg == "" {
		msg = Strerror(int(r))
	}
	return &Error{Code: int(r), Msg: msg}
}
