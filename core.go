// Package goome is the Go binding for the goome C library: overlays on top of everything, normal desktop windows or ui inside your own game. Open a session with Open, OpenApp or OpenOffscreen, then draw with Canvas or the UI kit.
package goome

/*
#include <stdlib.h>
#include "native/goome.h"
*/
import "C"

import (
	"errors"
	"iter"
	"math"
	"unsafe"
)

// Options are the settings for Open. Leave a field zero for its default.
type Options struct {
	// Title is the window title, empty means "goome".
	Title string
	// Backend forces a backend, BackendAuto picks one. The GOOME_BACKEND env var wins over it.
	Backend int
	// StateDir is the folder for the portal token and the control fifo. Empty means neither.
	StateDir string
	// NoCapture means an overlay never asks for screen access.
	NoCapture bool
	// CaptureFPS is the frame rate of the GNOME and portal capture streams, 0 means 5.
	CaptureFPS int
	// HideSelf keeps our overlay out of grabs where the OS lets us.
	HideSelf bool
	// TestFrames is for the test backend only, it can also come from the GOOME_TEST_FRAMES env var.
	TestFrames string
	// TestScript is for the test backend only, it can also come from the GOOME_TEST_SCRIPT env var.
	TestScript string
	// TestShots is for the test backend only, it can also come from the GOOME_TEST_SHOTS env var.
	TestShots string
	// Mode is one of ModeOverlay, ModeWindow or ModeOffscreen.
	Mode int
	// Width is the logical window or offscreen width, 0 means 960.
	Width int
	// Height is the logical window or offscreen height, 0 means 640.
	Height int
	// MinWidth is the smallest window width, 0 means no limit. Window mode only.
	MinWidth int
	// MinHeight is the smallest window height, 0 means no limit. Window mode only.
	MinHeight int
	// WindowFlags are the Window* bits ored together, like WindowResizable. Window mode only.
	WindowFlags int
	// Scale is real pixels per logical pixel, 0 means 1. Offscreen only.
	Scale float64
}

// Rect is a rectangle.
type Rect struct {
	// W is the width.
	// H is the height.
	// X is the left edge.
	// Y is the top edge.
	X, Y, W, H int
}

func (r Rect) c() C.goome_rect {
	return C.goome_rect{x: C.int(r.X), y: C.int(r.Y), w: C.int(r.W), h: C.int(r.H)}
}

func rectFrom(r C.goome_rect) Rect {
	return Rect{int(r.x), int(r.y), int(r.w), int(r.h)}
}

// Monitor describes one screen.
type Monitor struct {
	// ID is the monitor number you pass to Begin and Capture.
	ID int
	// Name is the name the system gives the monitor.
	Name string
	// Area is the monitor rect in logical desktop coordinates.
	Area Rect
	// Scale is real pixels per logical pixel.
	Scale float64
	// PH is the height in real pixels.
	// PW is the width in real pixels.
	PW, PH int
	// Rotation is 0, 90, 180 or 270, plus 360 if mirrored. Grabs and canvases already undo it.
	Rotation int
}

// Event is one thing that happened, like a click, key press, hotkey, drop or window event.
type Event struct {
	// Type is which kind of event this is, one of the Event* constants.
	Type int
	// X is the logical x position of the pointer.
	// Y is the logical y position of the pointer.
	X, Y float64
	// Mon is the monitor the event belongs to.
	Mon int
	// Button is the mouse button, 1 left, 2 middle, 3 right.
	Button int
	// DY is the scroll amount in notches, +1 is down.
	DY int
	// Key is the lowercase ascii key or one of the Key* constants.
	Key int
	// Mods holds the modifier bits (Shift, Ctrl, Alt, Super) that were down.
	Mods int
	// Text is the utf-8 the key typed, can be empty.
	Text string
	// ID depends on the type: the hotkey id, 1 or 0 for focus gained or lost, or the caret byte offset for a preedit.
	ID int
	// Msg is the text of a control line, or the text of an IME event, cut at 63 bytes.
	Msg string
	// Path is the full path of a dropped file, only set on EventDrop.
	Path string
}

