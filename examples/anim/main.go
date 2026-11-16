package main

import (
	"fmt"
	"math"
	"os"

	"github.com/goome-ui/go-goome"
)

var eases = []struct {
	name string
	kind int
}{
	{"linear", goome.EaseLinear},
	{"in", goome.EaseIn},
	{"out", goome.EaseOut},
	{"in out", goome.EaseInOut},
	{"smooth", goome.EaseSmooth},
	{"back", goome.EaseBack},
	{"spring", goome.EaseSpring},
	{"bounce", goome.EaseBounce},
	{"custom", goome.EaseCustom},
}

func main() {
	app := goome.NewApp("goome anim", 640, 520, goome.AppStatusbar)
	right, dark := false, true
	ms := 600.0
	app.OnStart = func(a *goome.App) {
		a.UI.Config().SetEase(func(t float64) float64 { return math.Sqrt(t) })
	}
	err := app.Run(func(ui *goome.UI) {
		if ui.Button("Go") {
			right = !right
		}
		ui.SameLine()
		ui.Toggle("Dark box", &dark)
		ui.SliderFmt("Time", &ms, 100, 2000, "%.0f ms")
		c := ui.Canvas()
		item, _ := ui.Custom("lanes", 600, float64(len(eases))*36, goome.ItemNoHit)
		for i, e := range eases {
			to := 0.0
			if right {
				to = 1
			}
			v := ui.Anim("lane "+e.name, to, ms, e.kind)
			y := item.Y + float64(i)*36
			if c != nil {
				c.Text(item.X, y+8, e.name, false, 0xFFB0B8C0)
				c.RRect(item.X+80, y+14, 500, 4, goome.Radius(2), 0x40FFFFFF)
				c.Circle(item.X+90+v*480, y+16, 10, 0xFF4FB3F7)
			}
		}
		col := uint32(0xFFF0F0F0)
		if dark {
			col = 0xFF202830
		}
		box := ui.AnimColor("box", col, 300, goome.EaseOut)
		if c != nil {
			c.RRect(item.X+520, item.Y-90, 60, 60, goome.Radius(12), box)
		}
		state := "idle"
		if ui.Animating() {
			state = "moving"
		}
		ui.Statusbar(fmt.Sprintf("%s, %d frames", state, app.Frames))
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "cant open a window:", err)
		os.Exit(1)
	}
}
