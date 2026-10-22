package chart

import "github.com/egoist/mygo/ui"

// GroupedBarsOptions tune a GroupedBars.
type GroupedBarsOptions = LineOptions

// GroupedBars sets each series' bars side by side per label.
func GroupedBars(c *ui.Context, label string, labels []string, series []Series, at *int, o GroupedBarsOptions) ui.Element {
	return grouped(c, "GroupedBars", label, labels, series, at, o)
}