func (e Event) c() C.goome_event {
	var ev C.goome_event
	ev._type = C.int(e.Type)
	ev.x, ev.y = C.double(e.X), C.double(e.Y)
	ev.mon, ev.button, ev.dy, ev.key, ev.mods, ev.id = C.int(e.Mon), C.int(e.Button), C.int(e.DY), C.int(e.Key), C.int(e.Mods), C.int(e.ID)
	copyCStr(ev.text[:], e.Text)
	copyCStr(ev.msg[:], e.Msg)
	if (e.Type == EventText || e.Type == EventPreedit) && e.Msg == "" {
		copyCStr(ev.text[:], "")
		copyCStr(ev.msg[:], e.Text)
	}
	return ev
}

func copyCStr(dst []C.char, s string) {
	n := len(s)
	if n > len(dst)-1 {
		n = len(dst) - 1
	}
	for i := 0; i < n; i++ {
		dst[i] = C.char(s[i])
	}
	dst[n] = 0
}

func goStr(src []C.char) string {
	b := make([]byte, 0, len(src))
	for _, c := range src {
		if c == 0 {
			break
		}
		b = append(b, byte(c))
	}
	return string(b)
}

// Goome is an open goome session, a window, overlay or offscreen surface plus its event queue.
type Goome struct {
	g      *C.goome
	strs   []*C.char
	shots  map[int]int
	frames int
	uis    []*UI
}

func fits(sizes ...int) bool {
	for _, v := range sizes {
		if v < 0 || v > math.MaxInt32 {
			return false
		}
	}
	return true
}

func cstr(s string) *C.char {
	if s == "" {
		return nil
	}
	return C.CString(s)
}

func b2i(b bool) C.int {
	if b {
		return 1
	}
	return 0
}

// Open starts a session with the given options. The error says why it failed, like no display.
func Open(o Options) (*Goome, error) {
	gm := &Goome{shots: map[int]int{}}
	keep := func(s string) *C.char {
		p := cstr(s)
		if p != nil {
			gm.strs = append(gm.strs, p)
		}
		return p
	}
	opts := (*C.goome_opts)(C.calloc(1, C.size_t(unsafe.Sizeof(C.goome_opts{}))))
	defer C.free(unsafe.Pointer(opts))
	opts.title = keep(o.Title)
	opts.backend = C.int(o.Backend)
	opts.state_dir = keep(o.StateDir)
	opts.no_capture = b2i(o.NoCapture)
	opts.capture_fps = C.int(o.CaptureFPS)
	opts.hide_self = b2i(o.HideSelf)
	opts.test_frames = keep(o.TestFrames)
	opts.test_script = keep(o.TestScript)
	opts.test_shots = keep(o.TestShots)
	opts.mode = C.int(o.Mode)
	opts.width, opts.height = C.int(o.Width), C.int(o.Height)
	opts.min_width, opts.min_height = C.int(o.MinWidth), C.int(o.MinHeight)
	opts.window_flags = C.int(o.WindowFlags)
	opts.scale = C.double(o.Scale)
	var err C.int
	gm.g = C.goome_open(opts, &err)
	if gm.g == nil {
		gm.freeStrs()
		if err == 0 {
			err = CodeNoDisplay
		}
		return nil, &Error{Code: int(err)}
	}
	return gm, nil
}

// OpenApp opens a normal resizable window, centered, with logical size w by h.
func OpenApp(title string, w, h int) (*Goome, error) {
	return Open(Options{Title: title, Mode: ModeWindow, Width: w, Height: h, WindowFlags: WindowResizable | WindowCenter})
}

// OpenOffscreen opens a session with no window, logical size w by h at the given scale (0 is 1).
func OpenOffscreen(w, h int, scale float64) (*Goome, error) {
	return Open(Options{Mode: ModeOffscreen, Width: w, Height: h, Scale: scale})
}

