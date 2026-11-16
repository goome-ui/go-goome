package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/goome-ui/go-goome"
)

var roles = map[int]string{
	goome.RoleWindow: "window", goome.RoleButton: "button", goome.RoleCheckbox: "checkbox", goome.RoleRadio: "radio",
	goome.RoleSlider: "slider", goome.RoleText: "text", goome.RoleInput: "input", goome.RoleCombo: "combo",
	goome.RoleList: "list", goome.RoleItem: "item", goome.RoleTab: "tab", goome.RoleGroup: "group",
	goome.RoleTree: "tree", goome.RoleProgress: "progress", goome.RoleDialog: "dialog",
}

func main() {
	g, err := goome.OpenOffscreen(360, 300, 1)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer g.Close()
	ui, _ := g.NewUI()
	defer ui.Free()
	ui.A11y(true)
	win := goome.NewWindow(10, 10, 340, 0)
	sound, volume, name := false, 0.5, ""
	var tree []goome.A11yNode
	for frame := 0; frame < 8; frame++ {
		switch frame {
		case 2:
			for _, n := range tree {
				if n.Label == "Sound" {
					g.A11yAction("press", n, "")
				}
				if n.Label == "Name" {
					g.A11yAction("focus", n, "")
				}
			}
		case 4:
			g.PushPreedit("こんに", 9)
		case 5:
			g.PushPreedit("", 0)
			g.PushIME("こんにちは")
		}
		for ev := range g.Events(0) {
			ui.Feed(ev)
		}
		ui.Frame(0, func(*goome.Canvas) {
			if ui.Window("Settings", win, goome.WinMove) {
				ui.Checkbox("Sound", &sound)
				ui.SliderFmt("Volume", &volume, 0, 1, "%.2f")
				ui.Input("Name", &name)
				ui.Button("Save")
			}
			ui.End()
		})
		tree = ui.A11yTree()
	}
	for _, n := range tree {
		depth := 0
		for p := n.Parent; p >= 0; p = tree[p].Parent {
			depth++
		}
		line := fmt.Sprintf("%s%s %q", strings.Repeat("  ", depth), roles[n.Role], n.Label)
		if n.Value != "" {
			line += fmt.Sprintf(" = %q", n.Value)
		}
		if n.Is(goome.A11yChecked) {
			line += " checked"
		}
		if n.Is(goome.A11yFocused) {
			line += " focused"
		}
		fmt.Println(line)
	}
	fmt.Printf("sound %v name %q\n", sound, name)
}
