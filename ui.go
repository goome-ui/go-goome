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

type Window struct {
	X, Y, W, H float64
	Closed     bool
}

func NewWindow(x, y, w, h float64) *Window { return &Window{X: x, Y: y, W: w, H: h} }

type Split struct {
	Frac     float64
	Vertical bool
}

type Item struct {
	X, Y, W, H                       float64
	Hovered, Pressed, Held, Released bool
	Clicked, DoubleClicked           bool
	RightClicked, Focused            bool
	MX, MY                           float64
	DX, DY                           float64
	DragX, DragY                     float64
	Key                              int
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

type UI struct {
	u       *C.goome_ui
	g       *Goome
	a       arena
	stack   []entry
	wins    map[*Window]*winState
	splits  map[*Split]*C.double
	edits   map[string]*editBuf
	used    map[string]int
	frame   int
	canvas  *Canvas
	Changed bool
	Entered bool
}

var liveUIs []*UI

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

func (u *UI) Theme() *Theme { return (*Theme)(unsafe.Pointer(C.goome_ui_theme(u.h()))) }

func (u *UI) Config() *UIConfig { return (*UIConfig)(unsafe.Pointer(C.goome_ui_cfg(u.h()))) }

func (u *UI) Bind(action int, combos string) error {
	if r := C.goome_ui_bind(u.h(), C.int(action), u.a.str(combos)); r < 0 {
		return newErr(int(r), "cant bind %q", combos)
	}
	return nil
}

func (u *UI) Scale(s float64) { C.goome_ui_scale(u.h(), C.double(s)) }

func (u *UI) Feed(ev Event) bool {
	e := (*C.goome_event)(u.a.alloc(int(unsafe.Sizeof(C.goome_event{}))))
	*e = ev.c()
	return C.goome_ui_feed(u.h(), e) > 0
}

func (u *UI) push(name string, end func()) { u.stack = append(u.stack, entry{name, end}) }

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

func (u *UI) Depth() int { return len(u.stack) }

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
	for k, e := range u.edits {
		if e.seen < u.frame-600 {
			C.free(unsafe.Pointer(e.buf))
			delete(u.edits, k)
		}
	}
}

func (u *UI) Frame(mon int, draw func(c *Canvas)) {
	c := u.Begin(mon)
	defer u.End()
	if c != nil {
		draw(c)
	}
}

func (u *UI) Hit(x, y, w, h float64) {
	C.goome_ui_hit(u.h(), C.double(x), C.double(y), C.double(w), C.double(h))
}

func (u *UI) Animating() bool  { return C.goome_ui_animating(u.h()) > 0 }
func (u *UI) Wants() int       { return int(C.goome_ui_wants(u.h())) }
func (u *UI) WantsMouse() bool { return u.Wants()&WantMouse != 0 }
func (u *UI) WantsKeyboard() bool {
	return u.Wants()&WantKeyboard != 0
}
func (u *UI) AutoGrab(wants int) { C.goome_ui_auto_grab(u.h(), C.int(wants)) }

func (u *UI) scope() string {
	var b strings.Builder
	for _, e := range u.stack {
		b.WriteString(e.name)
		b.WriteByte(0)
	}
	return b.String()
}

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

func (u *UI) Row(cols int)        { C.goome_ui_row(u.h(), C.int(cols)) }
func (u *UI) SameLine()           { C.goome_ui_same_line(u.h()) }
func (u *UI) Space(px float64)    { C.goome_ui_space(u.h(), C.double(px)) }
func (u *UI) Separator()          { C.goome_ui_separator(u.h()) }
func (u *UI) Width(w float64)     { C.goome_ui_width(u.h(), C.double(w)) }
func (u *UI) Pos(x, y float64)    { C.goome_ui_pos(u.h(), C.double(x), C.double(y)) }
func (u *UI) Size(w, h float64)   { C.goome_ui_size(u.h(), C.double(w), C.double(h)) }
func (u *UI) Label(text string)   { C.ggo_label(u.h(), u.a.str(text)) }
func (u *UI) Heading(text string) { C.goome_ui_heading(u.h(), u.a.str(text)) }
func (u *UI) Wrapped(text string) { C.goome_ui_wrapped(u.h(), u.a.str(text)) }

func (u *UI) changed(r C.int) bool {
	u.Changed = r != 0
	return u.Changed
}

func (u *UI) Button(label string) bool {
	return u.changed(C.goome_ui_button(u.h(), u.a.str(label)))
}

func (u *UI) ImageButton(id string, b *Bitmap, w, h float64) bool {
	return u.changed(C.goome_ui_image_button(u.h(), u.a.str(id), b.ptr(), C.double(w), C.double(h)))
}

