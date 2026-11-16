package main

import (
	"fmt"
	"os"

	"github.com/goome-ui/go-goome"
)

func main() {
	g, err := goome.Open(goome.Options{Title: "goome hello"})
	if err != nil {
		fmt.Fprintln(os.Stderr, "cant open an overlay:", err)
		os.Exit(1)
	}
	defer g.Close()
	mons := g.Monitors()
	fmt.Println("backend:", g.BackendName())
	for i, m := range mons {
		a := m.Area
		fmt.Printf("monitor %d: %s %dx%d at %d,%d scale %.2f\n", i, m.Name, a.W, a.H, a.X, a.Y, m.Scale)
	}
	if err := g.Hotkey("alt+h", 1); err != nil {
		fmt.Println("no global hotkeys here:", g.LastError())
	}

	bx, by := 80.0, 80.0
	if len(mons) > 0 {
		bx, by = float64(mons[0].Area.X+80), float64(mons[0].Area.Y+80)
	}
	clicks, hidden, typing, running := 0, false, false, true
	typed := ""
	pace := goome.NewPacer(30)
	defer pace.Free()
	for running {
		for ev := range g.Events(pace.Wait()) {
			switch ev.Type {
			case goome.EventQuit:
				running = false
			case goome.EventDown:
				fmt.Printf("click %.0f,%.0f\n", ev.X, ev.Y)
				switch {
				case ev.X > bx+250 && ev.Y < by+30:
					running = false
				case ev.Y > by+80:
					typing = !typing
					g.Keyboard(typing)
				default:
					clicks++
				}
			case goome.EventKey:
				fmt.Printf("key %d mods %x text '%s'\n", ev.Key, ev.Mods, ev.Text)
				switch {
				case ev.Key == goome.KeyEsc:
					running = false
				case ev.Key == goome.KeyBackspace && typed != "":
					r := []rune(typed)
					typed = string(r[:len(r)-1])
				case ev.Text != "" && len(typed+ev.Text) < 40:
					typed += ev.Text
				}
			case goome.EventHotkey:
				hidden = !hidden
			case goome.EventControl:
				if ev.Msg == "hide" {
					hidden = !hidden
				}
			case goome.EventScroll:
				by += float64(ev.DY * 10)
			}
		}
		c := g.Begin(0)
		if c == nil {
			continue
		}
		if !hidden {
			c.RoundRect(bx, by, 280, 120, 14, 0xE0202838)
			border := uint32(0xFF8090A0)
			if clicks&1 != 0 {
				border = 0xFF6FDC5F
			}
			c.Frame(bx, by, 280, 120, 2, border)
			c.Text(bx+16, by+12, "hello from goome", true, 0xFFFFFFFF)
			c.Line(bx+256, by+10, bx+268, by+22, 2, 0xFFFF6B8B)
			c.Line(bx+268, by+10, bx+256, by+22, 2, 0xFFFF6B8B)
			c.Text(bx+16, by+46, fmt.Sprintf("clicks: %d", clicks), false, 0xFFC0C8D0)
			box, line := uint32(0xFF303848), "click here to type"
			if typing {
				box, line = 0xFF2E4A30, "type something..."
				if typed != "" {
					line = typed
				}
			}
			c.RoundRect(bx+12, by+80, 256, 28, 8, box)
			c.Text(bx+22, by+86, line, false, 0xFFFFFFFF)
			c.Hit(bx, by, 280, 120)
		}
		g.End(c)
	}
}
