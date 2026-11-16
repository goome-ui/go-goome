package goome

/*
#include <stdint.h>
*/
import "C"

import "runtime/cgo"

//export goomeGoEase
func goomeGoEase(t C.double, user C.uintptr_t) C.double {
	f, ok := cgo.Handle(user).Value().(func(float64) float64)
	if !ok {
		return t
	}
	return C.double(f(float64(t)))
}
