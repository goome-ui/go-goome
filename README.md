<p align="center">
  <img src="https://raw.githubusercontent.com/goome-ui/.github/main/profile/logo.png" alt="goome" width="150">
</p>

<h1 align="center">🐹 go-goome 🐹</h1>

<p align="center">
<b>Overlays on top of everything, normal desktop apps, or ui inside your game, straight from Go.</b>
</p>

<p align="center">
  <img src="https://goome.wtf/badges/go-123.svg" alt="Go 1.23+">
  <a href="https://pkg.go.dev/github.com/goome-ui/go-goome"><img src="https://goome.wtf/badges/go-get.svg" alt="go get"></a>
  <img src="https://goome.wtf/badges/cgo-one-file.svg" alt="cgo, one file">
  <img src="https://goome.wtf/badges/dependencies-none.svg" alt="no dependencies">
  <img src="https://img.shields.io/badge/Linux-X11%20%2B%20Wayland-1b1d22?style=flat-square&logo=linux&logoColor=white" alt="Linux">
  <img src="https://img.shields.io/badge/Windows-10%20%2F%2011-1b1d22?style=flat-square&logo=windows11&logoColor=white" alt="Windows">
  <img src="https://img.shields.io/badge/macOS-14%2B-1b1d22?style=flat-square&logo=apple&logoColor=white" alt="macOS">
  <img src="https://img.shields.io/badge/license-GPL--3.0-1b1d22?style=flat-square" alt="GPL-3.0">
</p>

<p align="center">
  <a href="https://goome.wtf/docs/">Docs</a> &nbsp;·&nbsp;
  <a href="https://goome.wtf/builder/">Builder</a> &nbsp;·&nbsp;
  <a href="https://github.com/goome-ui/examples">Examples</a> &nbsp;·&nbsp;
  <a href="https://github.com/goome-ui/c-goome">C</a> &nbsp;·&nbsp;
  <a href="https://github.com/goome-ui/csharp-goome">C#</a> &nbsp;·&nbsp;
  <a href="https://github.com/goome-ui/python-goome">Python</a> &nbsp;·&nbsp;
  <a href="https://github.com/goome-ui/java-goome">Java</a> &nbsp;·&nbsp;
  <a href="https://github.com/goome-ui/godot-goome">Godot</a> &nbsp;·&nbsp;
  <a href="https://github.com/goome-ui/zig-goome">Zig</a> &nbsp;·&nbsp;
  <a href="https://github.com/goome-ui/lua-goome">Lua</a>
</p>

<p align="center">
  This is <a href="https://github.com/goome-ui/c-goome">goome</a> for Go. Draw over the whole screen and clicks only stop where you drew something.<br>
  It does app windows and offscreen frames for game engines too. It's cgo on top of the one file C library, so <code>go get</code> is all you do, theres no library to install.<br>
  The <a href="https://goome.wtf/builder/">builder</a> lets you drag a ui together in your browser.<br> <br>
  Questions go in <a href="https://github.com/goome-ui/c-goome/discussions">Discussions</a>.<br>
  Do you like this project? If so, please let me know!
  <a href="https://sevvy-wevvy.com/discord">Discord</a>
</p>

<h3 align="center">🫥 An overlay on top of everything 🫥</h3>

```go
g, err := goome.Open(goome.Options{Title: "my overlay"})
if err != nil {
	log.Fatal(err)
}
defer g.Close()
for running := true; running; {
	for ev := range g.Events(16) {
		if ev.Type == goome.EventQuit {
			running = false
		}
	}
	g.Frame(0, func(c *goome.Canvas) {
		c.RoundRect(100, 100, 220, 60, 12, 0xE0202838)
		c.Text(118, 120, "im on top of everything", false, 0xFFFFFFFF)
		c.Hit(100, 100, 220, 60)
	})
}
```

<h3 align="center">🪟 A normal app 🪟</h3>

```go
app := goome.NewApp("My app", 960, 640, 0)
name := ""
app.Run(func(ui *goome.UI) {
	ui.Input("Name", &name)
	if ui.Button("Say hi") {
		ui.Toast("hi "+name, 2000)
	}
})
```

<h3 align="center">🎮 Inside a game engine 🎮</h3>