func (g *Goome) freeStrs() {
	for _, p := range g.strs {
		C.free(unsafe.Pointer(p))
	}
	g.strs = nil
}

// Close frees the UIs and the session. It is safe to call twice.
func (g *Goome) Close() {
	if g == nil || g.g == nil {
		return
	}
	for len(g.uis) > 0 {
		g.uis[len(g.uis)-1].Free()
	}
	C.goome_close(g.g)
	g.g = nil
	g.frames++
	g.shots = map[int]int{}
	g.freeStrs()
}

var errClosed = &Error{Code: CodeGone, Msg: "this goome is closed"}

// Closed reports whether the session is already closed.
func (g *Goome) Closed() bool { return g == nil || g.g == nil }

func (g *Goome) h() *C.goome {
	if g == nil || g.g == nil {
		panic(&Error{Code: CodeGone, Msg: "this goome is closed"})
	}
	return g.g
}

// LastError returns the message of the last failure, or an empty string if closed.
func (g *Goome) LastError() string {
	if g.Closed() {
		return ""
	}
	return C.GoString(C.goome_error(g.g))
}

// BackendName returns the name of the backend that got picked.
func (g *Goome) BackendName() string { return C.GoString(C.goome_backend_name(g.h())) }

// Caps returns the capability bits the backend supports, the Cap* constants ored together.
func (g *Goome) Caps() int { return int(C.goome_caps(g.h())) }

// Can reports whether the backend has every capability bit in capability.
func (g *Goome) Can(capability int) bool {
	return g.Caps()&capability == capability
}

// Monitors lists the monitors, up to 32.
func (g *Goome) Monitors() []Monitor {
	var mons [32]C.goome_monitor
	n := int(C.goome_monitors(g.h(), &mons[0], 32))
	out := make([]Monitor, 0, n)
	for i := 0; i < n; i++ {
		m := mons[i]
		out = append(out, Monitor{int(m.id), goStr(m.name[:]), rectFrom(m.area), float64(m.scale), int(m.pw), int(m.ph), int(m.rotation)})
	}
	return out
}

// Poll waits up to timeoutMs for an event. The bool is false when none came.
func (g *Goome) Poll(timeoutMs int) (Event, bool, error) {
	if g.Closed() {
		return Event{}, false, errClosed
	}
	var ev C.goome_event
	r := C.goome_poll(g.h(), &ev, C.int(timeoutMs))
	if r < 0 {
		return Event{}, false, g.check(r)
	}
	if r == 0 {
		return Event{}, false, nil
	}
	e := Event{
		Type: int(ev._type), X: float64(ev.x), Y: float64(ev.y), Mon: int(ev.mon), Button: int(ev.button),
		DY: int(ev.dy), Key: int(ev.key), Mods: int(ev.mods), Text: goStr(ev.text[:]), ID: int(ev.id), Msg: goStr(ev.msg[:]),
	}
	if e.Type == EventDrop {
		if p := C.goome_drop_path(g.g); p != nil {
			e.Path = C.GoString(p)
		}
	}
	if e.Type == EventText || e.Type == EventPreedit {
		e.Text = C.GoString(C.goome_event_text(g.g))
	}
	return e, true, nil
}

// Events returns an iterator over incoming events. It waits up to timeoutMs for the first one, then drains the queue and stops.
func (g *Goome) Events(timeoutMs int) iter.Seq[Event] {
	return func(yield func(Event) bool) {
		wait := timeoutMs
		for !g.Closed() {
			ev, ok, err := g.Poll(wait)
			if err != nil || !ok {
				return
			}
			if !yield(ev) {
				return
			}
			wait = 0
		}
	}
}

