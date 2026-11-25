package goome

/*
#include <stdlib.h>
#include <string.h>
#include "native/goome.h"
#include "helpers.h"
*/
import "C"

import (
	"strings"
	"unsafe"
)

// Window holds the position, size and open state of a floating window that UI.Window reads and updates.
type Window struct {
	// X is the window position on the X axis, in logical units.
	// W is the window width, in logical units.
	// Y is the window position on the Y axis, in logical units.
	// H is the window height, in logical units.
	X, Y, W, H float64
	// Closed is true once the window has been closed.
	Closed bool
}

// NewWindow returns a Window at the given position and size, in logical units.
func NewWindow(x, y, w, h float64) *Window { return &Window{X: x, Y: y, W: w, H: h} }

// Split holds the state of a draggable split between two panes, which UI.Split reads and updates.
type Split struct {
	// Frac is how much of the space the first pane takes, from 0 to 1.
	Frac float64
	// Vertical splits the panes with a vertical divider when true.
	Vertical bool
}

// Item is what UI.Custom reports back about a custom widget: where it ended up and what the pointer and keyboard did to it this frame.
type Item struct {
	// X is where the item ended up on the X axis, in logical units.
	// W is the width the item ended up with, in logical units.
	// Y is where the item ended up on the Y axis, in logical units.
	// H is the height the item ended up with, in logical units.
	X, Y, W, H float64
	// Pressed is true on the frame the pointer button went down on the item.
	// Released is true on the frame the pointer button was let go.
	// Held is true while the pointer button is down on the item.
	// Hovered is true when the pointer is over the item.
	Hovered, Pressed, Held, Released bool
	// Clicked is true when the item was clicked this frame.
	// DoubleClicked is true when the item was double clicked this frame.
	Clicked, DoubleClicked bool
	// RightClicked is true when the item was right clicked this frame.
	// Focused is true when the item has keyboard focus.
	RightClicked, Focused bool
	// MY is the pointer Y position, in logical units.
	// MX is the pointer X position, in logical units.
	MX, MY float64
	// DY is how far the pointer moved on Y this frame while the item was held, after the axis lock.
	// DX is how far the pointer moved on X this frame while the item was held, after the axis lock.
	DX, DY float64
	// DragY is how far the pointer has moved on Y since the press, after the axis lock.
	// DragX is how far the pointer has moved on X since the press, after the axis lock.
	DragX, DragY float64
	// Key is the key pressed this frame while the item had focus, with modifiers like shift ored in, or 0 for none.
	Key int
}

type winState struct {
	x, y, w, h C.double
	open       C.int
}

type editBuf struct {
	buf  *C.char
	cap  int
	seen int
}

type entry struct {
	name string
	end  func()
}

// UI is the immediate mode interface. Feed it events, then draw every frame between Begin and End.
type UI struct {
	u      *C.goome_ui
	g      *Goome
	a      arena
	stack  []entry
	wins   map[*Window]*winState
	splits map[*Split]*C.double
	edits  map[string]*editBuf
	used   map[string]int
	frame  int
	canvas *Canvas
	held   []*Bitmap
	// Changed is true if the last widget changed its value.
	Changed bool
	// Entered is true if the last text input had Enter pressed in it.
	Entered bool
}

var liveUIs []*UI

// NewUI creates a UI bound to this Goome, or returns ErrGone if the Goome is already closed. Free it when you are done.
func (g *Goome) NewUI() (*UI, error) {
	if g.Closed() {
		return nil, ErrGone
	}
	p := C.goome_ui_new(g.g)
	if p == nil {
		return nil, newErr(CodeNoMem, "no memory for the ui")
	}
	u := &UI{u: p, g: g, wins: map[*Window]*winState{}, splits: map[*Split]*C.double{}, edits: map[string]*editBuf{}, used: map[string]int{}}
	g.uis = append(g.uis, u)
	liveUIs = append(liveUIs, u)
	return u, nil
}

// Free releases the UI and everything it holds. It is safe to call more than once.
func (u *UI) Free() {
	if u == nil || u.u == nil {
		return
	}
	u.Config().release()
	C.goome_ui_free(u.u)
	u.u = nil
	u.stack = nil
	for _, w := range u.wins {
		C.free(unsafe.Pointer(w))
	}
	for _, s := range u.splits {
		C.free(unsafe.Pointer(s))
	}
	for _, e := range u.edits {
		C.free(unsafe.Pointer(e.buf))
	}
	u.wins, u.splits, u.edits = nil, nil, nil
	u.a.free()
	for i, x := range u.g.uis {
		if x == u {
			u.g.uis = append(u.g.uis[:i], u.g.uis[i+1:]...)
			break
		}
	}
	for i, x := range liveUIs {
		if x == u {
			liveUIs = append(liveUIs[:i], liveUIs[i+1:]...)
			break
		}
	}
}