```go
g, _ := goome.OpenOffscreen(800, 600, 1)
ui, _ := g.NewUI()
win, volume := goome.NewWindow(20, 20, 260, 0), 0.5
for {
	g.Push(goome.Event{Type: goome.EventMotion, X: mouseX, Y: mouseY})
	for ev := range g.Events(0) {
		ui.Feed(ev)
	}
	ui.Frame(0, func(*goome.Canvas) {
		if ui.Window("Settings", win, goome.WinMove) {
			ui.Slider("Volume", &volume, 0, 1)
		}
		ui.End()
	})
	frame, _ := g.LastFrame()
	uploadTexture(frame.Pixels(), frame.W(), frame.H())
}
```

<p align="center">
Overlay is still the default, nothing changes unless you ask for an app or offscreen.<br>
Go cant really get inside a game for mods, the other bindings do that part. <a href="https://github.com/goome-ui/c-goome">c-goome</a> has the DirectX 9/11/12 and Vulkan backends plus Metal and GL, <a href="https://github.com/goome-ui/csharp-goome">csharp-goome</a> has Unity mod templates, <a href="https://github.com/goome-ui/java-goome">java-goome</a> does Fabric and NeoForge, and <a href="https://github.com/goome-ui/godot-goome">godot-goome</a> is a GDExtension.
</p>

<details>
<summary align="center"><b>⬇️ Install ⬇️</b></summary>

<br>

```sh
go get github.com/goome-ui/go-goome@latest
```

<p align="center">
Thats the newest release, <code>@v0.5.5</code> pins one and <code>@dev</code> gets the newest code. Go 1.23 or newer, with cgo on. The first build compiles goome with your C compiler, that takes a bit, after that Go caches it.<br> <br>
On Debian/Ubuntu thats <code>sudo apt install build-essential</code>, on Mac <code>xcode-select --install</code>, on Windows <a href="https://www.mingw-w64.org/">mingw-w64</a> gcc on your PATH. If cgo is off (<code>CGO_ENABLED=0</code>, or no compiler on Windows) turn it on with <code>CGO_ENABLED=1</code>.<br> <br>
goome loads everything it needs at runtime, so theres no <code>-l</code> flags, no pkg-config and nothing to ship next to your exe.
</p>

```sh
git clone https://github.com/goome-ui/go-goome
cd go-goome
go test ./...
go run ./examples/ui
```
</details>

<details>
<summary align="center"><b>🚀 Quick Start 🚀</b></summary>

<br>

<p align="center">
A whole overlay that draws a box and prints clicks on it:
</p>

```go
package main

import (
	"fmt"
	"log"

	"github.com/goome-ui/go-goome"
)

func main() {
	g, err := goome.Open(goome.Options{Title: "my overlay"})
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()
	pace := goome.NewPacer(60)
	defer pace.Free()
	for running := true; running; {
		for ev := range g.Events(pace.Wait()) {
			switch ev.Type {
			case goome.EventQuit:
				running = false
			case goome.EventDown:
				fmt.Printf("clicked at %.0f,%.0f\n", ev.X, ev.Y)
			}
		}
		c := g.Begin(0)
		if c == nil {
			continue
		}
		c.RoundRect(100, 100, 200, 80, 12, 0xE0202838)
		c.Text(116, 128, "click me", true, 0xFFFFFFFF)
		c.Hit(100, 100, 200, 80)
		g.End(c)
	}
}
```

<p align="center">
Colors are plain <code>0xAARRGGBB</code> uint32s and every position is in logical desktop coords, same as in C.<br>
<code>c.Hit</code> collects the parts that take clicks and <code>g.End(c)</code> sends them. You can also pass rects straight to <code>g.End(c, goome.Rect{...})</code>.
</p>
</details>

<details>
<summary align="center"><b>🎛️ UI Kit 🎛️</b></summary>

<br>

<p align="center">
The immediate mode ui works like in C. Widgets take a pointer to your value and give back true when they changed it this frame.
</p>