// Pointer returns the global pointer position in logical coordinates and the pressed button bits.
func (g *Goome) Pointer() (x, y float64, buttons int, err error) {
	if g.Closed() {
		return 0, 0, 0, errClosed
	}
	var cx, cy C.double
	var b C.int
	if e := g.check(C.goome_pointer(g.h(), &cx, &cy, &b)); e != nil {
		return 0, 0, 0, e
	}
	return float64(cx), float64(cy), int(b), nil
}

// Keyboard asks the system to send keys to us (on) or stops (off).
func (g *Goome) Keyboard(on bool) error {
	if g.Closed() {
		return errClosed
	}
	return g.check(C.goome_keyboard(g.h(), b2i(on)))
}

// Grab turns input grabbing on or off. In an overlay it takes focus and the pointer, in a window it keeps drags going outside it.
func (g *Goome) Grab(on bool) error { return g.check(C.goome_grab(g.h(), b2i(on))) }

// Grabbed reports whether a grab is on.
func (g *Goome) Grabbed() bool { return C.goome_grabbed(g.h()) > 0 }

// Hotkey registers a global hotkey like "ctrl+alt+k" that comes back as an EventHotkey carrying id.
func (g *Goome) Hotkey(combo string, id int) error {
	if g.Closed() {
		return errClosed
	}
	s := C.CString(combo)
	defer C.free(unsafe.Pointer(s))
	return g.check(C.goome_hotkey(g.h(), s, C.int(id)))
}

// ParseCombo turns a combo like "ctrl+shift+k" into a key and modifier bits.
func ParseCombo(combo string) (key, mods int, err error) {
	s := C.CString(combo)
	defer C.free(unsafe.Pointer(s))
	var k, m C.int
	if r := C.goome_parse_combo(s, &k, &m); r < 0 {
		return 0, 0, newErr(int(r), "cant read the combo %q", combo)
	}
	return int(k), int(m), nil
}

// ControlPath returns the path of the control fifo in the state dir, or an empty string if there is none.
func (g *Goome) ControlPath() string {
	p := C.goome_control_path(g.h())
	if p == nil {
		return ""
	}
	return C.GoString(p)
}

// Send writes a control line to the control fifo in stateDir, which reaches the running program as an EventControl.
func Send(stateDir, msg string) error {
	d, m := C.CString(stateDir), C.CString(msg)
	defer C.free(unsafe.Pointer(d))
	defer C.free(unsafe.Pointer(m))
	if r := C.goome_send(d, m); r < 0 {
		return &Error{Code: int(r)}
	}
	return nil
}

// Pacer keeps a loop at a steady frame rate.
type Pacer struct {
	p *C.goome_pacer
	// FPS is the target frames per second, you can change it between calls.
	FPS int
}

// NewPacer makes a Pacer for a frame loop that runs at fps frames per second. Free it when done.
func NewPacer(fps int) *Pacer {
	p := &Pacer{p: (*C.goome_pacer)(C.calloc(1, C.size_t(unsafe.Sizeof(C.goome_pacer{})))), FPS: fps}
	return p
}

// Wait returns how many milliseconds are left until the next frame is due, pass that to Poll.
func (p *Pacer) Wait() int { return int(C.goome_pace(p.p, C.int(p.FPS))) }

// Free releases the pacer.
func (p *Pacer) Free() {
	if p.p != nil {
		C.free(unsafe.Pointer(p.p))
		p.p = nil
	}
}

// NowMs returns the current time in milliseconds.
func NowMs() int64 { return int64(C.goome_now_ms()) }

// Mode returns the mode this session runs in, one of the Mode* constants.
func (g *Goome) Mode() int { return int(C.goome_mode(g.h())) }

// SetTitle changes the window title.
func (g *Goome) SetTitle(title string) error {
	if g.Closed() {
		return errClosed
	}
	s := C.CString(title)
	defer C.free(unsafe.Pointer(s))
	return g.check(C.goome_window_title(g.h(), s))
}

// SetSize sets the logical client size of the window.
func (g *Goome) SetSize(w, h int) error {
	if g.Closed() {
		return errClosed
	}
	return g.check(C.goome_window_size(g.h(), C.int(w), C.int(h)))
}

