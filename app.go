package goome

import "strings"

// App is the easy way to run a window: it opens the window and the UI, runs the loop and cleans up.
type App struct {
	// Title is the window title.
	Title string
	// Width is the starting window width.
	Width int
	// Height is the starting window height.
	Height int
	// Flags are the App* bits for the root panel, like the title bar, menu bar and status bar.
	Flags int
	// FPS is the frame rate while drawing, 60 by default.
	FPS int
	// IdleMs is how long Step waits for an event when nothing needs drawing, 500 by default.
	IdleMs int
	// Light switches to the light theme at start.
	Light bool
	// Animate keeps redrawing every frame even when nothing happened.
	Animate bool
	// Options are extra Open settings. Title, Mode and a zero size or WindowFlags are filled in from the other fields.
	Options Options

	// OnStart runs once the window and UI are up.
	OnStart func(a *App)
	// OnEvent gets every event before the UI does. Return true to eat it.
	OnEvent func(a *App, ev Event) bool
	// OnClose runs when the window is asked to close. Return false to keep it open.
	OnClose func(a *App) bool

	// G is the underlying Goome, set once Start has run.
	G *Goome
	// UI is the UI kit for the window, set once Start has run.
	UI *UI
	// Frames counts the frames drawn so far.
	Frames int

	running bool
	dirty   int
	every   bool
	pace    *Pacer
}

// NewApp makes an App with the given title, size and App* flags, 60 FPS and a 500ms idle wait.
func NewApp(title string, w, h, flags int) *App {
	return &App{Title: title, Width: w, Height: h, Flags: flags, FPS: 60, IdleMs: 500}
}

// Redraw asks for the next frame to be drawn.
func (a *App) Redraw() {
	if a.dirty < 1 {
		a.dirty = 1
	}
}

// Quit ends the loop after the current step.
func (a *App) Quit() { a.running = false }

// Running reports whether the loop is still going.
func (a *App) Running() bool { return a.running }

// Start opens the window and the UI without running the loop, for when you drive Step yourself.
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

// Step handles events and draws one frame if needed, and returns false once the app is done. With wait it sleeps until something happens.
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

// Stop frees the UI, the window and the pacer. It is safe to call twice.
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

// Run opens the window, calls draw every frame it needs one until the app quits, then cleans up.
func (a *App) Run(draw func(ui *UI)) error {
	if err := a.Start(); err != nil {
		return err
	}
	defer a.Stop()
	for a.Step(draw, true) {
	}
	return nil
}
