package gallery

import (
	"log/slog"
	"strings"
	"sync"

	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

const keepEvents = 6

// App is the gallery that desktop and browser share.
type App struct {
	router *ui.Router
	query  string
	mode   int
	events []string
	codes  map[string]*editor.State

	mu   sync.Mutex
	snap Snapshot
}

// Snapshot is what the browser hook reads between frames.
type Snapshot struct {
	Page  string
	Built bool
}

// New opens the gallery at path.
func New(path string) *App {
	return &App{router: ui.NewRouter(path), codes: map[string]*editor.State{}}
}

// Snapshot returns the state of the last built frame.
func (a *App) Snapshot() Snapshot {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.snap
}

func (a *App) View(c *ui.Context) {
	switch a.mode {
	case 1:
		c.SetTheme(ui.LightTheme())
	case 2:
		c.SetTheme(ui.DarkTheme())
	}
	ui.Column(c).Fill().Children(func() {
		a.topBar(c)
		ui.Row(c).Grow(1).AlignItems(ui.Stretch).Children(func() {
			a.sidebar(c)
			ui.Column(c).Grow(1).Children(func() {
				if strings.TrimSpace(a.query) != "" {
					a.results(c)
					return
				}
				a.pages(c)
			})
		})
		a.statusBar(c)
	})
	a.mu.Lock()
	a.snap = Snapshot{Page: a.router.Location(), Built: true}
	a.mu.Unlock()
}

// Go opens path and clears the search.
func (a *App) Go(path string) {
	a.query = ""
	if path != a.router.Location() {
		slog.Debug("crychic: gallery go", "path", path)
		a.router.Push(path)
	}
}

// logger keeps a demo's events for the status strip.
func (a *App) logger(e Entry) func(string) {
	return func(ev string) {
		// The strip shows the event; the log never does.
		slog.Debug("crychic: gallery event", "entry", e.Name)
		a.events = append(a.events, e.Name+": "+ev)
		a.events = a.events[max(0, len(a.events)-keepEvents):]
	}
}

func (a *App) code(e Entry) *editor.State {
	st, ok := a.codes[e.Path()]
	if !ok {
		buf := editor.NewBuffer(e.Code())
		buf.SetLanguage(editor.LanguageFor("demo.go"))
		st = &editor.State{Buffer: buf, ReadOnly: true}
		a.codes[e.Path()] = st
	}
	return st
}

func matches(q string) []Entry {
	q = strings.ToLower(strings.TrimSpace(q))
	var out []Entry
	for _, e := range entries {
		if strings.Contains(strings.ToLower(e.Name+" "+e.Category+" "+e.Summary), q) {
			out = append(out, e)
		}
	}
	return out
}
