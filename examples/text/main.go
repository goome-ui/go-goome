package main

import (
	"fmt"
	"os"

	"github.com/goome-ui/go-goome"
)

var lines = []string{
	"plain text with emoji 😀 👍🏽 ❤️ 🐱",
	"עברית: שלום עולם (123)",
	"عربي: مرحبا بالعالم",
	"mixed: goome ממש עובד 2026",
}

func main() {
	app := goome.NewApp("goome text", 640, 420, goome.AppStatusbar)
	var font *goome.Font
	dir, note := 0, "type here, the ime works too"
	app.OnStart = func(a *goome.App) {
		f, err := goome.SystemFont("")
		if err != nil {
			fmt.Println("no system font, using the built in one:", err)
			return
		}
		font = f
		a.UI.Theme().SetFont(goome.FontMain, f)
		a.UI.Theme().FontSize = 18
		if len(os.Args) > 1 {
			if extra, err := goome.LoadFont(os.Args[1]); err == nil {
				f.Fallback(extra)
			}
		}
	}
	err := app.Run(func(ui *goome.UI) {
		if ui.Combo("Direction", &dir, []string{"auto", "left to right", "right to left"}) {
			goome.SetTextDir([]int{goome.DirAuto, goome.DirLTR, goome.DirRTL}[dir])
		}
		ui.Separator()
		for _, l := range lines {
			ui.Label(l)
		}
		ui.Separator()
		ui.Input("Note", &note)
		ui.Statusbar(fmt.Sprintf("%d bytes typed", len(note)))
	})
	if font != nil {
		font.Free()
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "cant open a window:", err)
		os.Exit(1)
	}
}