// Freed reports whether the UI was already freed.
func (u *UI) Freed() bool { return u == nil || u.u == nil }

func (u *UI) h() *C.goome_ui {
	if u == nil || u.u == nil {
		panic(&Error{Code: CodeGone, Msg: "this ui is freed"})
	}
	return u.u
}

func (u *UI) forgetFont(f *Font) {
	if u.u == nil {
		return
	}
	u.Theme().forget(unsafe.Pointer(f.f), nil)
}

func (u *UI) forgetBitmap(b *Bitmap) {
	if u.u == nil {
		return
	}
	u.Theme().forget(nil, unsafe.Pointer(b.b))
}

// Theme returns the UI's current theme, which you can change in place.
func (u *UI) Theme() *Theme { return (*Theme)(unsafe.Pointer(C.goome_ui_theme(u.h()))) }

// Config returns the UI's tuning settings, like double click timing, cursors and key bindings. Changes apply right away.
func (u *UI) Config() *UIConfig { return (*UIConfig)(unsafe.Pointer(C.goome_ui_cfg(u.h()))) }

// Bind sets the key combos for an editing action such as copy or undo. It returns an error if the action or the combos are not valid.
func (u *UI) Bind(action int, combos string) error {
	if r := C.goome_ui_bind(u.h(), C.int(action), u.a.str(combos)); r < 0 {
		return newErr(int(r), "cant bind %q", combos)
	}
	return nil
}

// Scale sets the UI scale factor.
func (u *UI) Scale(s float64) { C.goome_ui_scale(u.h(), C.double(s)) }

// Feed hands an event to the UI and returns true if the UI used it.
func (u *UI) Feed(ev Event) bool {
	e := (*C.goome_event)(u.a.alloc(int(unsafe.Sizeof(C.goome_event{}))))
	*e = ev.c()
	return C.goome_ui_feed(u.h(), e) > 0
}

func (u *UI) push(name string, end func()) { u.stack = append(u.stack, entry{name, end}) }

// End closes the most recent scope you opened, such as a window, scroll area, popup or the frame itself. It does nothing if none is open.
func (u *UI) End() {
	if len(u.stack) == 0 {
		return
	}
	e := u.stack[len(u.stack)-1]
	u.stack = u.stack[:len(u.stack)-1]
	if e.end != nil {
		e.end()
	}
}

// Depth returns how many scopes are open right now.
func (u *UI) Depth() int { return len(u.stack) }

// Begin starts a frame on monitor mon and returns the canvas to draw on, or nil if there is nothing to draw on and you should skip the frame. Anything still open from before is closed first.
func (u *UI) Begin(mon int) *Canvas {
	h := u.h()
	for len(u.stack) > 0 {
		u.End()
	}
	u.a.reset()
	u.frame++
	u.used = map[string]int{}
	c := C.goome_ui_begin(h, C.int(mon))
	u.g.frames++
	u.push("frame", u.endFrame)
	if c == nil {
		u.canvas = nil
		return nil
	}
	u.canvas = &Canvas{c: c, g: u.g, gen: u.g.frames}
	return u.canvas
}

func (u *UI) endFrame() {
	C.goome_ui_end(u.h())
	if u.canvas != nil {
		u.canvas.c = nil
		u.canvas = nil
	}
	u.held = nil
	for k, e := range u.edits {
		if e.seen < u.frame-600 {
			C.free(unsafe.Pointer(e.buf))
			delete(u.edits, k)
		}
	}
}

// Frame runs draw with the canvas for monitor mon and closes the frame afterwards. If there is nothing to draw on, draw is skipped.
func (u *UI) Frame(mon int, draw func(c *Canvas)) {
	c := u.Begin(mon)
	defer u.End()
	if c != nil {
		draw(c)
	}
}

// Hit adds a hit area to the frame at x, y with size w by h, in logical units.
func (u *UI) Hit(x, y, w, h float64) {
	C.goome_ui_hit(u.h(), C.double(x), C.double(y), C.double(w), C.double(h))
}

// Animating reports whether something is still moving, so you know to draw again soon.
func (u *UI) Animating() bool { return C.goome_ui_animating(u.h()) > 0 }

// Anim eases a number toward target over ms milliseconds and returns the current value, keyed by id. ease picks the easing curve.
func (u *UI) Anim(id string, target, ms float64, ease int) float64 {
	return float64(C.goome_ui_anim(u.h(), u.a.str(id), C.double(target), C.double(ms), C.int(ease)))
}

// AnimColor eases an ARGB color toward target over ms milliseconds and returns the current color, keyed by id. ease picks the easing curve.
func (u *UI) AnimColor(id string, target uint32, ms float64, ease int) uint32 {
	return uint32(C.goome_ui_anim_color(u.h(), u.a.str(id), C.uint32_t(target), C.double(ms), C.int(ease)))
}

