package goome

import (
	"fmt"
	"strings"
	"testing"
)

type animScene struct {
	g        *Goome
	u        *UI
	on       bool
	box, sp  float64
	col, sum uint32
}

func (k *animScene) frame() {
	for ev := range k.g.Events(0) {
		k.u.Feed(ev)
	}
	c := k.u.Begin(0)
	if c == nil {
		k.u.End()
		return
	}
	px, w, h, st := c.Pixels(), c.W(), c.H(), c.Stride()
	r := k.u.Canvas()
	pick := func(a, b float64) float64 {
		if k.on {
			return a
		}
		return b
	}
	k.box = k.u.Anim("box", pick(200, 10), 160, EaseOut)
	k.sp = k.u.Anim("spring", pick(1, 0), 320, EaseSpring)
	col := uint32(0)
	if k.on {
		col = 0xFF20D060
	}
	k.col = k.u.AnimColor("col", col, 200, EaseLinear)
	k.u.Anim("gone", pick(5, 0), 100, EaseLinear)
	if r != nil {
		r.Rect(k.box, 20, 60, 40, 0xFFE0A030)
		r.RRect(20+k.sp*180, 90, 50, 30, Radius(15), 0xFF4FB3F7)
		r.Rect(20, 140, 260, 40, k.col)
	}
	k.u.End()
	k.sum = fnv(px, w, h, st)
}

func animLines(t *testing.T) ([]string, [4]uint32) {
	g, err := Open(Options{Backend: BackendTest, Mode: ModeWindow, Width: 300, Height: 200, TestScript: "end 100000"})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	u, _ := g.NewUI()
	defer u.Free()
	k := &animScene{g: g, u: u}
	k.frame()
	k.frame()
	if k.box != 10 || k.sp != 0 || k.col != 0 || u.Animating() {
		t.Errorf("anim didnt start at rest: %v %v %08x", k.box, k.sp, k.col)
	}
	var hs [4]uint32
	k.on = true
	moving, prev := 0, k.box
	for f := 0; f < 60; f++ {
		k.frame()
		switch f {
		case 1:
			hs[0] = k.sum
		case 4:
			hs[1] = k.sum
		case 8:
			hs[2] = k.sum
		}
		if k.box < prev {
			t.Errorf("box went back at frame %d", f)
		}
		prev = k.box
		if u.Animating() {
			moving = f + 1
		}
	}
	hs[3] = k.sum
	return []string{
		fmt.Sprintf("anim %08x %08x %08x %08x %d %.0f %.0f %08x", hs[0], hs[1], hs[2], hs[3], moving, k.box, k.sp, k.col),
		fmt.Sprintf("ease %.6f %.6f %.6f %.6f", Ease(EaseIn, 0.5), Ease(EaseBounce, 0.25), Ease(EaseSpring, 0.5), Ease(EaseBack, 0.3)),
	}, hs
}

func textHash(c *Canvas, f *Font, size float64, s string, argb uint32, dir int) uint32 {
	full := make([]uint32, 400*120)
	SetTextDir(dir)
	c.FontText(f, size, 4, 4, s, argb)
	SetTextDir(DirAuto)
	copy(full, c.Pixels())
	return fnv(full, 400, 120, 400)
}

func newTextCanvas(t *testing.T, w, h int) *Canvas {
	c, err := NewCanvas(w, h, 400, 1)
	if err != nil {
		t.Fatal(err)
	}
	c.SetClipPx(0, 0, w, h)
	return c
}

