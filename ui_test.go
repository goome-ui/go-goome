package goome

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"
)

const uiScript = "screen 400x360 1; 1 click 35,71; 2 click 29,106; 3 hold 40,157 6; 4 hover 120,157; " +
	"7 click 100,213; 8 key ab; 9 key \\n; 10 key zz; 11 click 120,273; 13 click 120,206; " +
	"15 click 300,20; 16 hold 100,20 20; 18 hover 150,40; end 22"

func TestUIScripted(t *testing.T) {
	g, err := Open(Options{Backend: BackendTest, TestScript: uiScript})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	u, err := g.NewUI()
	if err != nil {
		t.Fatal(err)
	}
	u.Theme().SetFont(FontMain, nil)
	u.Theme().FontSize = 14
	win := NewWindow(10, 10, 220, 0)
	val, check, pick, text := 0.0, false, 0, ""
	quit, frames, clicks, downs, keysOff, kbSeen := false, 0, 0, 0, 0, false
	hs := map[int]string{}
	for !quit && frames < 100 {
		for ev := range g.Events(0) {
			if ev.Type == EventQuit {
				quit = true
			}
			if ev.Type == EventDown {
				downs++
			}
			if ev.Type == EventKey && frames > 10 {
				keysOff++
			}
			u.Feed(ev)
		}
		c := u.Begin(0)
		if c == nil {
			t.Fatal("no canvas")
		}
		if u.Window("test", win, WinMove) {
			if u.Button("Go") {
				clicks++
			}
			u.Checkbox("Check", &check)
			u.Slider("Val", &val, 0, 1)
			u.Input("Text", &text)
			u.Combo("Pick", &pick, []string{"one", "two", "three"})
		}
		u.End()
		px, w, h, stride := c.Pixels(), c.W(), c.H(), c.Stride()
		u.End()
		if frames == 0 || frames == 12 || frames == 21 {
			hs[frames] = fmt.Sprintf("%08x", fnv(px, w, h, stride))
		}
		if c.Valid() {
			t.Fatal("canvas valid after the frame")
		}
		if frames == 8 {
			kbSeen = text != ""
		}
		frames++
	}
	if clicks != 1 || !check || val < 0.45 || val > 0.55 || !kbSeen || text != "ab" || keysOff != 0 || pick != 1 {
		t.Errorf("clicks %d check %v val %v kb %v text %q keysOff %d pick %d", clicks, check, val, kbSeen, text, keysOff, pick)
	}
	if win.X != 60 || win.Y != 30 || downs != 7 {
		t.Errorf("window at %v,%v downs %d", win.X, win.Y, downs)
	}
	want := map[int]string{0: "51713294", 12: "b52f16c5", 21: "ea01692c"}
	if fmt.Sprint(hs) != fmt.Sprint(want) {
		t.Errorf("hashes %v want %v", hs, want)
	}
	u.Free()
	u.Free()
}