// Wants returns the WantMouse and WantKeyboard bits for what the UI wants right now, good after Feed and End.
func (u *UI) Wants() int { return int(C.goome_ui_wants(u.h())) }

// WantsMouse reports whether the UI wants the mouse right now.
func (u *UI) WantsMouse() bool { return u.Wants()&WantMouse != 0 }

// WantsKeyboard reports whether the UI wants the keyboard right now.
func (u *UI) WantsKeyboard() bool {
	return u.Wants()&WantKeyboard != 0
}

// AutoGrab sets what the UI takes from the OS by itself, as WantMouse and WantKeyboard bits. The default is both, and 0 means the app does it.
func (u *UI) AutoGrab(wants int) { C.goome_ui_auto_grab(u.h(), C.int(wants)) }

func (u *UI) scope() string {
	var b strings.Builder
	for _, e := range u.stack {
		b.WriteString(e.name)
		b.WriteByte(0)
	}
	return b.String()
}

// Window draws a floating window with the given title, position and size taken from w, and returns true while it is open. The matching End closes it.
func (u *UI) Window(title string, w *Window, flags int) bool {
	h := u.h()
	s := u.wins[w]
	if s == nil {
		s = (*winState)(C.calloc(1, C.size_t(unsafe.Sizeof(winState{}))))
		u.wins[w] = s
	}
	s.x, s.y, s.w, s.h, s.open = C.double(w.X), C.double(w.Y), C.double(w.W), C.double(w.H), b2i(!w.Closed)
	back := func() {
		w.X, w.Y, w.W, w.H, w.Closed = float64(s.x), float64(s.y), float64(s.w), float64(s.h), s.open == 0
	}
	r := C.goome_ui_window(h, u.a.str(title), &s.x, &s.y, &s.w, &s.h, C.int(flags), &s.open)
	back()
	u.push("w:"+title, func() {
		C.goome_ui_window_end(h)
		back()
	})
	return r != 0
}

// Row lays the next widgets out in cols columns.
func (u *UI) Row(cols int) { C.goome_ui_row(u.h(), C.int(cols)) }

// SameLine puts the next widget on the same line as the last one.
func (u *UI) SameLine() { C.goome_ui_same_line(u.h()) }

// Space adds px of empty space.
func (u *UI) Space(px float64) { C.goome_ui_space(u.h(), C.double(px)) }

// Separator draws a horizontal divider line.
func (u *UI) Separator() { C.goome_ui_separator(u.h()) }

// Width sets the width of the next widget, where 0 fills the space.
func (u *UI) Width(w float64) { C.goome_ui_width(u.h(), C.double(w)) }

// Pos puts the next widget at x, y from the top left of the content, scroll included. The normal flow doesn't move.
func (u *UI) Pos(x, y float64) { C.goome_ui_pos(u.h(), C.double(x), C.double(y)) }

// Size sets the size of the next widget, where 0 keeps its own size for that side.
func (u *UI) Size(w, h float64) { C.goome_ui_size(u.h(), C.double(w), C.double(h)) }

// Label draws a line of text.
func (u *UI) Label(text string) { C.ggo_label(u.h(), u.a.str(text)) }

// Heading draws text in the heading style.
func (u *UI) Heading(text string) { C.goome_ui_heading(u.h(), u.a.str(text)) }

// Wrapped draws text that wraps to the available width.
func (u *UI) Wrapped(text string) { C.goome_ui_wrapped(u.h(), u.a.str(text)) }

func (u *UI) changed(r C.int) bool {
	u.Changed = r != 0
	return u.Changed
}

// Button draws a button and returns true when it was clicked.
func (u *UI) Button(label string) bool {
	return u.changed(C.goome_ui_button(u.h(), u.a.str(label)))
}

// ImageButton draws a button showing the bitmap at w by h and returns true when it was clicked. id keeps it apart from other buttons.
func (u *UI) ImageButton(id string, b *Bitmap, w, h float64) bool {
	u.held = append(u.held, b)
	return u.changed(C.goome_ui_image_button(u.h(), u.a.str(id), b.ptr(), C.double(w), C.double(h)))
}

// Image draws the bitmap at w by h with rounded corners of the given radius.
func (u *UI) Image(b *Bitmap, w, h, radius float64) {
	u.held = append(u.held, b)
	C.goome_ui_image(u.h(), b.ptr(), C.double(w), C.double(h), C.double(radius))
}

func (u *UI) cint(v int) *C.int {
	p := (*C.int)(u.a.alloc(4))
	*p = C.int(v)
	return p
}

func (u *UI) cbool(v bool) *C.int { return u.cint(int(b2i(v))) }

