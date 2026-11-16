package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/goome-ui/go-goome"
)

var fruits = []string{"Mango", "Banana", "Cherry", "Kiwi", "Apple", "Lemon", "Grape", "Fig", "Dragonfruit", "Pear"}
var tabs = []string{"Widgets", "Images", "Theme", "About"}

func gradient() *goome.Bitmap {
	px := make([]uint32, 128*128)
	for y := 0; y < 128; y++ {
		for x := 0; x < 128; x++ {
			r := uint32(0x2A + x*(0x6F-0x2A)/127)
			g := uint32(0x40 + y*(0xDC-0x40)/127)
			b := uint32(0x9A - (x+y)*0x50/254)
			px[y*128+x] = 0xFF000000 | r<<16 | g<<8 | b
		}
	}
	bm, _ := goome.BitmapFromARGB(px, 128, 128)
	return bm
}

func plural(n int, none, one, many string) string {
	switch n {
	case 0:
		return none
	case 1:
		return one
	}
	return fmt.Sprintf(many, n)
}

type demo struct {
	g                                 *goome.Goome
	u                                 *goome.UI
	logo                              *goome.Bitmap
	main, stats                       *goome.Window
	tab, clicks, mode, count, fruit   int
	picked, light, imgClicks          int
	grid, wifi, notify, more          bool
	ownAccent, solid, always, running bool
	volume, radius, anim              float64
	highlight, accent                 uint32
	name, search                      string
	fps, frameMs                      float64
}

func (d *demo) setTheme() {
	t := d.u.Theme()
	if d.light == 1 {
		t.Light()
	} else {
		t.Dark()
	}
	if d.ownAccent {
		t.Accent = d.accent
	}
	t.Radius, t.AnimMs = d.radius, d.anim
	if d.solid {
		t.Panel |= 0xFF000000
	}
}

func (d *demo) widgets() {
	u := d.u
	u.Heading("Buttons")
	u.Row(3)
	if u.Button("Click me") {
		d.clicks++
	}
	if u.Button("Reset") {
		d.clicks = 0
	}
	u.Button("Hover me")
	u.Tooltip("hi :)")
	u.Label(plural(d.clicks, "not clicked yet", "clicked once", "clicked %d times"))
	u.Separator()
	u.Heading("Switches")
	u.Checkbox("Show grid", &d.grid)
	u.SameLine()
	u.Checkbox("Notifications", &d.notify)
	u.Toggle("Wi-Fi", &d.wifi)
	for i, name := range []string{"Small", "Medium", "Large"} {
		if i > 0 {
			u.SameLine()
		}
		u.Radio(name, &d.mode, i)
	}
	u.Separator()
	u.Heading("Values")
	u.SliderFmt("Volume", &d.volume, 0, 1, "%.2f")
	u.SliderInt("Count", &d.count, 0, 100)
	u.Progress(d.volume, fmt.Sprintf("%d%%", int(d.volume*100+0.5)))
	u.Progress(-1, "")
	u.Row(2)
	u.InputEx("Name", &d.name, 0, "", 64)
	u.InputEx("Search", &d.search, 0, "", 64)
	u.Combo("Fruit", &d.fruit, fruits)
	u.Color("Highlight", &d.highlight)
	if u.Collapsing("More stuff", &d.more) {
		u.Wrapped("a long list so theres something to scroll. pick one, it shows up under it.")
		if u.Scroll("list", 150) {
			for i := 0; i < 30; i++ {
				if u.Selectable(fmt.Sprintf("Item %d", i+1), d.picked == i) {
					d.picked = i
				}
			}
		}
		u.End()
		if d.picked >= 0 {
			u.Label(fmt.Sprintf("picked item %d", d.picked+1))
		} else {
			u.Label("nothing picked")
		}
	}
}

func (d *demo) images() {
	u := d.u
	u.Heading("Generated gradient")
	for i, r := range []float64{18, 60, 0} {
		if i > 0 {
			u.SameLine()
		}
		u.Image(d.logo, 120, 120, r)
	}
	u.Label("image buttons")
	for i := 0; i < 4; i++ {
		if i > 0 {
			u.SameLine()
		}
		if u.ImageButton(fmt.Sprintf("img%d", i), d.logo, 56, 56) {
			d.imgClicks++
		}
	}
	u.Label(plural(d.imgClicks, "no image clicks yet", "one image click", "%d image clicks"))
	u.Image(d.logo, 0, 90, 10)
}