// Pos returns the window position. It gives ErrNotSupported where the OS won't say, like Wayland.
func (g *Goome) Pos() (x, y int, err error) {
	if g.Closed() {
		return 0, 0, errClosed
	}
	var cx, cy C.int
	if e := g.check(C.goome_window_pos(g.h(), &cx, &cy, 0)); e != nil {
		return 0, 0, e
	}
	return int(cx), int(cy), nil
}

// SetPos moves the window to x, y.
func (g *Goome) SetPos(x, y int) error {
	if g.Closed() {
		return errClosed
	}
	cx, cy := C.int(x), C.int(y)
	return g.check(C.goome_window_pos(g.h(), &cx, &cy, 1))
}

// SetState changes the window state, one of the State* constants.
func (g *Goome) SetState(state int) error {
	if g.Closed() {
		return errClosed
	}
	return g.check(C.goome_window_state(g.h(), C.int(state)))
}

// State returns the current window state, one of the State* constants.
func (g *Goome) State() int { return int(C.goome_window_get_state(g.h())) }

// Drag starts a native window move (edges 0) or resize (Edge* bits) from the press happening right now.
func (g *Goome) Drag(edges int) error { return g.check(C.goome_window_drag(g.h(), C.int(edges))) }

// Show shows or hides the window.
func (g *Goome) Show(on bool) error { return g.check(C.goome_window_show(g.h(), b2i(on))) }

// SetCursor changes the mouse cursor to one of the Cursor* shapes.
func (g *Goome) SetCursor(shape int) error {
	if g.Closed() {
		return errClosed
	}
	return g.check(C.goome_cursor(g.h(), C.int(shape)))
}

// Cursor returns the last cursor shape that was asked for, handy for offscreen embedders.
func (g *Goome) Cursor() int { return int(C.goome_cursor_shape(g.h())) }

// SetClipboard puts text on the clipboard.
func (g *Goome) SetClipboard(text string) error {
	if g.Closed() {
		return errClosed
	}
	s := C.CString(text)
	defer C.free(unsafe.Pointer(s))
	return g.check(C.goome_clipboard_set(g.h(), s))
}

// Clipboard returns the clipboard as text. It can block for around 250ms on X11 and Wayland.
func (g *Goome) Clipboard() (string, error) {
	if g.Closed() {
		return "", errClosed
	}
	cap := 4096
	for {
		buf := (*C.char)(C.malloc(C.size_t(cap)))
		n := C.goome_clipboard_get(g.h(), buf, C.int(cap))
		if n < 0 {
			C.free(unsafe.Pointer(buf))
			return "", g.check(n)
		}
		if int(n) < cap {
			s := C.GoStringN(buf, n)
			C.free(unsafe.Pointer(buf))
			return s, nil
		}
		C.free(unsafe.Pointer(buf))
		cap = int(n) + 1
	}
}

// DropPath returns the full path of the last EventDrop, good until the next poll.
func (g *Goome) DropPath() string {
	p := C.goome_drop_path(g.h())
	if p == nil {
		return ""
	}
	return C.GoString(p)
}

// EventText returns the full utf-8 of the last EventText or EventPreedit, good until the next poll.
func (g *Goome) EventText() string {
	return C.GoString(C.goome_event_text(g.h()))
}

// PreeditCursor returns the byte offset of the caret in the last preedit text.
func (g *Goome) PreeditCursor() int { return int(C.goome_preedit_cursor(g.h())) }

// TextInput turns IME input on or off. caret is where the candidate window goes, in logical coordinates, and can be nil.
func (g *Goome) TextInput(on bool, caret *Rect) error {
	if caret == nil {
		return g.check(C.goome_text_input(g.h(), b2i(on), nil))
	}
	r := (*C.goome_rect)(C.malloc(C.size_t(unsafe.Sizeof(C.goome_rect{}))))
	defer C.free(unsafe.Pointer(r))
	*r = caret.c()
	return g.check(C.goome_text_input(g.h(), b2i(on), r))
}

