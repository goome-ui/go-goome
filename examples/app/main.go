package main

import (
	"fmt"
	"os"

	"github.com/goome-ui/go-goome"
)

func main() {
	app := goome.NewApp("Notes", 960, 640, goome.AppMenubar|goome.AppStatusbar)
	notes := []string{"groceries", "ideas", "todo"}
	bodies := map[int]string{0: "milk\neggs", 1: "an overlay that shows the weather", 2: "finish the go bindings"}
	pick, about, saved := 0, false, 0
	side := &goome.Split{Frac: 0.28}
	app.OnClose = func(*goome.App) bool { return true }
	err := app.Run(func(ui *goome.UI) {
		if ui.Menubar() {
			if ui.Menu("File") {
				if ui.MenuItem("Save", "Ctrl+S", nil) {
					saved++
					ui.Toast("saved", 1500)
				}
				ui.MenuSeparator()
				if ui.MenuItem("Quit", "Ctrl+Q", nil) {
					app.Quit()
				}
			}
			ui.End()
			if ui.Menu("Help") {
				if ui.MenuItem("About", "", nil) {
					about = true
				}
			}
			ui.End()
		}
		ui.End()
		if ui.Split("main", side) {
			for i, n := range notes {
				if ui.Selectable(n, i == pick) {
					pick = i
				}
			}
			if ui.Button("New note") {
				notes = append(notes, fmt.Sprintf("note %d", len(notes)+1))
				pick = len(notes) - 1
			}
			ui.SplitNext()
			body := bodies[pick]
			if ui.TextEdit("body", &body, 0, goome.EditLineNumbers|goome.EditWrap) {
				bodies[pick] = body
			}
		}
		ui.End()
		if ui.Modal("About", &about) {
			ui.Label("a tiny notes app made with goome")
			if ui.Button("Close") {
				about = false
			}
		}
		ui.End()
		ui.Statusbar(fmt.Sprintf("%d notes, saved %d times", len(notes), saved))
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "cant open a window:", err)
		os.Exit(1)
	}
}
