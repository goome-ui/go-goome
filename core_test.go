package goome

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseCombo(t *testing.T) {
	cases := []struct {
		s         string
		key, mods int
	}{
		{"alt+h", 'h', ModAlt},
		{"Ctrl+Shift+F5", KeyF1 + 4, ModCtrl | ModShift},
		{"alt+?", '/', ModAlt | ModShift},
		{"super+space", ' ', ModSuper},
		{"cmd+left", KeyLeft, ModSuper},
		{"H", 'h', ModShift},
	}
	for _, c := range cases {
		k, m, err := ParseCombo(c.s)
		if err != nil || k != c.key || m != c.mods {
			t.Errorf("%s: got %d %x %v", c.s, k, m, err)
		}
	}
	for _, bad := range []string{"alt+", "alt+f13", ""} {
		if _, _, err := ParseCombo(bad); err == nil {
			t.Errorf("%q should fail", bad)
		}
	}
}

func TestPacer(t *testing.T) {
	p := NewPacer(50)
	defer p.Free()
	if w := p.Wait(); w != 0 {
		t.Errorf("first wait %d", w)
	}
	if w := p.Wait(); w <= 0 || w > 20 {
		t.Errorf("second wait %d", w)
	}
	a := NowMs()
	if NowMs() < a {
		t.Error("time went back")
	}
}

const script = "screen 640x360 2; 1 click 20,20; 1 click 300,300; 2 hotkey alt+j; 3 key Hi\\b; 4 send ping; 5 scroll 30,30 -1; 6 hold 25,25 8; end 10"

func TestScripted(t *testing.T) {
	g, err := Open(Options{Title: "t", Backend: BackendTest, TestScript: script})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(g.BackendName(), "test") {
		t.Errorf("backend %q", g.BackendName())
	}
	mons := g.Monitors()
	if len(mons) != 1 || mons[0].Area.W != 640 || mons[0].Scale != 2 || mons[0].PW != 1280 {
		t.Fatalf("monitors %+v", mons)
	}
	if err := g.Hotkey("alt+j", 7); err != nil {
		t.Fatal(err)
	}
	g.Keyboard(true)
	down, up, hot, keys, ctl, scroll, held := 0, 0, 0, 0, 0, 0, 0
	text, quit, frames := "", false, 0
	for !quit && frames < 100 {
		c := g.Begin(0)
		if c == nil {
			t.Fatal("no canvas")
		}
		if c.Scale() != 2 || c.W() != 1280 {
			t.Fatalf("canvas %v %d", c.Scale(), c.W())
		}
		c.Rect(10, 10, 40, 40, 0xFF336699)
		c.Hit(10, 10, 40, 40)
		if err := g.End(c); err != nil {
			t.Fatal(err)
		}
		if c.Valid() {
			t.Fatal("canvas still valid after End")
		}
		frames++
		if _, _, b, err := g.Pointer(); err == nil && b != 0 {
			held++
		}
		for ev := range g.Events(1) {
			switch {
			case ev.Type == EventDown:
				down++
			case ev.Type == EventUp:
				up++
			case ev.Type == EventHotkey && ev.ID == 7:
				hot++
			case ev.Type == EventKey:
				keys++
				text += ev.Text
			case ev.Type == EventControl && ev.Msg == "ping":
				ctl++
			case ev.Type == EventScroll && ev.DY == -1:
				scroll++
			case ev.Type == EventQuit:
				quit = true
			}
		}
	}
	if down != 2 || up != 2 || hot != 1 || keys != 3 || text != "Hi" || ctl != 1 || scroll != 1 || !quit || held < 1 {
		t.Errorf("down %d up %d hot %d keys %d text %q ctl %d scroll %d quit %v held %d", down, up, hot, keys, text, ctl, scroll, quit, held)
	}
	img, err := g.Capture(0, &Rect{10, 10, 40, 40})
	if err != nil {
		t.Fatal(err)
	}
	if img.W() != 80 || img.H() != 80 || !img.HasSelf() {
		t.Errorf("capture %dx%d self %v", img.W(), img.H(), img.HasSelf())
	}
	if img.Pixel(0, 0) != 0xFF336699 || img.Pixels()[0] != 0xFF336699 || len(img.Pixels()) != 80*img.Stride() {
		t.Errorf("capture pixels wrong %08x", img.Pixel(0, 0))
	}
	kept := img.Copy()
	if len(kept) != 80*80 || kept[0] != 0xFF336699 {
		t.Error("copy wrong")
	}
	if g.AlphaAt(20, 20) != 255 || g.AlphaAt(200, 200) != 0 {
		t.Errorf("alpha_at %d %d", g.AlphaAt(20, 20), g.AlphaAt(200, 200))
	}
	again, err := g.Capture(0, nil)
	if err != nil {
		t.Fatal(err)
	}
	if img.Fresh() || img.Pixels() != nil {
		t.Error("old grab still fresh")
	}
	if again.W() != 1280 || again.H() != 720 {
		t.Errorf("full grab %dx%d", again.W(), again.H())
	}
	p := filepath.Join(t.TempDir(), "grab.ppm")
	if err := again.SavePPM(p); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(p); len(b) < 2 || string(b[:2]) != "P6" {
		t.Error("not a ppm")
	}
	g.Close()
	g.Close()
	if !g.Closed() || again.Fresh() {
		t.Error("still open")
	}
}