func TestUIWidgets(t *testing.T) {
	g, err := Open(Options{Backend: BackendTest, TestScript: "screen 800x600 1; 2 click 60,80; end 6"})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	u, _ := g.NewUI()
	th := u.Theme()
	th.SetFont(FontMain, nil)
	th.Light()
	if th.Text>>24 == 0 {
		t.Error("light theme has no text color")
	}
	th.Dark()
	th.Accent = 0xFF6FDC5F
	th.Radius, th.Pad, th.AnimMs = 4, 8, 0
	if u.Theme().Accent != 0xFF6FDC5F || u.Theme().Pad != 8 {
		t.Error("theme edit didnt stick")
	}
	if err := th.Preset(PresetMidnight); err != nil {
		t.Error(err)
	}
	if err := th.Preset(99); err == nil {
		t.Error("preset 99 worked")
	}
	th.Make(0xFF4FA3FF, true, 10)
	win, stats := NewWindow(20, 20, 400, 0), NewWindow(500, 20, 200, 0)
	logo, _ := BitmapFromARGB(make([]uint32, 64), 8, 8)
	defer logo.Free()
	tab, vol, n, color, mode, on, more, name := 0, 0.5, 3, uint32(0xFF4FA3FF), 1, false, true, "hi"
	body, side := "line one\nline two", &Split{Frac: 0.3}
	about, tree, checked := true, true, false
	col, desc := 0, false
	for i := 0; i < 8; i++ {
		for ev := range g.Events(0) {
			u.Feed(ev)
		}
		if u.Begin(0) == nil {
			u.End()
			continue
		}
		if u.Window("main", win, WinMove|WinResize|WinCollapse) {
			u.Tabs("tabs", []string{"A", "B"}, &tab)
			u.ID(strconv.Itoa(tab))
			if u.Scroll("page", 0) {
				u.Heading("Head")
				u.Row(2)
				u.Button("One")
				u.Button("Two")
				u.Tooltip("tip")
				u.Hovered()
				u.Label("100% sure %s %d")
				u.Wrapped("some long text that wraps around a bit")
				u.Separator()
				u.Space(4)
				u.Width(120)
				u.Toggle("Toggle", &on)
				u.Radio("Small", &mode, 0)
				u.SameLine()
				u.Radio("Big", &mode, 1)
				u.SliderFmt("Vol", &vol, 0, 1, "%.2f")
				u.SliderInt("N", &n, 0, 10)
				u.Progress(vol, fmt.Sprintf("%d%%", int(vol*100)))
				u.Progress(-1, "")
				u.Color("Color", &color)
				u.Input("Name", &name)
				u.Password("Pass", &name, "8 or more")
				if u.Collapsing("More", &more) {
					for k := 0; k < 3; k++ {
						u.Selectable(fmt.Sprintf("Item %d", k), k == 1)
					}
				}
				u.Collapsing("Self kept", nil)
				u.Image(logo, 32, 32, 6)
				u.ImageButton("img", logo, 24, 24)
				if u.Tree("Tree", &tree) {
					u.Label("leaf")
				}
				u.End()
				if u.Table("files", 2, 80, TableSortable|TableStripes) {
					u.TableColumn("Name", 0)
					u.TableColumn("Size", 60)
					u.TableSort(&col, &desc)
					u.TableRow()
					u.Label("a")
					u.TableNext()
					u.Label("1")
				}
				u.End()
				if u.Split("main", side) {
					u.Label("left")
					u.SplitNext()
					u.TextEdit("body", &body, 60, EditLineNumbers)
				}
				u.End()
				u.DragFloat("Speed", &vol, 0.01, 0, 1, "")
				u.DragInt("Count", &n, 1, 0, 10)
				u.Number("Price", &vol, 0.01, "%.2f")
				u.PlotLines("fps", []float32{1, 3, 2, 5}, 0, 5, 40)
				u.PlotBars("bars", []float32{1, 3, 2, 5}, 0, 5, 40)
				u.Spinner(16)
				u.Link("docs")
				u.Badge("new", th.Accent)
				u.PushColor(StyleAccent, 0xFFF59E0B)
				u.PushNum(StyleRadius, 99)
				u.PushRadii(Radii{1, 2, 3, 4})
				u.Button("Styled")
				u.Pop(3)
				u.Custom("knob", 40, 40, ItemDragY)
				if cv := u.Canvas(); cv != nil {
					cv.Circle(10, 10, 5, 0xFFFFFFFF)
				}
				st := u.Store("knob", 2)
				st[0]++
				if first, last := u.ListClip(100, 20); last < first {
					t.Error("list clip")
				}
				u.Context("ctx")
				u.End()
			}
			u.End()
			u.End()
		}
		u.End()
		if u.Window("Stats", stats, WinMove|WinClose) {
			u.Label("stats")
		}
		u.End()
		if u.Modal("About", &about) {
			u.Label("made with goome")
		}
		u.End()
		if u.Menu("File") {
			u.MenuItem("Check", "Ctrl+K", &checked)
			u.MenuSeparator()
		}
		u.End()
		u.PopupOpen("pop")
		if u.Popup("pop") {
			u.Label("popped")
			u.PopupClose()
		}
		u.End()
		u.Toast("saved", 500)
		u.Hit(0, 0, 4, 4)
		u.Animating()
		if u.Depth() != 1 {
			t.Fatalf("stack depth %d", u.Depth())
		}
		u.End()
	}
	if stats.Closed {
		t.Error("stats window closed by itself")
	}
	stats.Closed = true
	u.Begin(0)
	if u.Window("Stats", stats, WinClose) {
		t.Error("closed window still shows")
	}
	u.End()
	u.End()
	cfg := u.Config()
	cfg.SetEase(func(x float64) float64 { return x * x })
	cfg.SetMask("*")
	if cfg.Mask() != "*" {
		t.Error("mask")
	}
	if err := u.Bind(ActCopy, "ctrl+c, ctrl+insert"); err != nil {
		t.Error(err)
	}
	u.Scale(1.25)
	u.Free()
}