func (u *UI) cdouble(v float64) *C.double {
	p := (*C.double)(u.a.alloc(8))
	*p = C.double(v)
	return p
}

func (u *UI) flag(fn func(*C.goome_ui, *C.char, *C.int) C.int, label string, v *bool) bool {
	p := u.cbool(*v)
	r := u.changed(fn(u.h(), u.a.str(label), p))
	*v = *p != 0
	return r
}

// Checkbox draws a checkbox bound to v and returns true when it was toggled.
func (u *UI) Checkbox(label string, v *bool) bool {
	return u.flag(func(h *C.goome_ui, l *C.char, p *C.int) C.int { return C.goome_ui_checkbox(h, l, p) }, label, v)
}

// Toggle draws an on/off switch bound to v and returns true when it was flipped.
func (u *UI) Toggle(label string, v *bool) bool {
	return u.flag(func(h *C.goome_ui, l *C.char, p *C.int) C.int { return C.goome_ui_toggle(h, l, p) }, label, v)
}

// Radio draws a radio button that sets v to value when picked, and returns true when it was.
func (u *UI) Radio(label string, v *int, value int) bool {
	p := u.cint(*v)
	r := u.changed(C.goome_ui_radio(u.h(), u.a.str(label), p, C.int(value)))
	*v = int(*p)
	return r
}

func numFmt(format string) string {
	convs := 0
	for i := 0; i < len(format); i++ {
		if format[i] != '%' {
			continue
		}
		i++
		if i < len(format) && format[i] == '%' {
			continue
		}
		for i < len(format) && strings.IndexByte("-+ #0", format[i]) >= 0 {
			i++
		}
		for i < len(format) && format[i] >= '0' && format[i] <= '9' {
			i++
		}
		if i < len(format) && format[i] == '.' {
			i++
			for i < len(format) && format[i] >= '0' && format[i] <= '9' {
				i++
			}
		}
		if convs++; i >= len(format) || strings.IndexByte("fFeEgG", format[i]) < 0 || convs > 1 {
			return ""
		}
	}
	return format
}

// Slider draws a slider for v between min and max and returns true when it changed.
func (u *UI) Slider(label string, v *float64, min, max float64) bool {
	return u.SliderFmt(label, v, min, max, "")
}

// SliderFmt is Slider with a printf style format for the shown value.
func (u *UI) SliderFmt(label string, v *float64, min, max float64, format string) bool {
	p := u.cdouble(*v)
	r := u.changed(C.goome_ui_slider(u.h(), u.a.str(label), p, C.double(min), C.double(max), u.a.optStr(numFmt(format))))
	*v = float64(*p)
	return r
}

// SliderInt draws an integer slider for v between min and max and returns true when it changed.
func (u *UI) SliderInt(label string, v *int, min, max int) bool {
	p := u.cint(*v)
	r := u.changed(C.goome_ui_slider_int(u.h(), u.a.str(label), p, C.int(min), C.int(max)))
	*v = int(*p)
	return r
}

func (u *UI) textBuf(s string, cap int) *C.char {
	if cap < 1 {
		cap = 256
	}
	if len(s) > cap-1 {
		s = s[:cap-1]
	}
	p := (*C.char)(u.a.alloc(cap))
	if len(s) > 0 {
		C.memcpy(unsafe.Pointer(p), unsafe.Pointer(unsafe.StringData(s)), C.size_t(len(s)))
	}
	return p
}

// Input draws a single line text input editing s and returns true when the text changed.
func (u *UI) Input(label string, s *string) bool {
	changed, _ := u.InputEx(label, s, 0, "", 256)
	return changed
}

// InputEx is Input with flags, a hint shown while empty, and a cap on the buffer size in bytes (256 if cap is below 1). It returns whether the text changed and whether Enter was pressed.
func (u *UI) InputEx(label string, s *string, flags int, hint string, cap int) (changed, entered bool) {
	if cap < 1 {
		cap = 256
	}
	buf := u.textBuf(*s, cap)
	r := C.goome_ui_input_ex(u.h(), u.a.str(label), buf, C.int(cap), C.int(flags), u.a.optStr(hint))
	*s = C.GoString(buf)
	u.Changed, u.Entered = r != 0, r == 2
	return u.Changed, u.Entered
}

// Password draws a text input that hides what you type, with a hint shown while empty. It returns true when the text changed.
func (u *UI) Password(label string, s *string, hint string) bool {
	changed, _ := u.InputEx(label, s, InputPassword, hint, 256)
	return changed
}

func (u *UI) strs(items []string) (**C.char, C.int) {
	n := len(items)
	arr := (**C.char)(u.a.alloc((n + 1) * int(unsafe.Sizeof(uintptr(0)))))
	s := unsafe.Slice(arr, n+1)
	for i, it := range items {
		s[i] = u.a.str(it)
	}
	return arr, C.int(n)
}