func (u *UI) Image(b *Bitmap, w, h, radius float64) {
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

func (u *UI) Checkbox(label string, v *bool) bool {
	return u.flag(func(h *C.goome_ui, l *C.char, p *C.int) C.int { return C.goome_ui_checkbox(h, l, p) }, label, v)
}

func (u *UI) Toggle(label string, v *bool) bool {
	return u.flag(func(h *C.goome_ui, l *C.char, p *C.int) C.int { return C.goome_ui_toggle(h, l, p) }, label, v)
}

func (u *UI) Radio(label string, v *int, value int) bool {
	p := u.cint(*v)
	r := u.changed(C.goome_ui_radio(u.h(), u.a.str(label), p, C.int(value)))
	*v = int(*p)
	return r
}

func (u *UI) Slider(label string, v *float64, min, max float64) bool {
	return u.SliderFmt(label, v, min, max, "")
}

func (u *UI) SliderFmt(label string, v *float64, min, max float64, format string) bool {
	p := u.cdouble(*v)
	r := u.changed(C.goome_ui_slider(u.h(), u.a.str(label), p, C.double(min), C.double(max), u.a.optStr(format)))
	*v = float64(*p)
	return r
}

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

func (u *UI) Input(label string, s *string) bool {
	changed, _ := u.InputEx(label, s, 0, "", 256)
	return changed
}

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

func (u *UI) Combo(label string, sel *int, items []string) bool {
	arr, n := u.strs(items)
	p := u.cint(*sel)
	r := u.changed(C.goome_ui_combo(u.h(), u.a.str(label), p, arr, n))
	*sel = int(*p)
	return r
}

func (u *UI) Selectable(label string, selected bool) bool {
	return u.changed(C.goome_ui_selectable(u.h(), u.a.str(label), b2i(selected)))
}

func (u *UI) Tabs(id string, names []string, sel *int) bool {
	arr, n := u.strs(names)
	p := u.cint(*sel)
	r := u.changed(C.goome_ui_tabs(u.h(), u.a.str(id), arr, n, p))
	*sel = int(*p)
	return r
}

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

func (u *UI) Progress(frac float64, text string) {
	C.goome_ui_progress(u.h(), C.double(frac), u.a.optStr(text))
}

func (u *UI) Color(label string, argb *uint32) bool {
	p := (*C.uint32_t)(u.a.alloc(4))
	*p = C.uint32_t(*argb)
	r := u.changed(C.goome_ui_color(u.h(), u.a.str(label), p))
	*argb = uint32(*p)
	return r
}

func (u *UI) Tooltip(text string) { C.goome_ui_tooltip(u.h(), u.a.str(text)) }
func (u *UI) Hovered() bool       { return C.goome_ui_hovered(u.h()) != 0 }

func (u *UI) Scroll(id string, h float64) bool {
	ui := u.h()
	r := C.goome_ui_scroll_begin(ui, u.a.str(id), C.double(h))
	u.push("s:"+id, func() { C.goome_ui_scroll_end(ui) })
	return r != 0
}

func (u *UI) ID(id string) {
	ui := u.h()
	C.goome_ui_push_id(ui, u.a.str(id))
	u.push("i:"+id, func() { C.goome_ui_pop_id(ui) })
}

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

func (u *UI) Menubar() bool {
	ui := u.h()
	return u.maybe("menubar", C.goome_ui_menubar_begin(ui), func() { C.goome_ui_menubar_end(ui) })
}

func (u *UI) Menu(label string) bool {
	ui := u.h()
	return u.maybe("m:"+label, C.goome_ui_menu_begin(ui, u.a.str(label)), func() { C.goome_ui_menu_end(ui) })
}

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

func (u *UI) MenuSeparator() { C.goome_ui_menu_separator(u.h()) }

func (u *UI) Context(id string) bool {
	ui := u.h()
	return u.maybe("c:"+id, C.goome_ui_context_begin(ui, u.a.str(id)), func() { C.goome_ui_context_end(ui) })
}

func (u *UI) PopupOpen(id string) { C.goome_ui_popup_open(u.h(), u.a.str(id)) }

func (u *UI) Popup(id string) bool {
	ui := u.h()
	return u.maybe("p:"+id, C.goome_ui_popup_begin(ui, u.a.str(id)), func() { C.goome_ui_popup_end(ui) })
}

func (u *UI) PopupClose() { C.goome_ui_popup_close(u.h()) }

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

func (u *UI) Toast(text string, ms int) { C.goome_ui_toast(u.h(), u.a.str(text), C.int(ms)) }
func (u *UI) Statusbar(text string)     { C.goome_ui_statusbar(u.h(), u.a.str(text)) }
func (u *UI) Focused() bool             { return C.goome_ui_focused(u.h()) != 0 }
func (u *UI) SetFocus()                 { C.goome_ui_set_focus(u.h()) }

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

func (u *UI) Table(id string, cols int, h float64, flags int) bool {
	ui := u.h()
	r := C.goome_ui_table_begin(ui, u.a.str(id), C.int(cols), C.double(h), C.int(flags))
	u.push("t:"+id, func() { C.goome_ui_table_end(ui) })
	return r != 0
}

func (u *UI) TableColumn(name string, width float64) {
	C.goome_ui_table_column(u.h(), u.a.str(name), C.double(width))
}

func (u *UI) TableSort(col *int, desc *bool) bool {
	c, d := u.cint(*col), u.cbool(*desc)
	r := C.goome_ui_table_sort(u.h(), c, d) != 0
	*col, *desc = int(*c), *d != 0
	return r
}

func (u *UI) TableRow()  { C.goome_ui_table_row(u.h()) }
func (u *UI) TableNext() { C.goome_ui_table_next(u.h()) }

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

func (u *UI) ListClip(count int, rowH float64) (first, last int) {
	f, l := u.cint(0), u.cint(0)
	C.goome_ui_list_clip(u.h(), C.int(count), C.double(rowH), f, l)
	return int(*f), int(*l)
}

func (u *UI) DragFloat(label string, v *float64, speed, min, max float64, format string) bool {
	p := u.cdouble(*v)
	r := u.changed(C.goome_ui_drag_float(u.h(), u.a.str(label), p, C.double(speed), C.double(min), C.double(max), u.a.optStr(format)))
	*v = float64(*p)
	return r
}

func (u *UI) DragInt(label string, v *int, speed float64, min, max int) bool {
	p := u.cint(*v)
	r := u.changed(C.goome_ui_drag_int(u.h(), u.a.str(label), p, C.double(speed), C.int(min), C.int(max)))
	*v = int(*p)
	return r
}

func (u *UI) Number(label string, v *float64, step float64, format string) bool {
	p := u.cdouble(*v)
	r := u.changed(C.goome_ui_number(u.h(), u.a.str(label), p, C.double(step), u.a.optStr(format)))
	*v = float64(*p)
	return r
}

func (u *UI) floats(v []float32) *C.float {
	p := (*C.float)(u.a.alloc(len(v) * 4))
	copy(unsafe.Slice((*float32)(unsafe.Pointer(p)), len(v)), v)
	return p
}

func (u *UI) PlotLines(label string, v []float32, min, max, h float64) {
	C.goome_ui_plot_lines(u.h(), u.a.str(label), u.floats(v), C.int(len(v)), C.double(min), C.double(max), C.double(h))
}

func (u *UI) PlotBars(label string, v []float32, min, max, h float64) {
	C.goome_ui_plot_bars(u.h(), u.a.str(label), u.floats(v), C.int(len(v)), C.double(min), C.double(max), C.double(h))
}

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

func (u *UI) SplitNext() { C.goome_ui_split_next(u.h()) }

func (u *UI) DragSource(typ string, data []byte) bool {
	p := u.a.alloc(len(data))
	if len(data) > 0 {
		C.memcpy(p, unsafe.Pointer(&data[0]), C.size_t(len(data)))
	}
	return C.goome_ui_drag_source(u.h(), u.a.str(typ), p, C.int(len(data))) != 0
}

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

func (u *UI) Spinner(size float64) { C.goome_ui_spinner(u.h(), C.double(size)) }

func (u *UI) Link(text string) bool { return u.changed(C.goome_ui_link(u.h(), u.a.str(text))) }

func (u *UI) Badge(text string, argb uint32) {
	C.goome_ui_badge(u.h(), u.a.str(text), C.uint32_t(argb))
}

func (u *UI) PushColor(what int, argb uint32) {
	C.goome_ui_push_color(u.h(), C.int(what), C.uint32_t(argb))
}
func (u *UI) PushNum(what int, v float64) { C.goome_ui_push_num(u.h(), C.int(what), C.double(v)) }
func (u *UI) PushRadii(r Radii)           { C.goome_ui_push_radii(u.h(), r.c()) }

func (u *UI) PushFont(f *Font) {
	if f != nil {
		pin(unsafe.Pointer(f.raw()), f)
	}
	C.goome_ui_push_font(u.h(), f.raw())
}

func (u *UI) PushSkin(kind int, s *Skin) {
	if s == nil {
		C.goome_ui_push_skin(u.h(), C.int(kind), nil)
		return
	}
	p := (*Skin)(u.a.alloc(int(unsafe.Sizeof(Skin{}))))
	*p = *s
	C.goome_ui_push_skin(u.h(), C.int(kind), (*C.goome_skin)(unsafe.Pointer(p)))
}

func (u *UI) PushTheme(t *Theme) {
	p := (*Theme)(u.a.alloc(int(unsafe.Sizeof(Theme{}))))
	*p = *t
	C.goome_ui_push_theme(u.h(), p.c())
}

func (u *UI) PushIcon(icon int, b *Bitmap) {
	var p *C.goome_bitmap
	if b != nil {
		p = b.ptr()
		pin(unsafe.Pointer(p), b)
	}
	C.goome_ui_push_icon(u.h(), C.int(icon), p)
}

func (u *UI) Pop(n int) { C.goome_ui_pop(u.h(), C.int(n)) }

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

func (u *UI) Canvas() *Canvas {
	c := C.goome_ui_canvas(u.h())
	if c == nil {
		return nil
	}
	return &Canvas{c: c, g: u.g, gen: u.g.frames}
}

func (u *UI) Store(id string, n int) []float64 {
	p := C.goome_ui_store(u.h(), u.a.str(id), C.int(n))
	if p == nil {
		return nil
	}
	return unsafe.Slice((*float64)(unsafe.Pointer(p)), n)
}