```go
ui, _ := g.NewUI()
defer ui.Free()
win := goome.NewWindow(40, 40, 300, 0)
volume, muted, name := 0.5, false, "goome"
for running {
	for ev := range g.Events(16) {
		ui.Feed(ev)
	}
	if ui.Begin(0) != nil {
		if ui.Window("Settings", win, goome.WinMove|goome.WinCollapse) {
			if ui.Button("Say hi") {
				fmt.Println("hi")
			}
			ui.SliderFmt("Volume", &volume, 0, 1, "%.2f")
			ui.Toggle("Muted", &muted)
			ui.Input("Name", &name)
		}
		ui.End()
	}
	ui.End()
}
```

<p align="center">
Everything with a begin in C (<code>Begin</code>, <code>Window</code>, <code>Scroll</code>, <code>ID</code>, <code>App</code>, <code>Menubar</code>, <code>Menu</code>, <code>Context</code>, <code>Popup</code>, <code>Modal</code>, <code>Tree</code>, <code>Table</code>, <code>Split</code>) gets closed by <code>ui.End()</code>. Call it once after each one, no matter what it returned. It knows if C wants the end or not, so you dont have to remember which ones only end when they're open.<br>
That also means <code>defer ui.End()</code> works when a part of your ui lives in its own func:
</p>

```go
func settings(ui *goome.UI) {
	defer ui.End()
	if !ui.Window("Settings", win, goome.WinMove) {
		return
	}
	ui.Checkbox("Vsync", &vsync)
}
```

<p align="center">
<code>ui.Frame(0, func(c *goome.Canvas) { ... })</code> does the begin and end of a whole frame for you.<br>
<code>goome.Window</code> keeps a windows spot and size between frames, <code>Closed</code> hides it and <code>WinClose</code> gives it an x that does the same. <code>goome.Split</code> keeps where a split bar is.<br>
<code>ui.Changed</code> is the last widgets changed flag, <code>ui.Entered</code> is true on the frame enter got pressed in an <code>InputEx</code> with <code>InputEnterReturns</code>.<br>
<code>ui.Theme()</code> points right at goomes theme, so <code>ui.Theme().Accent = 0xFF6FDC5F</code> just works, plus <code>Dark()</code>, <code>Light()</code>, <code>Make()</code>, <code>Preset()</code>, <code>SetFont()</code>, <code>SetSkin()</code> and <code>SetIcon()</code>.
</p>

```go
if ui.Table("files", 3, 200, goome.TableSortable|goome.TableStripes) {
	ui.TableColumn("Name", 0)
	ui.TableColumn("Size", 80)
	ui.TableColumn("Kind", 0)
	ui.TableSort(&col, &desc)
	for _, f := range files {
		ui.TableRow()
		ui.Label(f.Name)
		ui.TableNext()
		ui.Label(f.Size)
		ui.TableNext()
		ui.Label(f.Kind)
	}
}
ui.End()

if ui.Split("main", side) {
	ui.Label("left")
	ui.SplitNext()
	ui.TextEdit("body", &text, 0, goome.EditLineNumbers|goome.EditWrap)
}
ui.End()

ui.PlotLines("fps", fpsHistory, 0, 144, 60)
ui.DragFloat("Speed", &speed, 0.1, 0, 10, "")
ui.Number("Price", &price, 0.01, "%.2f")
ui.PushColor(goome.StyleAccent, 0xFFF59E0B)
ui.Button("Orange")
ui.Pop(1)
```
</details>

<details>
<summary align="center"><b>🪟 App Mode 🪟</b></summary>

<br>

<p align="center">
<code>goome.NewApp</code> is the easy way. <code>Run</code> opens the window, runs the loop, only redraws when something happened (or while the ui is animating, or with <code>app.Animate = true</code>) and closes when someone hits the x.<br>
<code>app.OnEvent</code> gets every event before the ui does (return true to eat it), <code>app.OnClose</code> returns false to keep it open, <code>app.OnStart</code> runs once the window is up. <code>app.Start()</code> / <code>app.Step()</code> / <code>app.Stop()</code> let you run the loop yourself.<br> <br>
Without the helper, <code>goome.OpenApp("title", 960, 640)</code> gives you a normal <code>*Goome</code>. The window is monitor 0, so <code>g.Begin(0)</code> and <code>ui.Begin(0)</code> work like always and resizing sends <code>EventMonitors</code>.
</p>