// Combo draws a dropdown over items with the chosen index in sel, and returns true when the choice changed.
func (u *UI) Combo(label string, sel *int, items []string) bool {
	arr, n := u.strs(items)
	p := u.cint(*sel)
	r := u.changed(C.goome_ui_combo(u.h(), u.a.str(label), p, arr, n))
	*sel = int(*p)
	return r
}

// Selectable draws a row that can be selected and returns true when it was clicked.
func (u *UI) Selectable(label string, selected bool) bool {
	return u.changed(C.goome_ui_selectable(u.h(), u.a.str(label), b2i(selected)))
}

// Tabs draws a tab bar over names with the chosen index in sel, and returns true when it changed.
func (u *UI) Tabs(id string, names []string, sel *int) bool {
	arr, n := u.strs(names)
	p := u.cint(*sel)
	r := u.changed(C.goome_ui_tabs(u.h(), u.a.str(id), arr, n, p))
	*sel = int(*p)
	return r
}

// Collapsing draws a collapsible header and returns true while it is open. open can be nil, otherwise it holds the open state.
func (u *UI) Collapsing(label string, open *bool) bool {
	if open == nil {
		u.Changed = false
		return C.goome_ui_collapsing(u.h(), u.a.str(label), nil) != 0
	}
	p := u.cbool(*open)
	r := C.goome_ui_collapsing(u.h(), u.a.str(label), p) != 0
	u.Changed = (*p != 0) != *open
	*open = *p != 0
	return r
}

// Progress draws a progress bar filled to frac, from 0 to 1, with optional text on it.
func (u *UI) Progress(frac float64, text string) {
	C.goome_ui_progress(u.h(), C.double(frac), u.a.optStr(text))
}

// Color draws a color picker for the ARGB value and returns true when it changed.
func (u *UI) Color(label string, argb *uint32) bool {
	p := (*C.uint32_t)(u.a.alloc(4))
	*p = C.uint32_t(*argb)
	r := u.changed(C.goome_ui_color(u.h(), u.a.str(label), p))
	*argb = uint32(*p)
	return r
}

// Tooltip sets the tooltip text for the widget right before it.
func (u *UI) Tooltip(text string) { C.goome_ui_tooltip(u.h(), u.a.str(text)) }

// Hovered returns whether the pointer is over the widget right before it.
func (u *UI) Hovered() bool { return C.goome_ui_hovered(u.h()) != 0 }

// Scroll starts a scroll area h tall and returns true if it is visible. The matching End closes it.
func (u *UI) Scroll(id string, h float64) bool {
	ui := u.h()
	r := C.goome_ui_scroll_begin(ui, u.a.str(id), C.double(h))
	u.push("s:"+id, func() { C.goome_ui_scroll_end(ui) })
	return r != 0
}

// ID pushes id onto the id stack, so repeated labels in a loop stay apart. Push the index first, and the matching End pops it.
func (u *UI) ID(id string) {
	ui := u.h()
	C.goome_ui_push_id(ui, u.a.str(id))
	u.push("i:"+id, func() { C.goome_ui_pop_id(ui) })
}

// App opens a root panel that fills monitor 0, call it after Begin. flags picks the titlebar, menubar, statusbar and padding. It returns whether the panel is shown, and the matching End closes it.
func (u *UI) App(title string, flags int) bool {
	ui := u.h()
	r := C.goome_ui_app(ui, u.a.optStr(title), C.int(flags))
	u.push("app", func() { C.goome_ui_app_end(ui) })
	return r != 0
}

func (u *UI) maybe(name string, r C.int, end func()) bool {
	if r != 0 {
		u.push(name, end)
	} else {
		u.push(name, nil)
	}
	return r != 0
}

// Menubar starts the menu bar and returns true while it is shown. The matching End closes it.
func (u *UI) Menubar() bool {
	ui := u.h()
	return u.maybe("menubar", C.goome_ui_menubar_begin(ui), func() { C.goome_ui_menubar_end(ui) })
}

// Menu opens a menu with the given label and returns true while it is open. Nest them for submenus, and the matching End closes it.
func (u *UI) Menu(label string) bool {
	ui := u.h()
	return u.maybe("m:"+label, C.goome_ui_menu_begin(ui, u.a.str(label)), func() { C.goome_ui_menu_end(ui) })
}

// MenuItem draws a menu entry with an optional shortcut label and returns true when it was clicked. checked can be nil, otherwise it shows and holds a check mark.
func (u *UI) MenuItem(label, shortcut string, checked *bool) bool {
	var p *C.int
	if checked != nil {
		p = u.cbool(*checked)
	}
	r := u.changed(C.goome_ui_menu_item(u.h(), u.a.str(label), u.a.optStr(shortcut), p))
	if checked != nil {
		*checked = *p != 0
	}
	return r
}

