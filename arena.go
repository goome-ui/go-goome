package goome

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"

import "unsafe"

type arena struct {
	chunks []unsafe.Pointer
	sizes  []int
	cur    int
	off    int
}

func (a *arena) alloc(n int) unsafe.Pointer {
	n = (n + 7) &^ 7
	if n == 0 {
		n = 8
	}
	for a.cur < len(a.chunks) {
		if a.off+n <= a.sizes[a.cur] {
			p := unsafe.Add(a.chunks[a.cur], a.off)
			a.off += n
			C.memset(p, 0, C.size_t(n))
			return p
		}
		a.cur++
		a.off = 0
	}
	size := 64 << 10
	if n > size {
		size = n
	}
	p := C.malloc(C.size_t(size))
	if p == nil {
		panic(ErrNoMem)
	}
	a.chunks = append(a.chunks, p)
	a.sizes = append(a.sizes, size)
	a.cur = len(a.chunks) - 1
	a.off = n
	C.memset(p, 0, C.size_t(n))
	return p
}

func (a *arena) str(s string) *C.char {
	p := a.alloc(len(s) + 1)
	if len(s) > 0 {
		C.memcpy(p, unsafe.Pointer(unsafe.StringData(s)), C.size_t(len(s)))
	}
	return (*C.char)(p)
}

func (a *arena) optStr(s string) *C.char {
	if s == "" {
		return nil
	}
	return a.str(s)
}

func (a *arena) reset() {
	for i := 1; i < len(a.chunks); i++ {
		C.free(a.chunks[i])
	}
	if len(a.chunks) > 1 {
		a.chunks, a.sizes = a.chunks[:1], a.sizes[:1]
	}
	a.cur, a.off = 0, 0
}

func (a *arena) free() {
	for _, p := range a.chunks {
		C.free(p)
	}
	a.chunks, a.sizes = nil, nil
	a.cur, a.off = 0, 0
}
