package goome

/*
#include <stdlib.h>
#include "native/goome.h"
*/
import "C"

import (
	"errors"
	"iter"
	"unsafe"
)

type Options struct {
	Title       string
	Backend     int
	StateDir    string
	NoCapture   bool
	CaptureFPS  int
	HideSelf    bool
	TestFrames  string
	TestScript  string
	TestShots   string
	Mode        int
	Width       int
	Height      int
	MinWidth    int
	MinHeight   int
	WindowFlags int
	Scale       float64
}

type Rect struct {
	X, Y, W, H int
}

func (r Rect) c() C.goome_rect {
	return C.goome_rect{x: C.int(r.X), y: C.int(r.Y), w: C.int(r.W), h: C.int(r.H)}
}

func rectFrom(r C.goome_rect) Rect {
	return Rect{int(r.x), int(r.y), int(r.w), int(r.h)}
}

type Monitor struct {
	ID       int
	Name     string
	Area     Rect
	Scale    float64
	PW, PH   int
	Rotation int
}

type Event struct {
	Type   int
	X, Y   float64
	Mon    int
	Button int
	DY     int
	Key    int
	Mods   int
	Text   string
	ID     int
	Msg    string
	Path   string
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

type Goome struct {
	g      *C.goome
	strs   []*C.char
	shots  map[int]int
	frames int
	uis    []*UI
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

func OpenApp(title string, w, h int) (*Goome, error) {
	return Open(Options{Title: title, Mode: ModeWindow, Width: w, Height: h, WindowFlags: WindowResizable | WindowCenter})
}

func OpenOffscreen(w, h int, scale float64) (*Goome, error) {
	return Open(Options{Mode: ModeOffscreen, Width: w, Height: h, Scale: scale})
}

func (g *Goome) freeStrs() {
	for _, p := range g.strs {
		C.free(unsafe.Pointer(p))
	}
	g.strs = nil
}

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

func (g *Goome) Closed() bool { return g == nil || g.g == nil }

func (g *Goome) h() *C.goome {
	if g == nil || g.g == nil {
		panic(&Error{Code: CodeGone, Msg: "this goome is closed"})
	}
	return g.g
}

func (g *Goome) LastError() string {
	if g.Closed() {
		return ""
	}
	return C.GoString(C.goome_error(g.g))
}

func (g *Goome) BackendName() string { return C.GoString(C.goome_backend_name(g.h())) }
func (g *Goome) Caps() int           { return int(C.goome_caps(g.h())) }
func (g *Goome) Can(capability int) bool {
	return g.Caps()&capability == capability
}

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

func (g *Goome) Keyboard(on bool) error {
	if g.Closed() {
		return errClosed
	}
	return g.check(C.goome_keyboard(g.h(), b2i(on)))
}
func (g *Goome) Grab(on bool) error { return g.check(C.goome_grab(g.h(), b2i(on))) }
func (g *Goome) Grabbed() bool      { return C.goome_grabbed(g.h()) > 0 }

func (g *Goome) Hotkey(combo string, id int) error {
	if g.Closed() {
		return errClosed
	}
	s := C.CString(combo)
	defer C.free(unsafe.Pointer(s))
	return g.check(C.goome_hotkey(g.h(), s, C.int(id)))
}

func ParseCombo(combo string) (key, mods int, err error) {
	s := C.CString(combo)
	defer C.free(unsafe.Pointer(s))
	var k, m C.int
	if r := C.goome_parse_combo(s, &k, &m); r < 0 {
		return 0, 0, newErr(int(r), "cant read the combo %q", combo)
	}
	return int(k), int(m), nil
}

func (g *Goome) ControlPath() string {
	p := C.goome_control_path(g.h())
	if p == nil {
		return ""
	}
	return C.GoString(p)
}

func Send(stateDir, msg string) error {
	d, m := C.CString(stateDir), C.CString(msg)
	defer C.free(unsafe.Pointer(d))
	defer C.free(unsafe.Pointer(m))
	if r := C.goome_send(d, m); r < 0 {
		return &Error{Code: int(r)}
	}
	return nil
}

type Pacer struct {
	p   *C.goome_pacer
	FPS int
}

func NewPacer(fps int) *Pacer {
	p := &Pacer{p: (*C.goome_pacer)(C.calloc(1, C.size_t(unsafe.Sizeof(C.goome_pacer{})))), FPS: fps}
	return p
}

func (p *Pacer) Wait() int { return int(C.goome_pace(p.p, C.int(p.FPS))) }

func (p *Pacer) Free() {
	if p.p != nil {
		C.free(unsafe.Pointer(p.p))
		p.p = nil
	}
}

func NowMs() int64 { return int64(C.goome_now_ms()) }

func (g *Goome) Mode() int { return int(C.goome_mode(g.h())) }

func (g *Goome) SetTitle(title string) error {
	if g.Closed() {
		return errClosed
	}
	s := C.CString(title)
	defer C.free(unsafe.Pointer(s))
	return g.check(C.goome_window_title(g.h(), s))
}

func (g *Goome) SetSize(w, h int) error {
	if g.Closed() {
		return errClosed
	}
	return g.check(C.goome_window_size(g.h(), C.int(w), C.int(h)))
}

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

func (g *Goome) SetPos(x, y int) error {
	if g.Closed() {
		return errClosed
	}
	cx, cy := C.int(x), C.int(y)
	return g.check(C.goome_window_pos(g.h(), &cx, &cy, 1))
}

func (g *Goome) SetState(state int) error {
	if g.Closed() {
		return errClosed
	}
	return g.check(C.goome_window_state(g.h(), C.int(state)))
}
func (g *Goome) State() int           { return int(C.goome_window_get_state(g.h())) }
func (g *Goome) Drag(edges int) error { return g.check(C.goome_window_drag(g.h(), C.int(edges))) }
func (g *Goome) Show(on bool) error   { return g.check(C.goome_window_show(g.h(), b2i(on))) }
func (g *Goome) SetCursor(shape int) error {
	if g.Closed() {
		return errClosed
	}
	return g.check(C.goome_cursor(g.h(), C.int(shape)))
}
func (g *Goome) Cursor() int { return int(C.goome_cursor_shape(g.h())) }

func (g *Goome) SetClipboard(text string) error {
	if g.Closed() {
		return errClosed
	}
	s := C.CString(text)
	defer C.free(unsafe.Pointer(s))
	return g.check(C.goome_clipboard_set(g.h(), s))
}

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

func (g *Goome) DropPath() string {
	p := C.goome_drop_path(g.h())
	if p == nil {
		return ""
	}
	return C.GoString(p)
}

func (g *Goome) EventText() string {
	return C.GoString(C.goome_event_text(g.h()))
}

func (g *Goome) PreeditCursor() int { return int(C.goome_preedit_cursor(g.h())) }

func (g *Goome) TextInput(on bool, caret *Rect) error {
	if caret == nil {
		return g.check(C.goome_text_input(g.h(), b2i(on), nil))
	}
	r := (*C.goome_rect)(C.malloc(C.size_t(unsafe.Sizeof(C.goome_rect{}))))
	defer C.free(unsafe.Pointer(r))
	*r = caret.c()
	return g.check(C.goome_text_input(g.h(), b2i(on), r))
}

func (g *Goome) PushIME(text string) error {
	return g.Push(Event{Type: EventText, Msg: text})
}

func (g *Goome) PushPreedit(text string, cursor int) error {
	return g.Push(Event{Type: EventPreedit, Msg: text, ID: cursor})
}

func (g *Goome) Push(ev Event) error {
	if g.Closed() {
		return errClosed
	}
	e := (*C.goome_event)(C.malloc(C.size_t(unsafe.Sizeof(C.goome_event{}))))
	defer C.free(unsafe.Pointer(e))
	*e = ev.c()
	return g.check(C.goome_push(g.h(), e))
}

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

func (g *Goome) Resize(w, h int, scale float64) error {
	if g.Closed() {
		return errClosed
	}
	return g.check(C.goome_resize(g.h(), C.int(w), C.int(h), C.double(scale)))
}

func (g *Goome) Hit(x, y float64) bool { return C.goome_hit(g.h(), C.double(x), C.double(y)) > 0 }

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

func (g *Goome) FrameDirty() (Rect, bool) {
	var r C.goome_rect
	if C.goome_frame_dirty(g.h(), &r) <= 0 {
		return Rect{}, false
	}
	return rectFrom(r), true
}

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
	w, h := 0, 0
	if r != nil {
		w, h = r.W, r.H
	} else if f, err := g.LastFrame(); err == nil {
		w, h = f.W(), f.H()
	} else {
		return err
	}
	if h > 0 && (strideBytes < w*4 || len(dst) < strideBytes*(h-1)+w*4) {
		return newErr(CodeArg, "dst is too small for %dx%d with stride %d", w, h, strideBytes)
	}
	if len(dst) == 0 {
		return nil
	}
	return g.check(C.goome_frame_copy(g.h(), unsafe.Pointer(&dst[0]), C.int(strideBytes), C.int(format), cr))
}

type Image struct {
	img C.goome_image
	g   *Goome
	mon int
	gen int
}

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

func (i *Image) Fresh() bool {
	if i.g.Closed() {
		return false
	}
	if i.mon < 0 {
		return i.gen == i.g.frames
	}
	return i.g.shots[i.mon] == i.gen
}

func (i *Image) W() int           { return int(i.img.w) }
func (i *Image) H() int           { return int(i.img.h) }
func (i *Image) Stride() int      { return int(i.img.stride) }
func (i *Image) Area() Rect       { return rectFrom(i.img.area) }
func (i *Image) Mon() int         { return int(i.img.mon) }
func (i *Image) HasSelf() bool    { return i.img.has_self != 0 }
func (i *Image) HasCursor() bool  { return i.img.has_cursor != 0 }
func (i *Image) CursorRect() Rect { return rectFrom(i.img.cursor) }

func (i *Image) Pixels() []uint32 {
	if !i.Fresh() || i.img.px == nil {
		return nil
	}
	return unsafe.Slice((*uint32)(unsafe.Pointer(i.img.px)), int(i.img.h)*int(i.img.stride))
}

func (i *Image) Pixel(x, y int) uint32 {
	px := i.Pixels()
	if px == nil || x < 0 || y < 0 || x >= i.W() || y >= i.H() {
		return 0
	}
	return px[y*i.Stride()+x]
}

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

func SavePPM(path string, px []uint32, w, h, stride int) error {
	if w <= 0 || h <= 0 || stride < w || len(px) < (h-1)*stride+w {
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

func (g *Goome) AlphaAt(x, y float64) int {
	return int(C.goome_alpha_at(g.h(), C.double(x), C.double(y)))
}

func (g *Goome) Begin(mon int) *Canvas {
	c := C.goome_begin(g.h(), C.int(mon))
	g.frames++
	if c == nil {
		return nil
	}
	return &Canvas{c: c, g: g, gen: g.frames}
}

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

func (g *Goome) Frame(mon int, draw func(c *Canvas)) error {
	c := g.Begin(mon)
	if c == nil {
		return nil
	}
	draw(c)
	return g.End(c)
}
