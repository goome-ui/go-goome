package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/goome-ui/go-goome"
)

type step struct {
	kind int
	x, y float64
	text string
}

var script = map[int][]step{
	2: {{goome.EventDown, 40, 109, ""}, {goome.EventUp, 40, 109, ""}},
	3: {{goome.EventDown, 60, 163, ""}},
	4: {{goome.EventMotion, 250, 163, ""}},
	5: {{goome.EventUp, 250, 163, ""}},
	6: {{goome.EventDown, 150, 223, ""}, {goome.EventUp, 150, 223, ""}},
	7: {{kind: -1, text: " from an engine"}},
	8: {{goome.EventMotion, 80, 279, ""}},
}

var cursorNames = []string{"arrow", "ibeam", "hand", "move", "resize_ew", "resize_ns", "resize_nwse", "resize_nesw", "wait", "no", "hidden"}

func writePNG(path string, f *goome.Image) error {
	w, h := f.W(), f.H()
	out := image.NewRGBA(image.Rect(0, 0, w, h))
	px := f.Pixels()
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			s := px[y*f.Stride()+x]
			bg := uint32(0xFF2E3C4A)
			if (x/16+y/16)&1 != 0 {
				bg = 0xFF3A4A5A
			}
			a := s >> 24
			mix := func(k uint) uint8 { return uint8((s>>k&255 + (bg>>k&255)*(255-a)/255) & 255) }
			out.SetRGBA(x, y, color.RGBA{mix(16), mix(8), mix(0), 255})
		}
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, out)
}

func main() {
	path := "offscreen.png"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	g, err := goome.OpenOffscreen(360, 330, 2)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer g.Close()
	ui, _ := g.NewUI()
	defer ui.Free()
	ui.Theme().AnimMs = 0
	win := goome.NewWindow(20, 20, 320, 0)
	check, volume, name := false, 0.25, "goome"
	for frame := 0; frame < 10; frame++ {
		for _, s := range script[frame] {
			if s.kind < 0 {
				g.PushText(s.text)
			} else {
				g.Push(goome.Event{Type: s.kind, X: s.x, Y: s.y, Button: 1})
			}
		}
		for ev := range g.Events(0) {
			ui.Feed(ev)
		}
		ui.Frame(0, func(*goome.Canvas) {
			if ui.Window("offscreen", win, goome.WinMove) {
				ui.Label("no window, the engine owns the screen")
				ui.Checkbox("Pushed a click", &check)
				ui.SliderFmt("Volume", &volume, 0, 1, "%.2f")
				ui.InputEx("Name", &name, 0, "", 64)
				ui.Progress(volume, "")
				ui.Button("Hover me")
			}
			ui.End()
		})
	}
	f, err := g.LastFrame()
	if err == nil {
		err = writePNG(path, f)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	b := 0
	if check {
		b = 1
	}
	fmt.Printf("check %d volume %.2f name '%s' cursor %s -> %s\n", b, volume, name, cursorNames[g.Cursor()], path)
}