func TestErrors(t *testing.T) {
	_, err := Open(Options{Backend: BackendTest, TestFrames: filepath.Join("no", "such", "frames")})
	if !errors.Is(err, ErrNoDisplay) {
		t.Errorf("got %v", err)
	}
	var ge *Error
	if !errors.As(err, &ge) || ge.Code != CodeNoDisplay {
		t.Errorf("not a goome error: %v", err)
	}
	if Strerror(CodeNotSupported) == "" {
		t.Error("no strerror")
	}
	g, err := Open(Options{Backend: BackendTest, TestScript: "end 3"})
	if err != nil {
		t.Fatal(err)
	}
	if g.ControlPath() != "" {
		t.Error("control path without a state dir")
	}
	err = g.Frame(0, func(c *Canvas) {
		c.Rect(0, 0, 10, 10, 0xFFFF0000)
		c.Hit(0, 0, 10, 10)
	})
	if err != nil {
		t.Fatal(err)
	}
	if g.AlphaAt(5, 5) != 255 {
		t.Error("frame didnt land")
	}
	c := g.Begin(0)
	g.Close()
	if c.Valid() {
		t.Error("canvas valid after close")
	}
	func() {
		defer func() {
			r := recover()
			if e, ok := r.(*Error); !ok || e.Code != CodeGone {
				t.Errorf("drawing on a dead canvas: %v", r)
			}
		}()
		c.Rect(0, 0, 1, 1, 0xFFFFFFFF)
	}()
	if _, _, err := g.Poll(0); !errors.Is(err, ErrGone) {
		t.Errorf("poll after close: %v", err)
	}
}

func TestControl(t *testing.T) {
	d := t.TempDir()
	g, err := Open(Options{Backend: BackendTest, StateDir: d, TestScript: "end 50"})
	if err != nil {
		t.Skipf("no state dir here: %v", err)
	}
	defer g.Close()
	if g.ControlPath() == "" {
		t.Skip("no control fifo on this os")
	}
	if err := Send(d, "hello there"); err != nil {
		t.Fatal(err)
	}
	got := ""
	for i := 0; i < 40 && got == ""; i++ {
		g.Frame(0, func(*Canvas) {})
		for ev := range g.Events(5) {
			if ev.Type == EventControl {
				got = ev.Msg
			}
		}
	}
	if got != "hello there" {
		t.Errorf("got %q", got)
	}
}

