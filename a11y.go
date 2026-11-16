package goome

/*
#include "native/goome.h"
*/
import "C"

import "fmt"

type A11yNode struct {
	Role, Parent, States int
	X, Y, W, H           float64
	Label, Value         string
	ID                   uint32
	Num, Min, Max        float64
	Caret                int
}

func (n A11yNode) Is(state int) bool { return n.States&state != 0 }

func (u *UI) A11y(on bool) { C.goome_ui_a11y(u.h(), b2i(on)) }

func (u *UI) A11yRole(role int, label string) {
	C.goome_ui_a11y_role(u.h(), C.int(role), u.a.str(label))
}

func (u *UI) A11yValue(value string) { C.goome_ui_a11y_value(u.h(), u.a.str(value)) }

func (u *UI) A11yTree() []A11yNode {
	h := u.h()
	n := int(C.goome_a11y_tree(h, nil, 0))
	if n <= 0 {
		return nil
	}
	buf := make([]C.goome_a11y_node, n)
	n = int(C.goome_a11y_tree(h, &buf[0], C.int(n)))
	if n > len(buf) {
		n = len(buf)
	}
	out := make([]A11yNode, n)
	for i, c := range buf[:n] {
		out[i] = A11yNode{
			Role: int(c.role), Parent: int(c.parent), States: int(c.states),
			X: float64(c.x), Y: float64(c.y), W: float64(c.w), H: float64(c.h),
			Label: goStr(c.label[:]), Value: goStr(c.value[:]), ID: uint32(c.id),
			Num: float64(c.num), Min: float64(c.min), Max: float64(c.max), Caret: int(c.caret),
		}
		if out[i].Parent >= n {
			out[i].Parent = -1
		}
	}
	return out
}

func (g *Goome) A11yAction(verb string, n A11yNode, value string) error {
	msg := fmt.Sprintf("a11y %s %d", verb, n.ID)
	if value != "" {
		msg += " " + value
	}
	return g.Push(Event{Type: EventControl, Msg: msg})
}
