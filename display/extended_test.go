package display

import (
	"fmt"
	"testing"

	"github.com/egoist/mygo/ui"
)

func TestFeature(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Feature(c, "sparkles", "Fast", "Motion and sound") }, 300, 200)
	settles(t, tt, "Feature")
	find(t, tt, "Fast")
	find(t, tt, "Motion and sound")
}

func TestMetric(t *testing.T) {
	value := 12.0
	tt := ui.NewTester(func(c *ui.Context) { Metric(c, value, "ms") }, 300, 200)
	value = 42
	tt.Frame()
	settles(t, tt, "Metric")
	find(t, tt, "42 ms")
	panics(t, "Metric unit", func() { ui.Render(func(c *ui.Context) { Metric(c, 1, "") }, 10, 10, 1) })
}

func TestStars(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) { Stars(c, 2.5, StarsOptions{}) }, 300, 200)
	settles(t, tt, "Stars")
	panics(t, "Stars rating", func() { ui.Render(func(c *ui.Context) { Stars(c, 6, StarsOptions{}) }, 10, 10, 1) })
}

func TestPrice(t *testing.T) {
	current := 10.0
	tt := ui.NewTester(func(c *ui.Context) { Price(c, current, 12, "USD") }, 300, 200)
	find(t, tt, "USD 10.00")
	find(t, tt, "USD 12.00")
	current = 11
	tt.Frame()
	settles(t, tt, "Price")
	find(t, tt, "USD 11.00")
	panics(t, "Price order", func() { ui.Render(func(c *ui.Context) { Price(c, 20, 10, "USD") }, 10, 10, 1) })
}

func TestVersion(t *testing.T) {
	version := "1.2.3"
	tt := ui.NewTester(func(c *ui.Context) { Version(c, version) }, 300, 200)
	version = "v2.0.0"
	tt.Frame()
	settles(t, tt, "Version")
	find(t, tt, "v2.0.0")
	panics(t, "Version empty", func() { ui.Render(func(c *ui.Context) { Version(c, "") }, 10, 10, 1) })
}

func TestLegend(t *testing.T) {
	var accent, success ui.Color
	tt := ui.NewTester(func(c *ui.Context) {
		accent, success = c.Theme().Accent, c.Theme().Success
		Legend(c, []LegendItem{{Label: "Now", Color: accent}, {Label: "Done", Color: success}})
	}, 300, 200)
	settles(t, tt, "Legend")
	find(t, tt, "Now")
	find(t, tt, "Done")
}

func TestSwatch(t *testing.T) {
	var accent, success ui.Color
	color := ui.Color{}
	tt := ui.NewTester(func(c *ui.Context) {
		accent, success = c.Theme().Accent, c.Theme().Success
		if color.A == 0 {
			color = accent
		}
		Swatch(c, color)
	}, 300, 200)
	find(t, tt, fmt.Sprintf("#%02X%02X%02X", accent.R, accent.G, accent.B))
	tt.Frame()
	color = success
	tt.Frame()
	settles(t, tt, "Swatch color")
	panics(t, "Swatch color", func() { ui.Render(func(c *ui.Context) { Swatch(c, ui.Color{}) }, 10, 10, 1) })
}

func TestStepList(t *testing.T) {
	current := 0
	tt := ui.NewTester(func(c *ui.Context) { StepList(c, []string{"One", "Two"}, current) }, 300, 200)
	current = 1
	tt.Frame()
	settles(t, tt, "StepList current")
	find(t, tt, "Two")
	panics(t, "StepList current", func() { ui.Render(func(c *ui.Context) { StepList(c, []string{"One"}, 1) }, 10, 10, 1) })
}

func TestChecklist(t *testing.T) {
	items := []ChecklistItem{{Label: "Write", Done: false}}
	tt := ui.NewTester(func(c *ui.Context) { Checklist(c, items) }, 300, 200)
	items[0].Done = true
	tt.Frame()
	settles(t, tt, "Checklist check")
	find(t, tt, "Write")
}
