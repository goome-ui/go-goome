package goome

/*
#include <stdlib.h>
#include "native/goome.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type Radii struct {
	TL float64 `c:"tl"`
	TR float64 `c:"tr"`
	BR float64 `c:"br"`
	BL float64 `c:"bl"`
}

func Radius(r float64) Radii { return Radii{r, r, r, r} }

func (r Radii) c() C.goome_radii { return *(*C.goome_radii)(unsafe.Pointer(&r)) }

type Canvas struct {
	c     *C.goome_canvas
	g     *Goome
	gen   int
	own   bool
	hits  []Rect
	layer bool
}

func NewCanvas(w, h, stride int, scale float64) (*Canvas, error) {
	if stride <= 0 {
		stride = w
	}
	if w < 0 || h < 0 || stride < w {
		return nil, newErr(CodeArg, "bad canvas size %dx%d stride %d", w, h, stride)
	}
	if scale <= 0 {
		scale = 1
	}
	c := (*C.goome_canvas)(C.calloc(1, C.size_t(unsafe.Sizeof(C.goome_canvas{}))))
	n := h * stride
	if n < 1 {
		n = 1
	}
	c.px = (*C.uint32_t)(C.calloc(C.size_t(n), 4))
	if c.px == nil {
		C.free(unsafe.Pointer(c))
		return nil, ErrNoMem
	}
	c.w, c.h, c.stride, c.scale = C.int(w), C.int(h), C.int(stride), C.double(scale)
	cv := &Canvas{c: c, own: true}
	runtime.SetFinalizer(cv, (*Canvas).Free)
	return cv, nil
}

func (c *Canvas) Free() {
	if c == nil || !c.own || c.c == nil {
		return
	}
	C.free(unsafe.Pointer(c.c.px))
	C.free(unsafe.Pointer(c.c))
	c.c = nil
	runtime.SetFinalizer(c, nil)
}

func (c *Canvas) Valid() bool {
	if c == nil || c.c == nil {
		return false
	}
	if c.own {
		return true
	}
	return !c.g.Closed() && c.g.frames == c.gen
}

func (c *Canvas) ptr() *C.goome_canvas {
	if !c.Valid() {
		panic(&Error{Code: CodeGone, Msg: "this canvas is done, get a new one from Begin"})
	}
	return c.c
}

func (c *Canvas) W() int         { return int(c.ptr().w) }
func (c *Canvas) H() int         { return int(c.ptr().h) }
func (c *Canvas) Stride() int    { return int(c.ptr().stride) }
func (c *Canvas) Scale() float64 { return float64(c.ptr().scale) }
func (c *Canvas) Origin() (x, y float64) {
	p := c.ptr()
	return float64(p.ox), float64(p.oy)
}
func (c *Canvas) SetOrigin(x, y float64) {
	p := c.ptr()
	p.ox, p.oy = C.double(x), C.double(y)
}
func (c *Canvas) Mon() int { return int(c.ptr().mon) }

func (c *Canvas) Pixels() []uint32 {
	p := c.ptr()
	return unsafe.Slice((*uint32)(unsafe.Pointer(p.px)), int(p.h)*int(p.stride))
}

func (c *Canvas) Drawn() (x0, y0, x1, y1 int) {
	p := c.ptr()
	return int(p.dx0), int(p.dy0), int(p.dx1), int(p.dy1)
}

func (c *Canvas) Cleared() (x0, y0, x1, y1 int) {
	p := c.ptr()
	return int(p.cx0), int(p.cy0), int(p.cx1), int(p.cy1)
}

func (c *Canvas) ClipPx() (x0, y0, x1, y1 int) {
	p := c.ptr()
	return int(p.kx0), int(p.ky0), int(p.kx1), int(p.ky1)
}

func (c *Canvas) SetClipPx(x0, y0, x1, y1 int) {
	p := c.ptr()
	p.kx0, p.ky0, p.kx1, p.ky1 = C.int(x0), C.int(y0), C.int(x1), C.int(y1)
}

func (c *Canvas) ClipDepth() int { return int(c.ptr().nclip) }
func (c *Canvas) SetClipDepth(n int) {
	c.ptr().nclip = C.int(n)
}

func (c *Canvas) Hit(x, y, w, h float64) {
	c.hits = append(c.hits, Rect{int(x), int(y), int(w + 0.999999), int(h + 0.999999)})
}

func (c *Canvas) Damage(r Rect) {
	cr := (*C.goome_rect)(C.malloc(C.size_t(unsafe.Sizeof(C.goome_rect{}))))
	defer C.free(unsafe.Pointer(cr))
	*cr = r.c()
	C.goome_damage(c.ptr(), cr)
}

func (c *Canvas) Clear() { C.goome_clear(c.ptr()) }

func (c *Canvas) Rect(x, y, w, h float64, argb uint32) {
	C.goome_draw_rect(c.ptr(), C.double(x), C.double(y), C.double(w), C.double(h), C.uint32_t(argb))
}

func (c *Canvas) RoundRect(x, y, w, h, radius float64, argb uint32) {
	C.goome_draw_round_rect(c.ptr(), C.double(x), C.double(y), C.double(w), C.double(h), C.double(radius), C.uint32_t(argb))
}

func (c *Canvas) Frame(x, y, w, h, width float64, argb uint32) {
	C.goome_draw_frame(c.ptr(), C.double(x), C.double(y), C.double(w), C.double(h), C.double(width), C.uint32_t(argb))
}

func (c *Canvas) Line(x0, y0, x1, y1, width float64, argb uint32) {
	C.goome_draw_line(c.ptr(), C.double(x0), C.double(y0), C.double(x1), C.double(y1), C.double(width), C.uint32_t(argb))
}

func (c *Canvas) Triangle(x0, y0, x1, y1, x2, y2 float64, argb uint32) {
	C.goome_draw_triangle(c.ptr(), C.double(x0), C.double(y0), C.double(x1), C.double(y1), C.double(x2), C.double(y2), C.uint32_t(argb))
}

func cPixels(px []uint32) (*C.uint32_t, func()) {
	if len(px) == 0 {
		return nil, func() {}
	}
	buf := (*C.uint32_t)(C.malloc(C.size_t(len(px) * 4)))
	copy(unsafe.Slice((*uint32)(unsafe.Pointer(buf)), len(px)), px)
	return buf, func() { C.free(unsafe.Pointer(buf)) }
}

func (c *Canvas) Image(x, y, w, h float64, px []uint32, iw, ih int) error {
	if iw <= 0 || ih <= 0 || len(px) < iw*ih {
		return newErr(CodeArg, "px is too small for %dx%d", iw, ih)
	}
	buf, done := cPixels(px[:iw*ih])
	defer done()
	C.goome_draw_image(c.ptr(), C.double(x), C.double(y), C.double(w), C.double(h), buf, C.int(iw), C.int(ih))
	return nil
}

func (c *Canvas) Erase(x, y, w, h float64) {
	C.goome_erase(c.ptr(), C.double(x), C.double(y), C.double(w), C.double(h))
}

func (c *Canvas) Text(x, y float64, s string, big bool, argb uint32) float64 {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return float64(C.goome_draw_text(c.ptr(), C.double(x), C.double(y), cs, b2i(big), C.uint32_t(argb)))
}

func (c *Canvas) TextWidth(s string, big bool) float64 {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return float64(C.goome_text_width(c.ptr(), cs, b2i(big)))
}

func SetFont(font int) { C.goome_set_font(C.int(font)) }
func FontCount() int   { return int(C.goome_font_count()) }

func (c *Canvas) Blit(src *Canvas, alpha, dx, dy float64) {
	C.goome_canvas_blit(c.ptr(), src.ptr(), C.double(alpha), C.double(dx), C.double(dy))
}

func (c *Canvas) ClipPush(x, y, w, h float64) {
	C.goome_clip_push(c.ptr(), C.double(x), C.double(y), C.double(w), C.double(h))
}

func (c *Canvas) ClipPop() { C.goome_clip_pop(c.ptr()) }

func (c *Canvas) Clip(x, y, w, h float64, draw func()) {
	c.ClipPush(x, y, w, h)
	defer c.ClipPop()
	draw()
}

func (c *Canvas) RRect(x, y, w, h float64, r Radii, argb uint32) {
	C.goome_draw_rrect(c.ptr(), C.double(x), C.double(y), C.double(w), C.double(h), r.c(), C.uint32_t(argb))
}

func (c *Canvas) RRectLine(x, y, w, h float64, r Radii, width float64, argb uint32) {
	C.goome_draw_rrect_line(c.ptr(), C.double(x), C.double(y), C.double(w), C.double(h), r.c(), C.double(width), C.uint32_t(argb))
}

func (c *Canvas) RRectGrad(x, y, w, h float64, r Radii, from, to uint32, sideways bool) {
	C.goome_draw_rrect_grad(c.ptr(), C.double(x), C.double(y), C.double(w), C.double(h), r.c(), C.uint32_t(from), C.uint32_t(to), b2i(sideways))
}

func (c *Canvas) Shadow(x, y, w, h, radius, blur float64, argb uint32) {
	C.goome_draw_shadow(c.ptr(), C.double(x), C.double(y), C.double(w), C.double(h), C.double(radius), C.double(blur), C.uint32_t(argb))
}

func (c *Canvas) Circle(cx, cy, r float64, argb uint32) {
	C.goome_draw_circle(c.ptr(), C.double(cx), C.double(cy), C.double(r), C.uint32_t(argb))
}

func (c *Canvas) Ring(cx, cy, r, width float64, argb uint32) {
	C.goome_draw_ring(c.ptr(), C.double(cx), C.double(cy), C.double(r), C.double(width), C.uint32_t(argb))
}

func (c *Canvas) Arc(cx, cy, r, width, a0, a1 float64, argb uint32) {
	C.goome_draw_arc(c.ptr(), C.double(cx), C.double(cy), C.double(r), C.double(width), C.double(a0), C.double(a1), C.uint32_t(argb))
}

func (c *Canvas) Poly(xy []float64, argb uint32) {
	n := len(xy) / 2
	if n == 0 {
		return
	}
	buf := (*C.double)(C.malloc(C.size_t(n * 2 * 8)))
	defer C.free(unsafe.Pointer(buf))
	copy(unsafe.Slice((*float64)(unsafe.Pointer(buf)), n*2), xy[:n*2])
	C.goome_draw_poly(c.ptr(), buf, C.int(n), C.uint32_t(argb))
}

func (c *Canvas) Bitmap(b *Bitmap, x, y, w, h float64, r Radii, tint uint32) {
	C.goome_draw_bitmap(c.ptr(), b.ptr(), C.double(x), C.double(y), C.double(w), C.double(h), r.c(), C.uint32_t(tint))
}

func (c *Canvas) BitmapPart(b *Bitmap, sx, sy, sw, sh int, x, y, w, h float64, tint uint32) {
	C.goome_draw_bitmap_part(c.ptr(), b.ptr(), C.int(sx), C.int(sy), C.int(sw), C.int(sh), C.double(x), C.double(y), C.double(w), C.double(h), C.uint32_t(tint))
}

func (c *Canvas) BitmapNine(b *Bitmap, x, y, w, h, l, t, r, bm float64) {
	C.goome_draw_bitmap_nine(c.ptr(), b.ptr(), C.double(x), C.double(y), C.double(w), C.double(h), C.double(l), C.double(t), C.double(r), C.double(bm))
}

func (c *Canvas) FontText(f *Font, size, x, y float64, s string, argb uint32) float64 {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	return float64(C.goome_text(c.ptr(), f.raw(), C.double(size), C.double(x), C.double(y), cs, C.uint32_t(argb)))
}

type Bitmap struct {
	b     *C.goome_bitmap
	wraps bool
}

func newBitmap() *Bitmap {
	b := &Bitmap{b: (*C.goome_bitmap)(C.calloc(1, C.size_t(unsafe.Sizeof(C.goome_bitmap{}))))}
	runtime.SetFinalizer(b, (*Bitmap).Free)
	return b
}

func bitmapErr(r C.int, what string) error {
	why := C.GoString(C.goome_bitmap_error())
	if why == "" {
		why = Strerror(int(r))
	}
	return &Error{Code: int(r), Msg: what + ": " + why}
}

func LoadBitmap(path string) (*Bitmap, error) {
	p := C.CString(path)
	defer C.free(unsafe.Pointer(p))
	b := newBitmap()
	if r := C.goome_bitmap_load(p, b.b); r < 0 {
		b.Free()
		return nil, bitmapErr(r, path)
	}
	return b, nil
}

func BitmapFromBytes(data []byte) (*Bitmap, error) {
	if len(data) == 0 {
		return nil, newErr(CodeArg, "no image data")
	}
	buf := C.CBytes(data)
	defer C.free(buf)
	b := newBitmap()
	if r := C.goome_bitmap_load_mem(buf, C.size_t(len(data)), b.b); r < 0 {
		b.Free()
		return nil, bitmapErr(r, "image")
	}
	return b, nil
}

func BitmapFromARGB(px []uint32, w, h int) (*Bitmap, error) {
	if w <= 0 || h <= 0 || len(px) < w*h {
		return nil, newErr(CodeArg, "px is too small for %dx%d", w, h)
	}
	buf, done := cPixels(px[:w*h])
	defer done()
	b := newBitmap()
	if r := C.goome_bitmap_from_argb(buf, C.int(w), C.int(h), b.b); r < 0 {
		b.Free()
		return nil, &Error{Code: int(r)}
	}
	return b, nil
}

func NewBitmap(w, h, stride int) (*Bitmap, error) {
	if stride <= 0 {
		stride = w
	}
	if w <= 0 || h <= 0 || stride < w {
		return nil, newErr(CodeArg, "bad bitmap size %dx%d stride %d", w, h, stride)
	}
	b := newBitmap()
	b.b.px = (*C.uint32_t)(C.calloc(C.size_t(h*stride), 4))
	if b.b.px == nil {
		b.Free()
		return nil, ErrNoMem
	}
	b.b.w, b.b.h, b.b.stride = C.int(w), C.int(h), C.int(stride)
	b.wraps = true
	return b, nil
}

func (b *Bitmap) Free() {
	if b == nil || b.b == nil {
		return
	}
	for _, u := range liveUIs {
		u.forgetBitmap(b)
	}
	unpin(unsafe.Pointer(b.b))
	if b.wraps {
		C.free(unsafe.Pointer(b.b.px))
	} else {
		C.goome_bitmap_free(b.b)
	}
	C.free(unsafe.Pointer(b.b))
	b.b = nil
	runtime.SetFinalizer(b, nil)
}

func (b *Bitmap) ptr() *C.goome_bitmap {
	if b == nil || b.b == nil {
		panic(&Error{Code: CodeGone, Msg: "this bitmap is freed"})
	}
	return b.b
}

func (b *Bitmap) Freed() bool { return b == nil || b.b == nil }
func (b *Bitmap) W() int      { return int(b.ptr().w) }
func (b *Bitmap) H() int      { return int(b.ptr().h) }
func (b *Bitmap) Stride() int { return int(b.ptr().stride) }
func (b *Bitmap) Owned() bool { return b.ptr().owned != 0 }
func (b *Bitmap) Pixels() []uint32 {
	p := b.ptr()
	return unsafe.Slice((*uint32)(unsafe.Pointer(p.px)), int(p.h)*int(p.stride))
}

type Font struct {
	f     *C.goome_font
	fb    []*Font
	users int
	freed bool
}

func wrapFont(f *C.goome_font, what string) (*Font, error) {
	if f == nil {
		return nil, newErr(CodeArg, "cant load the font %s", what)
	}
	ft := &Font{f: f}
	runtime.SetFinalizer(ft, (*Font).Free)
	return ft, nil
}

func LoadFont(path string) (*Font, error) {
	p := C.CString(path)
	defer C.free(unsafe.Pointer(p))
	return wrapFont(C.goome_font_load(p), path)
}

func LoadFontIndex(path string, index int) (*Font, error) {
	p := C.CString(path)
	defer C.free(unsafe.Pointer(p))
	return wrapFont(C.goome_font_load_index(p, C.int(index)), path)
}

func FontFromBytes(data []byte) (*Font, error) {
	if len(data) == 0 {
		return nil, newErr(CodeArg, "no font data")
	}
	buf := C.CBytes(data)
	defer C.free(buf)
	return wrapFont(C.goome_font_load_mem(buf, C.size_t(len(data)), 1), "from memory")
}

func SystemFont(name string) (*Font, error) {
	var p *C.char
	if name != "" {
		p = C.CString(name)
		defer C.free(unsafe.Pointer(p))
	}
	if name == "" {
		name = "ui font"
	}
	return wrapFont(C.goome_font_system(p), name)
}

func (f *Font) Free() {
	if f == nil || f.f == nil || f.freed {
		return
	}
	for _, u := range liveUIs {
		u.forgetFont(f)
	}
	f.freed = true
	runtime.SetFinalizer(f, nil)
	f.release()
}

func (f *Font) release() {
	if f.users > 0 || f.f == nil {
		return
	}
	unpin(unsafe.Pointer(f.f))
	C.goome_font_free(f.f)
	f.f = nil
	for _, b := range f.fb {
		b.users--
		if b.freed {
			b.release()
		}
	}
	f.fb = nil
}

func (f *Font) raw() *C.goome_font {
	if f == nil {
		return nil
	}
	if f.f == nil || f.freed {
		panic(&Error{Code: CodeGone, Msg: "this font is freed"})
	}
	return f.f
}

func (f *Font) Fallback(fallback *Font) {
	p, b := f.raw(), fallback.raw()
	if b == nil || b == p {
		return
	}
	for _, x := range f.fb {
		if x == fallback {
			return
		}
	}
	C.goome_font_fallback(p, b)
	f.fb = append(f.fb, fallback)
	fallback.users++
}

func (f *Font) Fallbacks() []*Font { return append([]*Font(nil), f.fb...) }

func SetTextDir(dir int) { C.goome_text_dir(C.int(dir)) }

func (f *Font) Measure(text string, size float64) (w, h float64) {
	cs := C.CString(text)
	defer C.free(unsafe.Pointer(cs))
	var ch C.double
	cw := C.goome_text_measure(f.raw(), C.double(size), cs, &ch)
	return float64(cw), float64(ch)
}

func (f *Font) Width(text string, size float64) float64 {
	w, _ := f.Measure(text, size)
	return w
}

func (f *Font) LineHeight(size float64) float64 {
	return float64(C.goome_text_line_height(f.raw(), C.double(size)))
}