// MenuSeparator draws a divider line between menu entries.
func (u *UI) MenuSeparator() { C.goome_ui_menu_separator(u.h()) }

// Context opens a context menu on a right click on the widget before it, or on empty space. It returns true while the menu is open, and the matching End closes it.
func (u *UI) Context(id string) bool {
	ui := u.h()
	return u.maybe("c:"+id, C.goome_ui_context_begin(ui, u.a.str(id)), func() { C.goome_ui_context_end(ui) })
}

// PopupOpen opens the popup with this id.
func (u *UI) PopupOpen(id string) { C.goome_ui_popup_open(u.h(), u.a.str(id)) }

// Popup returns true while the popup with this id is open, and the matching End closes it. Open it with PopupOpen.
func (u *UI) Popup(id string) bool {
	ui := u.h()
	return u.maybe("p:"+id, C.goome_ui_popup_begin(ui, u.a.str(id)), func() { C.goome_ui_popup_end(ui) })
}

// PopupClose closes the innermost open popup, menu or modal.
func (u *UI) PopupClose() { C.goome_ui_popup_close(u.h()) }

// Modal opens a centered dialog that dims everything else, and Esc closes it. It returns true while it is open, open can be nil, and the matching End closes it.
func (u *UI) Modal(title string, open *bool) bool {
	ui := u.h()
	var p *C.int
	if open != nil {
		p = u.cbool(*open)
	}
	r := C.goome_ui_modal_begin(ui, u.a.str(title), p)
	back := func() {
		if open != nil {
			*open = *p != 0
		}
	}
	back()
	return u.maybe("d:"+title, r, func() {
		C.goome_ui_modal_end(ui)
		back()
	})
}

// Toast pops up a short message for ms milliseconds.
func (u *UI) Toast(text string, ms int) { C.goome_ui_toast(u.h(), u.a.str(text), C.int(ms)) }

// Statusbar sets the status bar text, which shows with the AppStatusbar flag.
func (u *UI) Statusbar(text string) { C.goome_ui_statusbar(u.h(), u.a.str(text)) }

// Focused returns whether the widget right before it has keyboard focus.
func (u *UI) Focused() bool { return C.goome_ui_focused(u.h()) != 0 }

// SetFocus gives keyboard focus to the next widget.
func (u *UI) SetFocus() { C.goome_ui_set_focus(u.h()) }

// TextEdit draws a multiline text editor h tall that edits text and returns true when it changed. flags picks things like read only, line numbers, wrapping and a monospace font.
func (u *UI) TextEdit(id string, text *string, h float64, flags int) bool {
	key := u.scope() + id
	u.used[key]++
	key += string(rune('0' + u.used[key]))
	e := u.edits[key]
	need := len(*text)*4 + 1
	if need < 65536 {
		need = 65536
	}
	if e == nil || len(*text)*4 > e.cap*3 {
		if e != nil {
			C.free(unsafe.Pointer(e.buf))
		}
		e = &editBuf{buf: (*C.char)(C.calloc(C.size_t(need), 1)), cap: need}
		u.edits[key] = e
	}
	e.seen = u.frame
	if C.GoString(e.buf) != *text {
		cs := unsafe.Slice((*byte)(unsafe.Pointer(e.buf)), e.cap)
		n := copy(cs[:e.cap-1], *text)
		cs[n] = 0
	}
	r := u.changed(C.goome_ui_text_edit(u.h(), u.a.str(id), e.buf, C.int(e.cap), C.double(h), C.int(flags)))
	*text = C.GoString(e.buf)
	return r
}

// Table starts a table with cols columns and h tall, using the table flags. The matching End closes it.
func (u *UI) Table(id string, cols int, h float64, flags int) bool {
	ui := u.h()
	r := C.goome_ui_table_begin(ui, u.a.str(id), C.int(cols), C.double(h), C.int(flags))
	u.push("t:"+id, func() { C.goome_ui_table_end(ui) })
	return r != 0
}

// TableColumn adds a column with a header name and a width, where 0 shares the rest of the space.
func (u *UI) TableColumn(name string, width float64) {
	C.goome_ui_table_column(u.h(), u.a.str(name), C.double(width))
}

// TableSort reports whether the sort changed this frame, and updates col and desc. Call it after the columns.
func (u *UI) TableSort(col *int, desc *bool) bool {
	c, d := u.cint(*col), u.cbool(*desc)
	r := C.goome_ui_table_sort(u.h(), c, d) != 0
	*col, *desc = int(*c), *d != 0
	return r
}

// TableRow starts a new row.
func (u *UI) TableRow() { C.goome_ui_table_row(u.h()) }