| | |
|:--|:--|
| `g.SetTitle("x")`, `g.SetSize(w, h)`, `g.SetPos(x, y)`, `g.Pos()` | move it around. `Pos` gives `ErrNotSupported` on Wayland |
| `g.SetState(goome.StateMaximized)`, `g.State()` | window state |
| `g.Drag(edges)` | start a native move (0) or resize (`Edge*` bits) from the click thats happening |
| `g.Show(true)` | show or hide it |
| `g.SetCursor(goome.CursorHand)` | the mouse cursor |
| `g.Clipboard()`, `g.SetClipboard(s)` | plain strings |
| `EventClose`, `EventFocus`, `EventDrop` | close button, focus (`ev.ID` 1 or 0), a dropped file with the full path in `ev.Path` |
</details>

<details>
<summary align="center"><b>✨ Animations ✨</b></summary>

<br>

<p align="center">
<code>ui.Anim(id, target, ms, ease)</code> gives you a number that moves toward <code>target</code> over <code>ms</code>. goome keeps the state per id, so just call it every frame with where you want it to end up. Change the target halfway and it turns around from where it is.<br>
<code>ui.AnimColor</code> does the same for colors and <code>ui.Animating()</code> stays true while any of them move, so <code>App</code> keeps drawing until they're done.
</p>

```go
x := ui.Anim("panel", target, 250, goome.EaseOut)
col := ui.AnimColor("bg", 0xFF20D060, 200, goome.EaseLinear)
c := ui.Canvas()
c.RRect(x, 20, 120, 40, goome.Radius(8), col)

goome.Ease(goome.EaseBounce, 0.25)
```

<p align="center">
The eases are <code>EaseLinear</code>, <code>EaseIn</code>, <code>EaseOut</code>, <code>EaseInOut</code>, <code>EaseSmooth</code>, <code>EaseBack</code>, <code>EaseSpring</code> and <code>EaseBounce</code> (spring and back go a bit past the end and come back).<br>
<code>EaseCustom</code> uses your own func from <code>ui.Config().SetEase(func(t float64) float64 { ... })</code>.
</p>
</details>

<details>
<summary align="center"><b>🌍 Emoji And RTL 🌍</b></summary>

<br>

<p align="center">
Fonts do color emoji (COLR, CBDT and sbix), skin tones and ZWJ ones, plus Hebrew and Arabic with the right letter shapes. <code>goome.SystemFont</code> adds the systems emoji font by itself. For your own fonts, give them a fallback:
</p>

```go
main, _ := goome.LoadFont("myfont.ttf")
emoji, _ := goome.LoadFont("NotoColorEmoji.ttf")
main.Fallback(emoji)
c.FontText(main, 20, 10, 10, "hi 😀 שלום", 0xFFFFFFFF)

goome.SetTextDir(goome.DirRTL)
```

<p align="center">
The font holds on to its fallbacks, so freeing <code>emoji</code> first is fine, it only really goes away once <code>main</code> is freed too.<br>
<code>SetTextDir</code> sets the direction for all text. <code>DirAuto</code> (the default) picks it from the first real letter of each line.
</p>
</details>

<details>
<summary align="center"><b>⌨️ IME And Screen Readers ⌨️</b></summary>

<br>

<p align="center">
Inputs and text editors turn the IME on by themselves while they have focus and draw the text you're still composing under the caret. If you do your own text box, <code>g.TextInput(true, &goome.Rect{...})</code> turns it on and puts the candidate window by the caret. You get <code>EventPreedit</code> while someone is composing (<code>ev.Text</code> is the full text, <code>ev.ID</code> the caret) and <code>EventText</code> when its done.<br>
In tests and engines push them yourself with <code>g.PushPreedit("かな", 3)</code> and <code>g.PushIME("かな")</code>.<br> <br>
The ui builds an accessibility tree every frame, screen readers get it on each OS. Its on by default in a window, for offscreen turn it on with <code>ui.A11y(true)</code> and read it with <code>ui.A11yTree()</code>:
</p>

```go
for _, n := range ui.A11yTree() {
	fmt.Println(n.Role, n.Label, n.Value, n.Parent, n.Is(goome.A11yChecked))
}
ui.A11yRole(goome.RoleButton, "Knob")
ui.A11yValue("11")
ui.Custom("knob", 40, 40, 0)
g.A11yAction("press", node, "")
```

