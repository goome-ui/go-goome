package goome

import "strings"

type App struct {
	Title   string
	Width   int
	Height  int
	Flags   int
	FPS     int
	IdleMs  int
	Light   bool
	Animate bool
	Options Options

	OnStart func(a *App)
	OnEvent func(a *App, ev Event) bool
	OnClose func(a *App) bool

	G      *Goome
	UI     *UI
	Frames int

	running bool
	dirty   int
	every   bool
	pace    *Pacer
}

func NewApp(title string, w, h, flags int) *App {
	return &App{Title: title, Width: w, Height: h, Flags: flags, FPS: 60, IdleMs: 500}
}

func (a *App) Redraw() {
	if a.dirty < 1 {
		a.dirty = 1
	}
}

func (a *App) Quit() { a.running = false }

func (a *App) Running() bool { return a.running }

func (a *App) Start() error {
	o := a.Options
	o.Title, o.Mode = a.Title, ModeWindow
	if o.Width == 0 {
		o.Width, o.Height = a.Width, a.Height
	}
	if o.WindowFlags == 0 {
		o.WindowFlags = WindowResizable | WindowCenter
	}
	g, err := Open(o)
	if err != nil {
		return err
	}
	ui, err := g.NewUI()
	if err != nil {
		g.Close()
		return err
	}
	a.G, a.UI = g, ui
	if a.Light {
		ui.Theme().Light()
	}
	a.every = strings.HasPrefix(g.BackendName(), "test")
	if a.FPS <= 0 {
		a.FPS = 60
	}
	a.pace = NewPacer(a.FPS)
	a.dirty = 2
	a.running = true
	if a.OnStart != nil {
		a.OnStart(a)
	}
	return nil
}

func (a *App) handle(ev Event) {
	if a.OnEvent != nil && a.OnEvent(a, ev) {
		return
	}
	a.UI.Feed(ev)
	switch ev.Type {
	case EventQuit:
		a.running = false
	case EventClose:
		if a.OnClose == nil || a.OnClose(a) {
			a.running = false
		}
	}
}

func (a *App) Step(draw func(ui *UI), wait bool) bool {
	if !a.running {
		return false
	}
	busy := a.every || a.Animate || a.dirty > 0 || a.UI.Animating()
	if a.dirty > 0 {
		a.dirty--
	}
	timeout := 0
	if wait {
		if busy {
			timeout = a.pace.Wait()
		} else {
			timeout = a.IdleMs
		}
	}
	for ev := range a.G.Events(timeout) {
		a.handle(ev)
		a.Redraw()
		if !a.running {
			return false
		}
	}
	if a.running && (busy || a.dirty > 0) {
		a.paint(draw)
	}
	return a.running
}

func (a *App) paint(draw func(ui *UI)) {
	ui := a.UI
	if ui.Begin(0) == nil {
		ui.End()
		return
	}
	defer ui.End()
	ui.App(a.Title, a.Flags)
	defer ui.End()
	if draw != nil {
		draw(ui)
	}
	a.Frames++
}

func (a *App) Stop() {
	a.running = false
	if a.UI != nil {
		a.UI.Free()
		a.UI = nil
	}
	if a.G != nil {
		a.G.Close()
		a.G = nil
	}
	if a.pace != nil {
		a.pace.Free()
		a.pace = nil
	}
}

func (a *App) Run(draw func(ui *UI)) error {
	if err := a.Start(); err != nil {
		return err
	}
	defer a.Stop()
	for a.Step(draw, true) {
	}
	return nil
}