// TableNext moves to the next cell of the row.
func (u *UI) TableNext() { C.goome_ui_table_next(u.h()) }

// Tree draws a tree node and returns true while it is open. open can be nil, otherwise it holds the open state, and the matching End closes it.
func (u *UI) Tree(label string, open *bool) bool {
	ui := u.h()
	var p *C.int
	if open != nil {
		p = u.cbool(*open)
	}
	r := C.goome_ui_tree(ui, u.a.str(label), p)
	if open != nil {
		*open = *p != 0
	}
	return u.maybe("r:"+label, r, func() { C.goome_ui_tree_end(ui) })
}

// ListClip works out which rows of a count row list at rowH each are visible in a scroll area, so you only draw first through last.
func (u *UI) ListClip(count int, rowH float64) (first, last int) {
	f, l := u.cint(0), u.cint(0)
	C.goome_ui_list_clip(u.h(), C.int(count), C.double(rowH), f, l)
	return int(*f), int(*l)
}

// DragFloat draws a number you change by dragging, moving speed per pixel and clamped to min and max. format is an optional printf style format.
func (u *UI) DragFloat(label string, v *float64, speed, min, max float64, format string) bool {
	p := u.cdouble(*v)
	r := u.changed(C.goome_ui_drag_float(u.h(), u.a.str(label), p, C.double(speed), C.double(min), C.double(max), u.a.optStr(numFmt(format))))
	*v = float64(*p)
	return r
}

// DragInt draws an integer you change by dragging, moving speed per pixel and clamped to min and max.
func (u *UI) DragInt(label string, v *int, speed float64, min, max int) bool {
	p := u.cint(*v)
	r := u.changed(C.goome_ui_drag_int(u.h(), u.a.str(label), p, C.double(speed), C.int(min), C.int(max)))
	*v = int(*p)
	return r
}

// Number draws a number field with step buttons that change v by step. format is an optional printf style format.
func (u *UI) Number(label string, v *float64, step float64, format string) bool {
	p := u.cdouble(*v)
	r := u.changed(C.goome_ui_number(u.h(), u.a.str(label), p, C.double(step), u.a.optStr(numFmt(format))))
	*v = float64(*p)
	return r
}

func (u *UI) floats(v []float32) *C.float {
	p := (*C.float)(u.a.alloc(len(v) * 4))
	copy(unsafe.Slice((*float32)(unsafe.Pointer(p)), len(v)), v)
	return p
}

// PlotLines draws v as a line graph with the given label, scaled between min and max, h tall.
func (u *UI) PlotLines(label string, v []float32, min, max, h float64) {
	C.goome_ui_plot_lines(u.h(), u.a.str(label), u.floats(v), C.int(len(v)), C.double(min), C.double(max), C.double(h))
}

// PlotBars draws v as a bar chart with the given label, scaled between min and max, h tall.
func (u *UI) PlotBars(label string, v []float32, min, max, h float64) {
	C.goome_ui_plot_bars(u.h(), u.a.str(label), u.floats(v), C.int(len(v)), C.double(min), C.double(max), C.double(h))
}

// Split lays out two panes divided by a draggable bar, and returns true if it is shown. Call SplitNext to move to the second pane, and End when done.
func (u *UI) Split(id string, s *Split) bool {
	ui := u.h()
	f := u.splits[s]
	if f == nil {
		f = (*C.double)(C.calloc(1, 8))
		u.splits[s] = f
	}
	*f = C.double(s.Frac)
	r := C.goome_ui_split_begin(ui, u.a.str(id), b2i(s.Vertical), f)
	s.Frac = float64(*f)
	u.push("x:"+id, func() {
		C.goome_ui_split_end(ui)
		s.Frac = float64(*f)
	})
	return r != 0
}

// SplitNext moves on to the second pane of the current split.
func (u *UI) SplitNext() { C.goome_ui_split_next(u.h()) }

// DragSource makes the widget before it draggable, carrying data under the given type. It returns true while it is being dragged.
func (u *UI) DragSource(typ string, data []byte) bool {
	p := u.a.alloc(len(data))
	if len(data) > 0 {
		C.memcpy(p, unsafe.Pointer(&data[0]), C.size_t(len(data)))
	}
	return C.goome_ui_drag_source(u.h(), u.a.str(typ), p, C.int(len(data))) != 0
}

// DropTarget makes the widget before it accept drops of the given type and returns the bytes dropped, or nil if nothing was. cap is the largest drop it takes, and 4096 is used if cap is below 1.
func (u *UI) DropTarget(typ string, cap int) []byte {
	if cap < 1 {
		cap = 4096
	}
	p := u.a.alloc(cap)
	n := int(C.goome_ui_drop_target(u.h(), u.a.str(typ), p, C.int(cap)))
	if n <= 0 {
		return nil
	}
	if n > cap {
		n = cap
	}
	return C.GoBytes(p, C.int(n))
}