func (d *demo) theme() {
	u := d.u
	u.Heading("Theme")
	was := d.light
	u.Radio("Dark", &d.light, 0)
	u.SameLine()
	u.Radio("Light", &d.light, 1)
	ch := d.light != was
	if u.Color("Accent", &d.accent) {
		d.ownAccent, ch = true, true
	}
	ch = u.SliderFmt("Corner radius", &d.radius, 0, 20, "%.0f px") || ch
	ch = u.SliderFmt("Animation", &d.anim, 0, 400, "%.0f ms") || ch
	ch = u.Toggle("Solid panels", &d.solid) || ch
	if u.Button("Reset theme") {
		d.light, d.ownAccent, d.radius, d.anim, d.solid, d.accent = 0, false, 10, 120, false, 0xFF6FDC5F
		ch = true
	}
	if ch {
		d.setTheme()
	}
	if !d.ownAccent {
		d.accent = u.Theme().Accent
	}
}

func (d *demo) about() {
	u := d.u
	u.Heading("About")
	u.Wrapped("made with goome. try clicking next to this window, your desktop still gets it.")
	u.Label(fmt.Sprintf("goome %s on %s", goome.Version, d.g.BackendName()))
	u.Label(fmt.Sprintf("%.0f fps, %.2f ms a frame", d.fps, d.frameMs))
	if d.stats.Closed && u.Button("Show stats") {
		d.stats.Closed = false
	}
}

func (d *demo) draw() {
	u := d.u
	flags := goome.WinMove | goome.WinResize | goome.WinCollapse
	if u.Window("goome ui", d.main, flags) {
		u.Tabs("tabs", tabs, &d.tab)
		u.ID(tabs[d.tab])
		if u.Scroll("page", 0) {
			[]func(){d.widgets, d.images, d.theme, d.about}[d.tab]()
		}
		u.End()
		u.End()
	}
	u.End()
	if u.Window("Stats", d.stats, flags|goome.WinClose) {
		u.Label(fmt.Sprintf("%.0f fps", d.fps))
		u.Label(fmt.Sprintf("%.2f ms a frame", d.frameMs))
		u.Progress(d.frameMs/16.7, fmt.Sprintf("%d%% of a 60 fps frame", int(d.frameMs*100/16.7+0.5)))
		u.Toggle("Redraw every frame", &d.always)
		if u.Button("Quit") {
			d.running = false
		}
	}
	u.End()
}

func (d *demo) run() {
	g, u := d.g, d.u
	test := strings.HasPrefix(g.BackendName(), "test")
	var spent int64
	tframes, accN, fpsN := 0, 0, 0
	accMs := 0.0
	fpsT := goome.NowMs()
	pace := goome.NewPacer(60)
	defer pace.Free()
	dirty := true
	for d.running {
		busy := dirty || d.always || test || u.Animating()
		wait := 500
		if busy {
			wait = pace.Wait()
		}
		for ev := range g.Events(wait) {
			if ev.Type == goome.EventQuit {
				d.running = false
			}
			u.Feed(ev)
			dirty = true
		}
		if !d.running {
			break
		}
		if !(dirty || d.always || test || u.Animating()) {
			continue
		}
		dirty = false
		t0 := goome.NowMs()
		if u.Begin(0) == nil {
			u.End()
			continue
		}
		d.draw()
		u.End()
		t1 := goome.NowMs()
		spent += t1 - t0
		tframes++
		accMs += float64(t1 - t0)
		accN++
		fpsN++
		if t1-fpsT >= 500 {
			d.fps = float64(fpsN) * 1000 / float64(t1-fpsT)
			d.frameMs = accMs / float64(accN)
			fpsT, fpsN, accMs, accN = t1, 0, 0, 0
		}
	}
	if tframes > 0 {
		fmt.Printf("%d frames, %.2f ms a frame\n", tframes, float64(spent)/float64(tframes))
	}
}

func main() {
	g, err := goome.Open(goome.Options{Title: "goome ui", NoCapture: true})
	if err != nil {
		fmt.Fprintln(os.Stderr, "cant open an overlay:", err)
		os.Exit(1)
	}
	defer g.Close()
	u, err := g.NewUI()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer u.Free()
	a := goome.Rect{W: 1280, H: 720}
	if m := g.Monitors(); len(m) > 0 {
		a = m[0].Area
	}
	d := &demo{
		g: g, u: u, logo: gradient(),
		main:  goome.NewWindow(float64(a.X+40), float64(a.Y+40), 440, 620),
		stats: goome.NewWindow(float64(a.X+500), float64(a.Y+40), 280, 0),
		grid:  true, wifi: true, mode: 1, count: 42, fruit: 2, picked: -1, more: true,
		volume: 0.65, radius: 10, anim: 120, highlight: 0xFF4FA3FF, accent: 0xFF6FDC5F,
		name: "goome", running: true,
	}
	defer d.logo.Free()
	d.run()
}