func TestUIFonts(t *testing.T) {
	g, _ := Open(Options{Backend: BackendTest, TestScript: "end 3"})
	defer g.Close()
	u, _ := g.NewUI()
	f, err := FontFromBytes(data(t, "dejavu-sub.ttf"))
	if err != nil {
		t.Fatal(err)
	}
	u.Theme().SetFont(FontMain, f)
	if !u.Theme().HasFont(FontMain, f) {
		t.Error("font not set")
	}
	f.Free()
	if u.Theme().HasFont(FontMain, nil) == false {
		t.Error("freed font still in the theme")
	}
	b, _ := BitmapFromARGB(make([]uint32, 16), 4, 4)
	u.Theme().SetSkin(SkinButton, NewSkin(b, 1, 1, 1, 1))
	u.Theme().SetIcon(IconCheck, b)
	b.Free()
	if !u.Theme().Skin(SkinButton).Empty() || u.Theme().HasIcon(IconCheck) {
		t.Error("freed bitmap still in the theme")
	}
	u.Frame(0, func(c *Canvas) {
		if u.Window("w", NewWindow(0, 0, 100, 0), 0) {
			u.Label("hi")
		}
		u.End()
	})
}

func TestApp(t *testing.T) {
	a := NewApp("app", 700, 500, AppStatusbar)
	a.Options = Options{Backend: BackendTest, TestScript: "screen 1x1 1; 3 resize 700x500; 5 close; end 10"}
	closed := false
	a.OnClose = func(*App) bool {
		closed = true
		return true
	}
	err := a.Run(func(ui *UI) {
		ui.Label("hello")
		ui.Statusbar("ready")
	})
	if err != nil {
		t.Fatal(err)
	}
	if !closed || a.Frames < 2 {
		t.Errorf("closed %v frames %d", closed, a.Frames)
	}
}

func TestGoThemePush(t *testing.T) {
	g, err := Open(Options{Backend: BackendTest, TestScript: "screen 320x240 1; end 4"})
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	u, err := g.NewUI()
	if err != nil {
		t.Fatal(err)
	}
	u.Theme().SetFont(FontMain, nil)
	u.Theme().FontSize = 14
	for _, mk := range []func(*Theme){
		func(th *Theme) { th.Light() },
		func(th *Theme) { th.Dark() },
		func(th *Theme) { th.Make(0xFF4FA3FF, true, 8) },
	} {
		th := new(Theme)
		mk(th)
		if th.Text>>24 == 0 {
			t.Fatal("go theme has no text color")
		}
		th.SetFont(FontMain, nil)
		th.FontSize = 14
		if c := u.Begin(0); c == nil {
			t.Fatal("no canvas")
		}
		u.PushTheme(th)
		u.Label("themed")
		u.End()
	}
}

func TestUISloppyArguments(t *testing.T) {
	g, err := OpenOffscreen(200, 200, 1)
	if err != nil {
		t.Fatal(err)
	}
	defer g.Close()
	u, err := g.NewUI()
	if err != nil {
		t.Fatal(err)
	}
	u.Begin(0)
	func() {
		b, err := BitmapFromARGB(make([]uint32, 64*64), 64, 64)
		if err != nil {
			t.Fatal(err)
		}
		u.Image(b, 64, 64, 0)
	}()
	for i := 0; i < 4; i++ {
		runtime.GC()
		time.Sleep(10 * time.Millisecond)
	}
	v := 0.5
	u.SliderFmt("fmt", &v, 0, 1, "%s%s%n")
	u.Number("fmt2", &v, 1, "%d")
	u.DragFloat("fmt3", &v, 1, 0, 10, "%5s")
	if s := u.Store("big", 1<<32+1); s != nil {
		t.Errorf("store handed back %d numbers", len(s))
	}
	u.End()
	for in, want := range map[string]string{"%.2f": "%.2f", "%5.1f%%": "%5.1f%%", "%d": "", "%f%f": "", "%": "", "plain": "plain"} {
		if got := numFmt(in); got != want {
			t.Errorf("numFmt(%q) = %q", in, got)
		}
	}
}