// Spinner draws a loading spinner of the given size.
func (u *UI) Spinner(size float64) { C.goome_ui_spinner(u.h(), C.double(size)) }

// Link draws a link and returns true when it was clicked.
func (u *UI) Link(text string) bool { return u.changed(C.goome_ui_link(u.h(), u.a.str(text))) }

// Badge draws a small label filled with the ARGB color.
func (u *UI) Badge(text string, argb uint32) {
	C.goome_ui_badge(u.h(), u.a.str(text), C.uint32_t(argb))
}

// PushColor overrides one of the theme's color fields with an ARGB value until it is popped.
func (u *UI) PushColor(what int, argb uint32) {
	C.goome_ui_push_color(u.h(), C.int(what), C.uint32_t(argb))
}

// PushNum overrides one of the theme's number fields, such as the corner radius, until it is popped.
func (u *UI) PushNum(what int, v float64) { C.goome_ui_push_num(u.h(), C.int(what), C.double(v)) }

// PushRadii sets a radius per corner until it is popped, and beats the single radius number.
func (u *UI) PushRadii(r Radii) { C.goome_ui_push_radii(u.h(), r.c()) }

// PushFont switches to the font until it is popped. nil means the built in font.
func (u *UI) PushFont(f *Font) {
	if f != nil {
		pin(unsafe.Pointer(f.raw()), f)
	}
	C.goome_ui_push_font(u.h(), f.raw())
}

// PushSkin sets the skin for the given widget kind until it is popped. nil means no skin.
func (u *UI) PushSkin(kind int, s *Skin) {
	if s == nil {
		C.goome_ui_push_skin(u.h(), C.int(kind), nil)
		return
	}
	p := (*Skin)(u.a.alloc(int(unsafe.Sizeof(Skin{}))))
	*p = *s
	C.goome_ui_push_skin(u.h(), C.int(kind), (*C.goome_skin)(unsafe.Pointer(p)))
}

// PushTheme swaps in a whole theme with one push until it is popped.
func (u *UI) PushTheme(t *Theme) {
	p := (*Theme)(u.a.alloc(int(unsafe.Sizeof(Theme{}))))
	*p = *t
	C.goome_ui_push_theme(u.h(), p.c())
}

// PushIcon sets the bitmap used for the given icon until it is popped.
func (u *UI) PushIcon(icon int, b *Bitmap) {
	var p *C.goome_bitmap
	if b != nil {
		p = b.ptr()
		pin(unsafe.Pointer(p), b)
	}
	C.goome_ui_push_icon(u.h(), C.int(icon), p)
}

// Pop removes the last n pushes of any kind, such as colors, fonts or skins.
func (u *UI) Pop(n int) { C.goome_ui_pop(u.h(), C.int(n)) }

// Custom reserves a w by h area and reports what happened to it as an Item. Setting only DragX or DragY in flags locks dragging to that axis. The bool is true when something happened.
func (u *UI) Custom(id string, w, h float64, flags int) (Item, bool) {
	it := (*C.goome_ui_item)(u.a.alloc(int(unsafe.Sizeof(C.goome_ui_item{}))))
	r := C.goome_ui_custom(u.h(), u.a.str(id), C.double(w), C.double(h), C.int(flags), it)
	return Item{
		X: float64(it.x), Y: float64(it.y), W: float64(it.w), H: float64(it.h),
		Hovered: it.hovered != 0, Pressed: it.pressed != 0, Held: it.held != 0, Released: it.released != 0,
		Clicked: it.clicked != 0, DoubleClicked: it.double_clicked != 0, RightClicked: it.right_clicked != 0, Focused: it.focused != 0,
		MX: float64(it.mx), MY: float64(it.my), DX: float64(it.dx), DY: float64(it.dy),
		DragX: float64(it.drag_x), DragY: float64(it.drag_y), Key: int(it.key),
	}, r != 0
}

// Canvas returns the canvas for the spot you are at, so drawing on it lands in the UI and is clipped like the widgets. It returns nil if there is none.
func (u *UI) Canvas() *Canvas {
	c := C.goome_ui_canvas(u.h())
	if c == nil {
		return nil
	}
	return &Canvas{c: c, g: u.g, gen: u.g.frames}
}

// Store returns n float64s kept for this id between frames, all zero the first time.
func (u *UI) Store(id string, n int) []float64 {
	if n <= 0 || n > 1<<24 {
		return nil
	}
	p := C.goome_ui_store(u.h(), u.a.str(id), C.int(n))
	if p == nil {
		return nil
	}
	return unsafe.Slice((*float64)(unsafe.Pointer(p)), n)
}
