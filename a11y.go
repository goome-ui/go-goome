package goome

/*
#include "native/goome.h"
*/
import "C"

import "fmt"

// A11yNode is one entry in the accessibility tree, the same thing a screen reader sees.
type A11yNode struct {
	// Parent is the index of the parent node in the tree, -1 for the top.
	// States holds the state bits, test them with Is.
	// Role is what the node is, one of the Role* constants.
	Role, Parent, States int
	// H is the height in logical coordinates.
	// Y is the top edge in logical coordinates.
	// W is the width in logical coordinates.
	// X is the left edge in logical coordinates.
	X, Y, W, H float64
	// Value is the text value, like the contents of an input.
	// Label is the accessible name.
	Label, Value string
	// ID identifies the node, use it with Goome.A11yAction.
	ID uint32
	// Min is the bottom of the range for sliders and progress bars.
	// Num is the numeric value of a slider or progress bar.
	// Max is the top of the range for sliders and progress bars.
	Num, Min, Max float64
	// Caret is the caret position in an input.
	Caret int
}

// Is reports whether the given state bit (checked, focused, expanded and so on) is set.
func (n A11yNode) Is(state int) bool { return n.States&state != 0 }

// A11y turns accessibility on or off. It is on by default in a window and off for overlays and offscreen.
func (u *UI) A11y(on bool) { C.goome_ui_a11y(u.h(), b2i(on)) }

// A11yRole gives the next custom item a role and a label.
func (u *UI) A11yRole(role int, label string) {
	C.goome_ui_a11y_role(u.h(), C.int(role), u.a.str(label))
}

// A11yValue sets the accessible value of the next custom item.
func (u *UI) A11yValue(value string) { C.goome_ui_a11y_value(u.h(), u.a.str(value)) }

// A11yTree returns the accessibility tree built from the last frame, only for monitor 0.
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

// A11yAction does what a screen reader would: verb is press, focus, expand, collapse or set (with value) on node n.
func (g *Goome) A11yAction(verb string, n A11yNode, value string) error {
	msg := fmt.Sprintf("a11y %s %d", verb, n.ID)
	if value != "" {
		msg += " " + value
	}
	return g.Push(Event{Type: EventControl, Msg: msg})
}