// PushIME queues a finished IME text, like a composed or dead key result.
func (g *Goome) PushIME(text string) error {
	return g.Push(Event{Type: EventText, Msg: text})
}

// PushPreedit queues a preedit change with text and the caret byte offset in it. Empty text ends the composing.
func (g *Goome) PushPreedit(text string, cursor int) error {
	return g.Push(Event{Type: EventPreedit, Msg: text, ID: cursor})
}

// Push queues an event as if it came from the system, in any mode.
func (g *Goome) Push(ev Event) error {
	if g.Closed() {
		return errClosed
	}
	e := (*C.goome_event)(C.malloc(C.size_t(unsafe.Sizeof(C.goome_event{}))))
	defer C.free(unsafe.Pointer(e))
	*e = ev.c()
	return g.check(C.goome_push(g.h(), e))
}

// PushClick queues a left click at x, y: a move, a press and a release.
func (g *Goome) PushClick(x, y float64) error {
	if g.Closed() {
		return errClosed
	}
	if err := g.Push(Event{Type: EventMotion, X: x, Y: y}); err != nil {
		return err
	}
	if err := g.Push(Event{Type: EventDown, X: x, Y: y, Button: 1}); err != nil {
		return err
	}
	return g.Push(Event{Type: EventUp, X: x, Y: y, Button: 1})
}

// PushText queues each character of text as a key event.
func (g *Goome) PushText(text string) error {
	if g.Closed() {
		return errClosed
	}
	for _, r := range text {
		if err := g.Push(Event{Type: EventKey, Key: keyFor(r), Text: string(r)}); err != nil {
			return err
		}
	}
	return nil
}

func keyFor(r rune) int {
	if r >= 'A' && r <= 'Z' {
		return int(r - 'A' + 'a')
	}
	if r < 128 {
		return int(r)
	}
	return 0
}

// Resize changes the size (logical) and scale of an offscreen session. Offscreen only.
func (g *Goome) Resize(w, h int, scale float64) error {
	if g.Closed() {
		return errClosed
	}
	return g.check(C.goome_resize(g.h(), C.int(w), C.int(h), C.double(scale)))
}

// Hit reports whether x, y lands in the hit rects of the last End.
func (g *Goome) Hit(x, y float64) bool { return C.goome_hit(g.h(), C.double(x), C.double(y)) > 0 }

// LastFrame returns the last frame End sent, as premultiplied argb, good until the next Begin.
func (g *Goome) LastFrame() (*Image, error) {
	if g.Closed() {
		return nil, errClosed
	}
	img := &Image{g: g, mon: -1, gen: g.frames}
	if err := g.check(C.goome_frame(g.h(), &img.img)); err != nil {
		return nil, err
	}
	return img, nil
}

// FrameDirty returns the real pixel rect that changed since the last call, and false if nothing did. The first call after open or resize is the whole frame.
func (g *Goome) FrameDirty() (Rect, bool) {
	var r C.goome_rect
	if C.goome_frame_dirty(g.h(), &r) <= 0 {
		return Rect{}, false
	}
	return rectFrom(r), true
}

