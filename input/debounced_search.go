package input

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DebouncedSearchOptions tune a DebouncedSearch.
type DebouncedSearchOptions struct {
	// Placeholder and Searching are copy; Delay defaults to 300 ms.
	Placeholder, Searching string
	Delay                  time.Duration
}

// DebouncedSearch reports a query after a pause.
func DebouncedSearch(c *ui.Context, label string, query, reported *string, o DebouncedSearchOptions) ui.Element {
	if query == nil || reported == nil {
		panic("input: DebouncedSearch needs query and reported values")
	}
	if o.Placeholder == "" {
		o.Placeholder = "Search"
	}
	if o.Searching == "" {
		o.Searching = "Searching"
	}
	if o.Delay == 0 {
		o.Delay = 300 * time.Millisecond
	}
	if o.Delay < 0 {
		panic("input: DebouncedSearch delay must be positive")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	pending := ui.Local(root, "pending", func() bool { return false })
	due := ui.Local(root, "due", func() time.Time { return time.Time{} })
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		in := ui.TextInput(c, query).FillWidth().Placeholder(o.Placeholder).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		if changed(root, in) {
			*pending = true
			*due = c.Now().Add(o.Delay)
			c.After(o.Delay)
		}
		if *pending && !c.Now().Before(*due) {
			*pending = false
			if *reported != *query {
				*reported = *query
				root.SetChanged()
				kit.Log("DebouncedSearch", "reported", "count", len([]rune(*reported)))
			}
		} else if *pending {
			c.After(due.Sub(c.Now()))
		}
		if *pending {
			kit.Muted(c, o.Searching).FontSize(t.Rem(0.85)).ContentTransition(ui.ContentTransition{})
		}
	})
}