<p align="center">
<code>Parent</code> is the index of the parent node, -1 for the top. <code>A11yRole</code> and <code>A11yValue</code> describe the next <code>Custom</code> item, and <code>A11yAction</code> sends press, set, expand or focus like a screen reader would.
</p>
</details>

<details>
<summary align="center"><b>🎮 Engines And Offscreen 🎮</b></summary>

<br>

<p align="center">
<code>goome.OpenOffscreen(w, h, scale)</code> never opens a window. You push your engines input in, goome draws into a buffer it owns, and you copy that onto a texture. It never blocks and never starts a thread.
</p>

| | |
|:--|:--|
| `g.Push(goome.Event{...})` | feed it input, any mode. `g.PushClick(x, y)` and `g.PushText("hi")` are shortcuts |
| `g.LastFrame()` | the frame from the last end, `Pixels()` is a `[]uint32` right over goomes buffer. Good until the next begin |
| `g.FrameDirty()` | the `Rect` (real px) that changed since you last asked |
| `g.FrameCopy(dst, stride, goome.FmtRGBA, &rect)` | copies part of the frame into your `[]byte`, the non premul formats unpremultiply |
| `g.Resize(w, h, scale)` | new size |
| `g.Hit(x, y)` | true when that spot is on the ui, only send clicks to goome then |
| `g.Cursor()` | what cursor the ui wants |
| `ui.WantsMouse()`, `ui.WantsKeyboard()` | true while the ui is using them, so your game shouldnt |

<p align="center">
<code>examples/offscreen</code> renders a scripted ui with no window and saves it as a PNG with the standard library.
</p>
</details>

<details>
<summary align="center"><b>📚 The API 📚</b></summary>

<br>

<p align="center">
Everything from <a href="https://github.com/goome-ui/c-goome/blob/main/include/goome.h">goome.h</a> is in here, just in Go style. <a href="https://goome.wtf/docs/">The docs</a> explain every part of it. The names are the C ones in CamelCase without <code>goome_</code>: <code>goome_draw_rrect</code> is <code>c.RRect</code>, <code>GOOME_WIN_MOVE</code> is <code>goome.WinMove</code>, <code>GOOME_K_LEFT</code> is <code>goome.KeyLeft</code>, <code>GOOME_S_ACCENT</code> is <code>goome.StyleAccent</code>, event types are <code>goome.EventDown</code> and so on.<br> <br>
Errors come back as values. They're a <code>*goome.Error</code> with the C code on <code>.Code</code>, and <code>errors.Is(err, goome.ErrAgain)</code> works for <code>ErrNotSupported</code>, <code>ErrDenied</code>, <code>ErrNoDisplay</code>, <code>ErrGone</code>, <code>ErrAgain</code>, <code>ErrNoMem</code> and <code>ErrArg</code>.<br>
Drawing on a canvas thats already been sent, or a bitmap or font you freed, panics. Thats a bug in the program, like using a closed channel.
</p>

```go
g.Monitors()                              // areas, scale, real px
g.Hotkey("alt+h", 1)                      // comes back as an EventHotkey
img, err := g.Capture(0, &goome.Rect{0, 0, 200, 100})
img.Pixels()                              // []uint32 over goomes buffer, nil once you grab again
img.Copy()                                // your own copy
bmp, err := goome.LoadBitmap("logo.png")  // Free() it when youre done
font, err := goome.SystemFont("")         // font.Measure(text, size)
```

<p align="center">
<code>c.Pixels()</code>, <code>img.Pixels()</code> and <code>bmp.Pixels()</code> are slices straight over the C memory, nothing gets copied. Write into a canvas and call <code>c.Damage(rect)</code>.<br>
Things that hold C memory (<code>Bitmap</code>, <code>Font</code>, <code>Canvas</code> from <code>NewCanvas</code>, <code>Pacer</code>, <code>UI</code>) have a <code>Free()</code>, freeing twice is fine. <code>g.Close()</code> frees the uis it made too.
</p>
</details>

<details>
<summary align="center"><b>🧩 Examples 🧩</b></summary>

<br>