// FrameCopy copies the rect r of the last frame, or all of it if r is nil, into dst in the given pixel format. The stride is in bytes.
func (g *Goome) FrameCopy(dst []byte, strideBytes int, format int, r *Rect) error {
	if g.Closed() {
		return errClosed
	}
	var cr *C.goome_rect
	if r != nil {
		cr = (*C.goome_rect)(C.malloc(C.size_t(unsafe.Sizeof(C.goome_rect{}))))
		defer C.free(unsafe.Pointer(cr))
		*cr = r.c()
	}
	f, err := g.LastFrame()
	if err != nil {
		return err
	}
	w, h := f.W(), f.H()
	if r != nil {
		if !fits(r.X, r.Y, r.W, r.H) || r.X+r.W > w || r.Y+r.H > h {
			return newErr(CodeArg, "that rect isnt inside the %dx%d frame", w, h)
		}
		w, h = r.W, r.H
	}
	if h > 0 && (strideBytes < w*4 || strideBytes > math.MaxInt32 || len(dst) < strideBytes*(h-1)+w*4) {
		return newErr(CodeArg, "dst is too small for %dx%d with stride %d", w, h, strideBytes)
	}
	if len(dst) == 0 {
		return nil
	}
	return g.check(C.goome_frame_copy(g.h(), unsafe.Pointer(&dst[0]), C.int(strideBytes), C.int(format), cr))
}

// Image is a grabbed screen image of 0xFFRRGGBB pixels that goome owns.
type Image struct {
	img C.goome_image
	g   *Goome
	mon int
	gen int
}

// Capture grabs monitor mon, or just the rect r if it isn't nil. The pixels are good until the next grab of that monitor.
func (g *Goome) Capture(mon int, r *Rect) (*Image, error) {
	if g.Closed() {
		return nil, errClosed
	}
	img := (*C.goome_image)(C.calloc(1, C.size_t(unsafe.Sizeof(C.goome_image{}))))
	defer C.free(unsafe.Pointer(img))
	var cr *C.goome_rect
	if r != nil {
		cr = (*C.goome_rect)(C.malloc(C.size_t(unsafe.Sizeof(C.goome_rect{}))))
		defer C.free(unsafe.Pointer(cr))
		*cr = r.c()
	}
	if err := g.check(C.goome_capture(g.h(), C.int(mon), cr, img)); err != nil {
		return nil, err
	}
	g.shots[mon]++
	return &Image{img: *img, g: g, mon: mon, gen: g.shots[mon]}, nil
}

// CaptureWait is Capture that keeps polling and retrying while the backend says try again, for up to about five seconds.
func (g *Goome) CaptureWait(mon int, r *Rect) (*Image, error) {
	var err error
	for try := 0; try < 50; try++ {
		var img *Image
		if img, err = g.Capture(mon, r); err == nil || !errors.Is(err, ErrAgain) {
			return img, err
		}
		g.Poll(100)
	}
	return nil, err
}

// Fresh reports whether the pixels are still valid, they stop being once a newer grab of that monitor replaces them.
func (i *Image) Fresh() bool {
	if i.g.Closed() {
		return false
	}
	if i.mon < 0 {
		return i.gen == i.g.frames
	}
	return i.g.shots[i.mon] == i.gen
}

// W is the width in real pixels.
func (i *Image) W() int { return int(i.img.w) }

// H is the height in real pixels.
func (i *Image) H() int { return int(i.img.h) }

// Stride is the row length in pixels, which can be more than W.
func (i *Image) Stride() int { return int(i.img.stride) }

// Area is the logical desktop rect the image covers.
func (i *Image) Area() Rect { return rectFrom(i.img.area) }

// Mon is the monitor the image came from.
func (i *Image) Mon() int { return int(i.img.mon) }

// HasSelf reports whether our own overlay is in the image, mask it with AlphaAt.
func (i *Image) HasSelf() bool { return i.img.has_self != 0 }

// HasCursor reports whether the mouse cursor got baked into the image.
func (i *Image) HasCursor() bool { return i.img.has_cursor != 0 }

// CursorRect is where the cursor was baked in, in logical coordinates.
func (i *Image) CursorRect() Rect { return rectFrom(i.img.cursor) }

// Pixels returns the pixels without copying them, row by row with Stride, or nil if the image isn't fresh.
func (i *Image) Pixels() []uint32 {
	if !i.Fresh() || i.img.px == nil {
		return nil
	}
	return unsafe.Slice((*uint32)(unsafe.Pointer(i.img.px)), int(i.img.h)*int(i.img.stride))
}

