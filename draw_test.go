package goome

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var star = []float64{80, 5, 98, 85, 40, 35, 120, 35, 62, 85}

var shapeHashes = [][]uint32{
	{0xA7B3B4C5, 0x0FB29A82, 0x1A98CF16, 0x1E7B729F, 0x21FD3B05, 0x241CE7FD, 0x6910B0D5, 0x7180B28E, 0x6F33AC30, 0xB48ADFA3, 0x796B9105, 0x2D93DC92, 0x1F908FAD, 0x23A6CF0C},
	{0x73429A18, 0xD9BD97AB, 0xD031F56C, 0xD2D9FA9A, 0xD639A1A7, 0xA77C3E9E, 0x7BD3FEAD, 0xAC76561E, 0x5BD0648D, 0xD0AC70E1, 0x0BB5DDA0, 0xE7C304C7, 0xCCEF8860, 0xEFC1F212},
	{0x9175788F, 0x6E99EA87, 0x935CC8E1, 0x0B73F68D, 0x23891BE5, 0xA28476C1, 0xEA64F655, 0x87DC9430, 0x909D11B5, 0x0550BB7A, 0x6455EE86, 0xD8C08E2B, 0xFD0DD1FD, 0x28C27977},
}

var imageCases = []struct {
	name string
	w, h int
	hash uint32
}{
	{"png-gray-1-bit.png", 9, 3, 0x01975AD8},
	{"png-palette-4-bit-trns-interlaced.png", 6, 6, 0xAF77E7E2},
	{"png-rgb-16-bit.png", 5, 4, 0xDD60C23E},
	{"png-gray-alpha-stored-interlaced.png", 5, 5, 0x52D68985},
	{"png-rgba-fixed-huffman-split.png", 8, 6, 0x66DB9DCD},
	{"png-gray-key.png", 7, 4, 0x843ED605},
	{"png-rgba-16-bit-interlaced.png", 9, 9, 0xDC2F91D5},
	{"jpeg-420.jpg", 16, 16, 0x32E2C9CE},
	{"jpeg-422-restarts.jpg", 13, 9, 0x0E230BAF},
	{"jpeg-444.jpg", 10, 10, 0x34FC5275},
	{"jpeg-progressive-420-restarts.jpg", 14, 11, 0x5BCEE932},
	{"jpeg-gray.jpg", 9, 7, 0xF010B3C8},
	{"jpeg-exif-rotated.jpg", 6, 10, 0x65909E25},
	{"bmp-4-bit.bmp", 5, 3, 0x9256D097},
	{"bmp-32-bit-alpha-top-down.bmp", 3, 2, 0x6BDB3A29},
	{"bmp-24-bit.bmp", 5, 3, 0x09D94C56},
	{"ppm.ppm", 3, 2, 0xC6FDF725},
	{"pgm-16-bit.pgm", 3, 1, 0x35029201},
	{"qoi.qoi", 6, 4, 0x23AF45AD},
}

func shape(c *Canvas, k int, b *Bitmap) {
	switch k {
	case 0:
		c.RRect(10.3, 8.6, 120, 70.2, Radii{4, 20, 0, 33}, 0xE02060C0)
	case 1:
		c.RRectLine(10.3, 8.6, 120, 70.2, Radius(18), 2.5, 0xFF20A060)
	case 2:
		c.RRectGrad(10, 10, 130, 70, Radius(12), 0xFF102040, 0x8080C0FF, false)
	case 3:
		c.RRectGrad(10, 10, 130, 70, Radius(35), 0xFFFF4080, 0x004080FF, true)
	case 4:
		c.Shadow(30, 20, 90, 50, 10, 16, 0x80000000)
	case 5:
		c.Circle(50.5, 50.2, 30.7, 0xFFFFC020)
		c.Circle(120, 30, 1.3, 0xFF000000)
	case 6:
		c.Ring(80, 50, 30, 5, 0xC0208080)
	case 7:
		c.Arc(50, 50, 30, 6, 0.3, 4.5, 0xFF3060E0)
		c.Arc(120, 50, 20, 30, 2, 2.5, 0xFF905020)
	case 8:
		c.Poly(star, 0xFF2040A0)
	case 9:
		c.Bitmap(b, 10, 10, 100.5, 75.3, Radius(20), 0xFFFFFFFF)
	case 10:
		c.Bitmap(b, 10, 10, 24, 16, Radius(0), 0x80FF8080)
	case 11:
		c.Bitmap(b, 10, 10, 9, 5, Radius(2), 0xFFFFFFFF)
	case 12:
		c.BitmapPart(b, 4, 2, 8, 8, 100, 10, 40, 40, 0xFFFFFFFF)
	default:
		c.Clip(20, 15, 80, 50, func() {
			c.RRect(0, 0, 160, 100, Radius(30), 0xFF4080FF)
			c.Clip(50, 0, 100, 40, func() {
				c.Circle(70, 40, 30, 0xFFFF4040)
				c.Text(30, 30, "clip", true, 0xFF000000)
				c.Shadow(40, 20, 60, 30, 8, 10, 0xFF000000)
				c.Bitmap(b, 40, 10, 48, 32, Radius(6), 0xFFFFFFFF)
			})
			c.Line(0, 0, 160, 100, 3, 0xFF00A000)
			c.Poly(star, 0x80000000)
		})
		c.Frame(20, 15, 80, 50, 1, 0x80000000)
	}
}