- `examples/hello` - a box you can click and type into, alt+h hides it, esc or the x quits
- `examples/grab` - `go run ./examples/grab [monitor] [x y w h] out.ppm` saves a screen grab
- `examples/ui` - the ui kit demo with every widget, a theme editor and a stats window
- `examples/app` - a small notes app in a normal window with menus, a split pane, a text editor and a modal
- `examples/offscreen` - draws the ui with no window and saves it as a PNG
- `examples/anim` - every ease racing next to each other, plus a color fade
- `examples/text` - emoji, Hebrew and Arabic with a direction picker, `go run ./examples/text extra.ttf` adds a fallback font
- `examples/a11y` - an offscreen ui that gets pressed and typed into through the accessibility tree and the IME, then prints the tree

They're ports of the C examples. Way more small projects live in [goome-ui/examples](https://github.com/goome-ui/examples).

```sh
go run ./examples/ui
```
</details>

<details>
<summary align="center"><b>🖥️ Where It Works 🖥️</b></summary>

<br>

Linux (X11 and Wayland), Windows 10 / 11 and macOS 14+, same as goome. Not every desktop can do everything, the <a href="https://github.com/goome-ui/c-goome#readme">c-goome README</a> has the table and <code>g.Can(goome.CapHotkeys)</code> tells you at runtime.

On Mac goome has to run on the main thread. Call <code>runtime.LockOSThread()</code> in an <code>init()</code> of your main package and open goome from <code>main</code>.
</details>

<details>
<summary align="center"><b>🧪 Testing Without A Screen 🧪</b></summary>

<br>

<p align="center">
goome has a <b>test</b> backend that fakes a screen and plays a script of what the "user" does. Thats what all the tests here run on, no display needed:
</p>

```sh
go test ./...
```

<p align="center">
You can test your own overlay the same way, with env vars and no code changes:
</p>

```sh
GOOME_BACKEND=test GOOME_TEST_SCRIPT="1 click 400,300; 2 key hi\n; end 10" go run .
```

<p align="center">
or from Go with <code>goome.Options{Backend: goome.BackendTest, TestScript: "screen 800x600 1; 1 click 120,90; end 5"}</code>.<br>
The script format is in the <a href="https://github.com/goome-ui/c-goome#readme">c-goome README</a>. The drawing and ui tests check the exact same pixel hashes as the tests in c-goome, on every OS, and the layout test builds <code>testdata/layout.c</code> and checks every mirrored struct against what C says.
</p>
</details>

<details>
<summary align="center"><b>🔧 Updating goome 🔧</b></summary>

<br>

```sh
tools/update_native.sh ../c-goome      # or a path to any goome.h, or nothing to clone main
go test ./...
```

<p align="center">
<code>native/VERSION</code> keeps the commit it came from and <code>consts.go</code> gets made again from the new header. If the layout test fails after that, a struct changed and the Go one needs the same change.
</p>
</details>

<details>
<summary align="center"><b>🛠 Pull Requests 🛠</b></summary>

<br>

<p align="center">
Want to help with go-goome? Make A pull request to <b>dev</b>!<br>
Small fixes are the easiest to merge, for anything bigger open an issue first so we can talk about it. Bugs in the C side go to <a href="https://github.com/goome-ui/c-goome">c-goome</a>.<br> <br>
As long as it works and doesn't add a dependency its very likely to get accepted!<br>
Read <a href="CONTRIBUTING.md">CONTRIBUTING.md</a> before you send one.
</p>
</details>

<details>
<summary align="center"><b>⚖️ License ⚖️</b></summary>

<br>

<p align="center">
go-goome is under the <a href="LICENSE">GPL-3.0</a>.<br> <br>
<sub>
• You can use it, change it and share it<br>
• If you share something made with it, you have to share your source too, under the same license<br>
• Keep the credits
</sub><br> <br>
In short: play fair, keep it open, and give credit where it's due.
</p>
</details>

<details>
<summary align="center"><b>🩷 Donate 🩷</b></summary>

<br>

<p align="center">
<a href="https://sevvy-wevvy.com/kofi">Ko-fi</a> (Preferred)<br>
<a href="https://sevvy-wevvy.com/patreon">Patreon</a><br>
<a href="https://sevvy-wevvy.com/throne">Throne</a>
</p>
</details>