func font(t *testing.T, name string) *Font {
	f, err := FontFromBytes(data(t, name))
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func emojiLines(t *testing.T) ([]string, [6]uint32) {
	f, cb, sb, d := font(t, "colr.ttf"), font(t, "cbdt.ttf"), font(t, "sbix.ttf"), font(t, "dejavu-sub.ttf")
	defer f.Free()
	defer cb.Free()
	defer sb.Free()
	defer d.Free()
	var got [6]uint32
	got[0] = textHash(newTextCanvas(t, 400, 60), f, 30, "ab 😀 😀🏻 😀‍😀 ❤ ❤️", 0xFF2040A0, DirAuto)
	c, _ := NewCanvas(400, 120, 400, 1)
	c.FontText(cb, 16, 4, 4, "🐱", 0xFFFFFFFF)
	c.FontText(cb, 64, 30, 4, "🐱", 0xFFFFFFFF)
	c.FontText(sb, 40, 120, 4, "🐱🐱", 0x80FFFFFF)
	mid := c.Pixels()[(4+40)*400+30+35]
	if mid>>24 != 255 || mid>>16&255 <= 200 {
		t.Errorf("cbdt cat pixel %08x", mid)
	}
	got[1] = fnv(c.Pixels(), 400, 120, 400)
	c.Free()
	plain := d.Width("goome! ", 20)
	d.Fallback(f)
	d.Fallback(d)
	d.Fallback(f)
	if len(d.Fallbacks()) != 1 {
		t.Errorf("fallbacks %d", len(d.Fallbacks()))
	}
	if w := d.Width("goome! 😀", 20); w-plain-16 > 0.01 || plain+16-w > 0.01 {
		t.Errorf("fallback emoji width %v", w-plain)
	}
	heb := "ab \u05d0\u05d1\u05d2 (12) c"
	got[2] = textHash(newTextCanvas(t, 400, 40), d, 20, "goome! 😀 Hi", 0xFFFFFFFF, DirAuto)
	got[3] = textHash(newTextCanvas(t, 400, 40), f, 24, heb, 0xFFFFFFFF, DirAuto)
	if got[3] != textHash(newTextCanvas(t, 400, 40), f, 24, heb, 0xFFFFFFFF, DirLTR) {
		t.Error("ltr should look like auto here")
	}
	got[4] = textHash(newTextCanvas(t, 400, 40), f, 24, heb, 0xFFFFFFFF, DirRTL)
	got[5] = textHash(newTextCanvas(t, 400, 50), f, 30, "\u0628\u0644\u0627\u0645 \u0628\u0628\u0628 \u0644\u0645 \u0628\u064e\u0628", 0xFFFFFFFF, DirAuto)
	return []string{
		fmt.Sprintf("emoji %08x %08x %08x %08x %08x %08x", got[0], got[1], got[2], got[3], got[4], got[5]),
		fmt.Sprintf("measure %.2f %.2f %.2f", f.Width("😀🏻", 30), d.Width("goome! 😀", 20), f.Width("\u0628\u0644\u0627\u0645 \u0628\u0628\u0628 \u0644\u0645", 20)),
	}, got
}

type imeScene struct {
	g      *Goome
	u      *UI
	in, ed string
	sum    uint32
}

func (k *imeScene) frame(focus int) {
	for ev := range k.g.Events(0) {
		k.u.Feed(ev)
	}
	c := k.u.Begin(0)
	if c == nil {
		k.u.End()
		return
	}
	px, w, h, st := c.Pixels(), c.W(), c.H(), c.Stride()
	if focus == 1 {
		k.u.SetFocus()
	}
	k.u.InputEx("in", &k.in, InputNoSelectOnFocus, "", 64)
	if focus == 2 {
		k.u.SetFocus()
	}
	k.u.TextEdit("ed", &k.ed, 80, 0)
	k.u.End()
	k.sum = fnv(px, w, h, st)
}

func openIME(t *testing.T, in, ed string, focus int) *imeScene {
	g, err := Open(Options{Backend: BackendTest, Mode: ModeWindow, Width: 300, Height: 200, TestScript: "end 100000"})
	if err != nil {
		t.Fatal(err)
	}
	u, _ := g.NewUI()
	u.Theme().SetFont(FontMain, nil)
	u.Theme().FontSize = 14
	u.Theme().AnimMs = 0
	u.Config().BlinkMs = 0
	k := &imeScene{g: g, u: u, in: in, ed: ed}
	k.frame(focus)
	k.frame(0)
	return k
}

func (k *imeScene) push(typ int, s string, cur int) {
	if typ == EventText {
		k.g.PushIME(s)
	} else {
		k.g.PushPreedit(s, cur)
	}
	k.frame(0)
}

func (k *imeScene) close() {
	k.u.Free()
	k.g.Close()
}

func imeLines(t *testing.T) ([]string, [3]uint32) {
	k := openIME(t, "ab", "", 1)
	var h [3]uint32
	h0 := k.sum
	k.push(EventPreedit, "xyz", 1)
	h[0] = k.sum
	k.push(EventPreedit, "xyz", 3)
	h[1] = k.sum
	if k.in != "ab" {
		t.Errorf("preedit changed the text: %q", k.in)
	}
	k.push(EventPreedit, "", 0)
	back := 0
	if k.sum == h0 {
		back = 1
	}
	k.push(EventText, "日本", 0)
	lines := []string{fmt.Sprintf("ime %08x %08x %08x %d %s", h0, h[0], h[1], back, k.in)}
	k.close()
	k = openIME(t, "", "line", 2)
	k.push(EventPreedit, "kana", 2)
	h[2] = k.sum
	k.push(EventPreedit, "", 0)
	k.push(EventText, "かな\n", 0)
	ok := 0
	if k.ed == "かな\nline" {
		ok = 1
	}
	lines = append(lines, fmt.Sprintf("edit %08x %d", h[2], ok))
	k.close()
	return lines, h
}

type a11yScene struct {
	g                *Goome
	u                *UI
	ok, mode, fruit  int
	tab              int
	dark, open       bool
	vol              float64
	name, note       string
	nodes            []A11yNode
	knobHits, frames int
}

func (k *a11yScene) frame() {
	for ev := range k.g.Events(0) {
		k.u.Feed(ev)
	}
	u := k.u
	if u.Begin(0) == nil {
		u.End()
		return
	}
	u.App("Demo", AppTitlebar|AppStatusbar)
	u.Label("Hello there")
	if u.Button("OK") {
		k.ok++
	}
	u.Checkbox("Dark", &k.dark)
	u.Radio("Big", &k.mode, 1)
	u.SliderFmt("Volume", &k.vol, 0, 10, "%.1f")
	u.InputEx("Name", &k.name, 0, "", 64)
	u.Combo("Fruit", &k.fruit, []string{"Mango", "Kiwi", "Fig"})
	u.Tabs("tabs", []string{"One", "Two"}, &k.tab)
	if u.Tree("Folder", &k.open) {
		u.Label("inside")
	}
	u.End()
	u.Progress(0.25, "")
	u.TextEdit("Notes", &k.note, 60, 0)
	u.A11yRole(RoleButton, "Knob")
	u.A11yValue("11")
	if it, _ := u.Custom("knob", 40, 40, 0); it.Clicked {
		k.knobHits++
	}
	u.Statusbar("ready")
	u.End()
	u.End()
	k.nodes = u.A11yTree()
}

func (k *a11yScene) find(role int, label string) int {
	for i, n := range k.nodes {
		if n.Role == role && n.Label == label {
			return i
		}
	}
	return -1
}

func (k *a11yScene) act(verb string, role int, label, value string) {
	if i := k.find(role, label); i >= 0 {
		k.g.A11yAction(verb, k.nodes[i], value)
	}
}

func a11yLines(t *testing.T) []string {
	g, err := Open(Options{Backend: BackendTest, Mode: ModeWindow, Width: 420, Height: 700, TestScript: "end 100000"})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	u, _ := g.NewUI()
	defer u.Free()
	u.Theme().SetFont(FontMain, nil)
	k := &a11yScene{g: g, u: u, vol: 2, name: "Bo"}
	k.frame()
	lines := []string{fmt.Sprintf("a11y %d", len(k.nodes))}
	for _, n := range k.nodes {
		lines = append(lines, fmt.Sprintf("node %d %d %d %.1f %.1f %.1f %.1f [%s] [%s] %.2f %.2f %.2f %d",
			n.Role, n.Parent, n.States, n.X, n.Y, n.W, n.H, n.Label, n.Value, n.Num, n.Min, n.Max, n.Caret))
	}
	if len(k.nodes) < 15 || k.nodes[0].Parent != -1 || k.nodes[0].Role != RoleWindow {
		t.Fatalf("tree looks wrong: %+v", k.nodes)
	}
	if i := k.find(RoleTab, "One"); i < 0 || !k.nodes[i].Is(A11ySelected) || k.nodes[k.nodes[i].Parent].Role != RoleGroup {
		t.Error("tab One isnt selected under a group")
	}
	if i := k.find(RoleButton, "Knob"); i < 0 || k.nodes[i].Value != "11" {
		t.Error("custom knob role/value missing")
	}
	k.act("press", RoleButton, "OK", "")
	k.act("press", RoleCheckbox, "Dark", "")
	k.act("set", RoleSlider, "Volume", "7.5")
	k.act("set", RoleInput, "Name", "hello there")
	k.act("expand", RoleTree, "Folder", "")
	k.act("press", RoleTab, "Two", "")
	k.frame()
	k.frame()
	b := func(v bool) int {
		if v {
			return 1
		}
		return 0
	}
	lines = append(lines, fmt.Sprintf("acted %d %d %.1f %s %d %d %d %d", k.ok, b(k.dark), k.vol, k.name, b(k.open), k.tab, len(k.nodes), k.find(RoleText, "inside")))
	if i := k.find(RoleText, "inside"); i < 0 || k.nodes[k.nodes[i].Parent].Label != "Folder" {
		t.Error("tree child not under the tree")
	}
	if i := k.find(RoleCheckbox, "Dark"); i < 0 || !k.nodes[i].Is(A11yChecked) {
		t.Error("checkbox not checked after press")
	}
	k.act("focus", RoleInput, "Name", "")
	k.frame()
	k.frame()
	if i := k.find(RoleInput, "Name"); i < 0 || !k.nodes[i].Is(A11yFocused) || k.nodes[i].Caret != 11 {
		t.Error("focus action didnt land")
	}
	k.act("press", RoleButton, "Knob", "")
	k.frame()
	if k.knobHits != 1 {
		t.Errorf("knob hits %d", k.knobHits)
	}
	u.A11y(false)
	k.frame()
	if len(k.nodes) != 0 {
		t.Error("a11y off still has nodes")
	}
	return lines
}

func TestScenesMatchC(t *testing.T) {
	want := strings.TrimSpace(runC(t, "testdata/scenes.c"))
	var got []string
	al, ah := animLines(t)
	el, eh := emojiLines(t)
	il, ih := imeLines(t)
	got = append(got, al...)
	got = append(got, el...)
	got = append(got, il...)
	got = append(got, a11yLines(t)...)
	if g := strings.Join(got, "\n"); g != want {
		t.Errorf("go:\n%s\nC:\n%s", g, want)
	}
	if ah != [4]uint32{0x0226f2f5, 0x906d3ce5, 0xdb5848ad, 0x12bea72d} {
		t.Errorf("anim hashes %08x", ah)
	}
	if eh != [6]uint32{0x3bab3875, 0x85fe883a, 0xd3add387, 0xb79518dd, 0x49a72d7d, 0x469719e1} {
		t.Errorf("emoji rtl hashes %08x", eh)
	}
	if ih != [3]uint32{0x2e20540f, 0xceafcbad, 0xc9c1add5} {
		t.Errorf("ime hashes %08x", ih)
	}
}

func TestEase(t *testing.T) {
	if Ease(EaseLinear, 0.25) != 0.25 || Ease(EaseIn, 0.5) != 0.125 || Ease(EaseOut, 0.5) != 0.875 || Ease(EaseInOut, 0.75) != 0.9375 {
		t.Error("ease values")
	}
	if Ease(EaseBounce, 0.25) != 0.47265625 || Ease(EaseOut, -3) != 0 || Ease(EaseIn, 7) != 1 {
		t.Error("bounce or clamping")
	}
	peak := 0.0
	for i := 1; i < 100; i++ {
		if v := Ease(EaseSpring, float64(i)/100); v > peak {
			peak = v
		}
	}
	if peak < 1.1 || peak > 1.25 {
		t.Errorf("spring peak %v", peak)
	}
	g, _ := Open(Options{Backend: BackendTest, Mode: ModeWindow, Width: 100, Height: 100, TestScript: "end 100000"})
	defer g.Close()
	u, _ := g.NewUI()
	defer u.Free()
	calls := 0
	u.Config().SetEase(func(x float64) float64 {
		calls++
		return x
	})
	var v float64
	for f := 0; f < 30; f++ {
		u.Begin(0)
		target := 0.0
		if f > 1 {
			target = 50
		}
		v = u.Anim("x", target, 100, EaseCustom)
		u.End()
	}
	if calls == 0 || v != 50 {
		t.Errorf("custom ease calls %d v %v", calls, v)
	}
}

func TestIMEEvents(t *testing.T) {
	long := strings.Repeat("q", 80)
	g, err := Open(Options{Backend: BackendTest, Mode: ModeWindow, Width: 100, Height: 100,
		TestScript: "1 text " + long + ";2 preedit にほん 3;3 preedit abc;4 preedit;end 10"})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	if err := g.PushIME("hi"); err != nil {
		t.Fatal(err)
	}
	ev, ok, _ := g.Poll(0)
	if !ok || ev.Type != EventText || ev.Text != "hi" || g.EventText() != "hi" {
		t.Errorf("pushed text %+v", ev)
	}
	if g.TextInput(true, &Rect{10, 20, 1, 16}) != nil || g.TextInput(false, nil) != nil {
		t.Error("text input")
	}
	got, quit := 0, false
	for f := 0; f < 20 && !quit; f++ {
		for ev := range g.Events(0) {
			switch {
			case ev.Type == EventQuit:
				quit = true
			case ev.Type == EventText:
				got |= 1
				if len(ev.Msg) != 63 || ev.Text != long {
					t.Errorf("long text cut wrong: %d %d", len(ev.Msg), len(ev.Text))
				}
			case ev.Type == EventPreedit && ev.Text == "にほん":
				got |= 2
				if ev.ID != 3 || g.PreeditCursor() != 3 {
					t.Error("preedit caret")
				}
			case ev.Type == EventPreedit && ev.Text == "abc":
				got |= 4
			case ev.Type == EventPreedit && ev.Text == "":
				got |= 8
			}
		}
		if c := g.Begin(0); c != nil {
			g.End(c)
		}
	}
	if got != 15 {
		t.Errorf("got %b", got)
	}
}

func TestFallbackLifetime(t *testing.T) {
	d, f := font(t, "dejavu-sub.ttf"), font(t, "colr.ttf")
	d.Fallback(f)
	f.Free()
	if d.Width("goome! 😀", 20) <= d.Width("goome! ", 20) {
		t.Error("fallback gone after freeing the go side")
	}
	f.Free()
	d.Free()
	if f.f != nil || d.f != nil {
		t.Error("fonts not freed")
	}
	defer func() {
		if recover() == nil {
			t.Error("freed fallback still usable")
		}
	}()
	f.Width("a", 10)
}