// Pixel returns the pixel at x, y, or 0 if it is out of range or the image is stale.
func (i *Image) Pixel(x, y int) uint32 {
	px := i.Pixels()
	if px == nil || x < 0 || y < 0 || x >= i.W() || y >= i.H() {
		return 0
	}
	return px[y*i.Stride()+x]
}

// Copy returns the pixels as a new w*h slice you own, without the stride padding, or nil if the image isn't fresh.
func (i *Image) Copy() []uint32 {
	px := i.Pixels()
	if px == nil {
		return nil
	}
	w, h, s := i.W(), i.H(), i.Stride()
	out := make([]uint32, w*h)
	for y := 0; y < h; y++ {
		copy(out[y*w:(y+1)*w], px[y*s:y*s+w])
	}
	return out
}

// SavePPM writes the image to path as a PPM file.
func (i *Image) SavePPM(path string) error {
	if !i.Fresh() {
		return newErr(CodeGone, "this image got replaced by a newer grab")
	}
	p := C.CString(path)
	defer C.free(unsafe.Pointer(p))
	if r := C.goome_save_ppm(p, i.img.px, i.img.w, i.img.h, i.img.stride); r < 0 {
		return newErr(int(r), "cant write %s", path)
	}
	return nil
}

// SavePPM writes w by h pixels of px (stride in pixels) to path as a PPM file.
func SavePPM(path string, px []uint32, w, h, stride int) error {
	if w <= 0 || h <= 0 || stride < w || !fits(h, stride) || len(px) < (h-1)*stride+w {
		return newErr(CodeArg, "px is too small for %dx%d with stride %d", w, h, stride)
	}
	p := C.CString(path)
	defer C.free(unsafe.Pointer(p))
	buf := (*C.uint32_t)(C.malloc(C.size_t(len(px) * 4)))
	defer C.free(unsafe.Pointer(buf))
	copy(unsafe.Slice((*uint32)(unsafe.Pointer(buf)), len(px)), px)
	if r := C.goome_save_ppm(p, buf, C.int(w), C.int(h), C.int(stride)); r < 0 {
		return newErr(int(r), "cant write %s", path)
	}
	return nil
}

// AlphaAt returns the alpha of our own overlay at x, y, so you can mask it out of a grab that has HasSelf set.
func (g *Goome) AlphaAt(x, y float64) int {
	return int(C.goome_alpha_at(g.h(), C.double(x), C.double(y)))
}

// Begin starts a frame on monitor mon and returns the canvas to draw on, or nil if there is nothing to draw on. Finish it with End.
func (g *Goome) Begin(mon int) *Canvas {
	c := C.goome_begin(g.h(), C.int(mon))
	g.frames++
	if c == nil {
		return nil
	}
	return &Canvas{c: c, g: g, gen: g.frames}
}

// End sends the frame from Begin, with the hit rects of the canvas plus any extra hits you pass.
func (g *Goome) End(c *Canvas, hits ...Rect) error {
	if g.Closed() {
		return errClosed
	}
	if c == nil {
		return nil
	}
	cc := c.ptr()
	all := append(c.hits, hits...)
	c.hits = nil
	c.c = nil
	var buf *C.goome_rect
	if len(all) > 0 {
		buf = (*C.goome_rect)(C.malloc(C.size_t(len(all)) * C.size_t(unsafe.Sizeof(C.goome_rect{}))))
		defer C.free(unsafe.Pointer(buf))
		s := unsafe.Slice(buf, len(all))
		for i, r := range all {
			s[i] = r.c()
		}
	}
	return g.check(C.goome_end(g.h(), cc, buf, C.int(len(all))))
}

// Frame runs draw on a fresh canvas for monitor mon and sends it, a Begin and End in one call.
func (g *Goome) Frame(mon int, draw func(c *Canvas)) error {
	c := g.Begin(mon)
	if c == nil {
		return nil
	}
	draw(c)
	return g.End(c)
}
