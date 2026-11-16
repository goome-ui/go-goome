package goome

/*
#include "native/goome.h"
*/
import "C"

import "fmt"

type Error struct {
	Code int
	Msg  string
}

func (e *Error) Error() string {
	if e.Msg == "" {
		return "goome: " + Strerror(e.Code)
	}
	return "goome: " + e.Msg
}

func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	return ok && t.Msg == "" && t.Code == e.Code
}

var (
	ErrNotSupported = &Error{Code: CodeNotSupported}
	ErrDenied       = &Error{Code: CodeDenied}
	ErrNoDisplay    = &Error{Code: CodeNoDisplay}
	ErrGone         = &Error{Code: CodeGone}
	ErrAgain        = &Error{Code: CodeAgain}
	ErrNoMem        = &Error{Code: CodeNoMem}
	ErrArg          = &Error{Code: CodeArg}
)

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