func testBitmap(t *testing.T) *Bitmap {
	b, err := NewBitmap(24, 16, 27)
	if err != nil {
		t.Fatal(err)
	}
	px := b.Pixels()
	for y := 0; y < 16; y++ {
		for x := 0; x < 24; x++ {
			a := uint32(x * 255 / 23)
			v := uint32(0x2050C0)
			if (x/4+y/4)&1 != 0 {
				v = 0xFFFFFF
			}
			px[y*27+x] = a<<24 | (v>>16&255)*a/255<<16 | (v>>8&255)*a/255<<8 | (v&255)*a/255
		}
	}
	return b
}

func newCanvas(t *testing.T, w, h, stride int, scale float64) *Canvas {
	c, err := NewCanvas(w, h, stride, scale)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func zeroCols(px []uint32, stride, h, x0, x1 int) bool {
	for y := 0; y < h; y++ {
		for _, p := range px[y*stride+x0 : y*stride+x1] {
			if p != 0 {
				return false
			}
		}
	}
	return true
}

func TestBasic(t *testing.T) {
	c := newCanvas(t, 160, 100, 200, 1.5)
	defer c.Free()
	c.RoundRect(5, 5, 60, 40, 8, 0xC0203040)
	c.Line(0, 0, 100, 60, 3, 0xFFFF0000)
	c.Triangle(10, 50, 60, 50, 35, 10, 0x8000FF00)
	c.Text(8, 8, "goome!", true, 0xFFFFFFFF)
	c.Frame(70, 10, 30, 30, 2, 0xFF00FFFF)
	px := c.Pixels()
	for y := 0; y < 100; y++ {
		if px[y*200+170] != 0 {
			t.Fatal("drew past w")
		}
	}
	x0, y0, x1, y1 := c.Drawn()
	if x1 <= x0 || y1 <= y0 {
		t.Error("nothing drawn")
	}
	if h := fnv(px, 200, 100, 200); h != 0xB93B60CF {
		t.Errorf("basic hash %08x", h)
	}
	c.Clear()
	if h := fnv(c.Pixels(), 200, 100, 200); h != fnv(make([]uint32, 200*100), 200, 100, 200) {
		t.Error("clear left pixels")
	}
	if c.TextWidth("abc", false) <= 0 {
		t.Error("no text width")
	}
}

func TestShapes(t *testing.T) {
	b := testBitmap(t)
	defer b.Free()
	for s, sc := range []float64{1, 1.25, 2} {
		var got []string
		for k := 0; k < 14; k++ {
			w, h := int(160*sc), int(100*sc)
			c := newCanvas(t, w, h, w+3, sc)
			shape(c, k, b)
			px := c.Pixels()
			got = append(got, fmt.Sprintf("%08x", fnv(px, w+3, h, w+3)))
			if x0, _, x1, _ := c.Drawn(); x1 <= x0 {
				t.Errorf("shape %d drew nothing", k)
			}
			if !zeroCols(px, w+3, h, w, w+3) {
				t.Errorf("shape %d drew past w", k)
			}
			for _, p := range px {
				a := p >> 24
				if p>>16&255 > a || p>>8&255 > a || p&255 > a {
					t.Fatalf("shape %d not premultiplied", k)
				}
			}
			c.Free()
			if k == 13 {
				cx0, cy0 := int(20*sc+0.5), int(15*sc+0.5)
				cx1, cy1 := int(100*sc+0.5), int(65*sc+0.5)
				c := newCanvas(t, w, h, w+3, sc)
				c.ClipPush(20, 15, 80, 50)
				shape(c, 8, b)
				shape(c, 4, b)
				shape(c, 9, b)
				c.Text(0, 20, "outside", true, 0xFFFFFFFF)
				x0, y0, x1, y1 := c.Drawn()
				if x0 < cx0 || y0 < cy0 || x1 > cx1 || y1 > cy1 || x1 <= x0 {
					t.Errorf("drawn box %d %d %d %d outside the clip", x0, y0, x1, y1)
				}
				px := c.Pixels()
				for y := 0; y < h; y++ {
					for x := 0; x < w; x++ {
						if (x < cx0 || x >= cx1 || y < cy0 || y >= cy1) && px[y*(w+3)+x] != 0 {
							t.Fatal("drew outside the clip")
						}
					}
				}
				c.ClipPop()
				if c.ClipDepth() != 0 {
					t.Error("clip still pushed")
				}
				c.Free()
			}
		}
		var want []string
		for _, v := range shapeHashes[s] {
			want = append(want, fmt.Sprintf("%08x", v))
		}
		if strings.Join(got, " ") != strings.Join(want, " ") {
			t.Errorf("scale %g\n got %v\nwant %v", sc, got, want)
		}
	}
}

func TestImageAndBlit(t *testing.T) {
	src := newCanvas(t, 20, 20, 0, 1)
	defer src.Free()
	green := make([]uint32, 16)
	for i := range green {
		green[i] = 0xFF00FF00
	}
	if err := src.Image(0, 0, 20, 20, green, 4, 4); err != nil {
		t.Fatal(err)
	}
	if src.Pixels()[0] != 0xFF00FF00 {
		t.Error("image didnt draw")
	}
	dst := newCanvas(t, 40, 40, 0, 1)
	defer dst.Free()
	dst.Blit(src, 1, 10, 10)
	if dst.Pixels()[10*40+10] != 0xFF00FF00 || dst.Pixels()[0] != 0 {
		t.Error("blit wrong")
	}
	dst.Erase(0, 0, 40, 40)
	if dst.Pixels()[10*40+10] != 0 {
		t.Error("erase didnt")
	}
	if err := dst.Image(0, 0, 1, 1, green[:3], 2, 2); !errors.Is(err, ErrArg) {
		t.Errorf("short image: %v", err)
	}
}

func TestFonts(t *testing.T) {
	ttf, otf := data(t, "dejavu-sub.ttf"), data(t, "dejavu-sub.otf")
	if _, err := FontFromBytes([]byte("nope, not a font")); err == nil {
		t.Error("loaded junk")
	}
	if _, err := FontFromBytes(ttf[:40]); err == nil {
		t.Error("loaded a cut font")
	}
	f, err := FontFromBytes(ttf)
	if err != nil {
		t.Fatal(err)
	}
	o, err := FontFromBytes(otf)
	if err != nil {
		t.Fatal(err)
	}
	w, h := f.Measure("goome!", 20)
	if w < 76.97 || w > 76.98 || h < 23.2 || h > 23.3 {
		t.Errorf("measure %v %v", w, h)
	}
	for _, font := range []*Font{f, o} {
		if font.Width("AV", 20) >= font.Width("A", 20)+font.Width("V", 20)-1 {
			t.Error("no kerning")
		}
	}
	if _, h := f.Measure("a\nbb\n", 20); h < 3*23.2 || h > 3*23.3 {
		t.Errorf("3 lines %v", h)
	}
	if f.Width("\xff\xc3(\xed\xa0\x80\xf4\x90\x80\x80", 20) <= 0 || f.LineHeight(20) <= 0 {
		t.Error("bad utf8 or line height")
	}
	want := []uint32{0x412C299F, 0xBCD42095}
	var bufs [][]uint32
	for k, font := range []*Font{f, o} {
		c := newCanvas(t, 200, 120, 240, 1.5)
		c.SetClipPx(0, 0, 200, 120)
		c.FontText(font, 13, 4, 2, "goome! Hié\nAVAWAY Ty.@ ☃", 0xFFFFFFFF)
		c.FontText(font, 30, 2.3, 24, "Hig", 0xC0FF8040)
		c.SetClipPx(30, 60, 90, 100)
		c.SetClipDepth(1)
		c.FontText(font, 40, 10, 50, "WoW", 0xFF2080FF)
		px := c.Pixels()
		if !zeroCols(px, 240, 120, 200, 240) {
			t.Error("text past w")
		}
		for _, p := range px[100*240:] {
			if p != 0 {
				t.Fatal("text below the clip")
			}
		}
		if hv := fnv(px, 240, 120, 240); hv != want[k] {
			t.Errorf("font %d hash %08x want %08x", k, hv, want[k])
		}
		c.Free()
		c = newCanvas(t, 240, 120, 0, 1)
		c.SetClipPx(0, 0, 240, 120)
		c.FontText(font, 31, 3, 3, "goome Hi\nmoge", 0xFFFFFFFF)
		bufs = append(bufs, append([]uint32(nil), c.Pixels()...))
		c.Free()
	}
	ink, diff := 0, 0
	for i := range bufs[0] {
		a, b := int(bufs[0][i]>>24), int(bufs[1][i]>>24)
		ink += a
		if a > b {
			diff += a - b
		} else {
			diff += b - a
		}
	}
	if ink == 0 || diff*50 >= ink {
		t.Errorf("ttf and otf differ too much: ink %d diff %d", ink, diff)
	}
	c := newCanvas(t, 200, 120, 240, 1)
	c.SetClipPx(0, 0, 200, 120)
	for cut := 0; cut < len(otf); cut += 7 {
		bad := append([]byte(nil), otf...)
		bad[cut] ^= 0x5A
		if cut&1 == 0 && cut+12 < len(bad) {
			bad = bad[:cut+12]
		}
		bf, _ := FontFromBytes(bad)
		c.FontText(bf, 20, 0, 0, "goome! AV", 0xFFFFFFFF)
		bf.Free()
	}
	c.Free()
	c = newCanvas(t, 200, 120, 240, 1)
	var builtin *Font
	if got, w := c.FontText(nil, 12, 0, 0, "ab\nabc", 0xFFFFFFFF), builtin.Width("abc", 12); got != w {
		t.Errorf("builtin %v %v", got, w)
	}
	c.Free()
	f.Free()
	f.Free()
	o.Free()
	if FontCount() < 1 {
		t.Error("no builtin fonts")
	}
	SetFont(0)
	if builtin.LineHeight(12) <= 0 {
		t.Error("builtin line height")
	}
	p := filepath.Join(t.TempDir(), "f.ttf")
	os.WriteFile(p, ttf, 0o644)
	for _, load := range []func() (*Font, error){func() (*Font, error) { return LoadFont(p) }, func() (*Font, error) { return LoadFontIndex(p, 0) }} {
		lf, err := load()
		if err != nil || lf.Width("goome", 12) <= 0 {
			t.Errorf("load file %v", err)
		}
		lf.Free()
	}
	if _, err := LoadFont(filepath.Join(t.TempDir(), "gone.ttf")); err == nil {
		t.Error("loaded a missing font")
	}
}

func TestBitmaps(t *testing.T) {
	for _, ic := range imageCases {
		raw := data(t, ic.name)
		b, err := BitmapFromBytes(raw)
		if err != nil {
			t.Errorf("%s: %v", ic.name, err)
			continue
		}
		if b.W() != ic.w || b.H() != ic.h || !b.Owned() {
			t.Errorf("%s: %dx%d owned %v", ic.name, b.W(), b.H(), b.Owned())
		}
		if hv := fnv(b.Pixels(), b.W(), b.H(), b.Stride()); hv != ic.hash {
			t.Errorf("%s: hash %08x want %08x", ic.name, hv, ic.hash)
		}
		b.Free()
		b.Free()
		if !b.Freed() {
			t.Error("not freed")
		}
		if half, err := BitmapFromBytes(raw[:len(raw)/2]); err == nil {
			if half.W() <= 0 {
				t.Errorf("%s half has no size", ic.name)
			}
			half.Free()
		}
	}
	_, err := BitmapFromBytes(data(t, "jpeg-arithmetic.jpg"))
	if !errors.Is(err, ErrNotSupported) || !strings.Contains(err.Error(), "arithmetic") {
		t.Errorf("arithmetic jpeg: %v", err)
	}
	if _, err := BitmapFromBytes([]byte("GIF89a......")); !errors.Is(err, ErrNotSupported) {
		t.Errorf("gif: %v", err)
	}
	if _, err := BitmapFromBytes([]byte("\x89PNG\r\n\x1a\n")); err == nil {
		t.Error("loaded a cut png")
	}
	if _, err := LoadBitmap(filepath.Join("no", "such", "file.png")); err == nil {
		t.Error("loaded a missing file")
	}
	b, err := BitmapFromARGB([]uint32{0xFFFF0000, 0x80FFFFFF, 0x00123456, 0x40804020}, 2, 2)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprint(b.Pixels()) != fmt.Sprint([]uint32{0xFFFF0000, 0x80808080, 0, 0x40201008}) {
		t.Errorf("from argb %x", b.Pixels())
	}
	b.Free()
	if _, err := BitmapFromARGB([]uint32{1, 2}, 2, 2); !errors.Is(err, ErrArg) {
		t.Errorf("short argb: %v", err)
	}
	p := filepath.Join(t.TempDir(), "x.qoi")
	os.WriteFile(p, data(t, "qoi.qoi"), 0o644)
	lb, err := LoadBitmap(p)
	if err != nil || lb.W() != 6 || lb.H() != 4 {
		t.Errorf("load file %v", err)
	}
	lb.Free()
}