func TestOffscreen(t *testing.T) {
	g, err := OpenOffscreen(200, 100, 1)
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	g.Frame(0, func(c *Canvas) {
		c.Rect(10, 10, 50, 20, 0xFF00FF00)
		c.Hit(10, 10, 50, 20)
	})
	f, err := g.LastFrame()
	if err != nil {
		t.Fatal(err)
	}
	if f.W() != 200 || f.Pixel(20, 15) != 0xFF00FF00 {
		t.Errorf("frame %d %08x", f.W(), f.Pixel(20, 15))
	}
	if r, ok := g.FrameDirty(); !ok || r.W != 200 {
		t.Errorf("first dirty %v %v", r, ok)
	}
	buf := make([]byte, 50*20*4)
	if err := g.FrameCopy(buf, 50*4, FmtRGBA, &Rect{10, 10, 50, 20}); err != nil {
		t.Fatal(err)
	}
	if buf[0] != 0 || buf[1] != 255 || buf[3] != 255 {
		t.Errorf("rgba %v", buf[:4])
	}
	if err := g.FrameCopy(buf[:10], 50*4, FmtRGBA, &Rect{10, 10, 50, 20}); !errors.Is(err, ErrArg) {
		t.Errorf("short buffer: %v", err)
	}
	if !g.Hit(20, 15) || g.Hit(150, 80) {
		t.Error("hit")
	}
	if err := g.Push(Event{Type: EventMotion, X: 5, Y: 5}); err != nil {
		t.Fatal(err)
	}
	if err := g.FrameCopy(buf, 50*4, FmtRGBA, &Rect{1<<31 - 1, 0, 1, 1}); !errors.Is(err, ErrArg) {
		t.Errorf("rect off the frame: %v", err)
	}
	if err := SavePPM(filepath.Join(t.TempDir(), "x.ppm"), make([]uint32, 10), 1, 4, 0x5555555555555556); !errors.Is(err, ErrArg) {
		t.Errorf("stride that wraps: %v", err)
	}
	if err := g.Resize(300, 200, 2); err != nil {
		t.Fatal(err)
	}
	if f.Fresh() || f.Pixels() != nil {
		t.Error("frame from before the resize still hands out pixels")
	}
	g.SetClipboard("hi")
	if s, err := g.Clipboard(); err != nil || s != "hi" {
		t.Errorf("clipboard %q %v", s, err)
	}
}

func TestFrameMatchesC(t *testing.T) {
	want := strings.TrimSpace(runC(t, "testdata/frame.c"))
	g, err := Open(Options{Mode: ModeOffscreen, Title: "engine", Width: 200, Height: 200, Scale: 2})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	c := g.Begin(0)
	c.Rect(10, 10, 20, 20, 0x80FF0000)
	c.Circle(120, 50, 30, 0xFF3A7AFE)
	c.RRect(40, 60, 70, 30, Radius(8), 0xCC20C060)
	c.Text(44, 66, "offscreen", false, 0xFFFFFFFF)
	g.End(c, Rect{10, 10, 20, 20})
	f, err := g.LastFrame()
	if err != nil {
		t.Fatal(err)
	}
	lines := []string{fmt.Sprintf("draw %08x", fnv(f.Pixels(), f.W(), f.H(), f.Stride()))}
	u, _ := g.NewUI()
	u.Theme().AnimMs = 0
	u.Theme().SetFont(FontMain, nil)
	u.Theme().FontSize = 14
	win := &Window{X: 6, Y: 6, W: 180}
	val, check, name := 0.25, false, "hi"
	for frame := 0; frame < 6; frame++ {
		switch frame {
		case 1:
			g.PushClick(25, 64)
		case 2:
			g.Push(Event{Type: EventDown, X: 90, Y: 168, Button: 1})
			g.Push(Event{Type: EventUp, X: 90, Y: 168, Button: 1})
		case 3:
			g.Push(Event{Type: EventKey, Key: 'z', Text: "z"})
		}
		for ev := range g.Events(0) {
			u.Feed(ev)
		}
		if u.Begin(0) == nil {
			u.End()
			continue
		}
		if u.Window("engine", win, WinMove) {
			u.Checkbox("Check", &check)
			u.SliderFmt("Val", &val, 0, 1, "%.2f")
			u.InputEx("Name", &name, 0, "", 32)
		}
		u.End()
		u.End()
	}
	f, _ = g.LastFrame()
	b, h := 0, 0
	if check {
		b = 1
	}
	if g.Hit(10, 10) {
		h = 1
	}
	lines = append(lines, fmt.Sprintf("ui %08x %d %s %d", fnv(f.Pixels(), f.W(), f.H(), f.Stride()), b, name, h))
	if got := strings.Join(lines, "\n"); got != want {
		t.Errorf("go:\n%s\nC:\n%s", got, want)
	}
}
