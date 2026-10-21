package display

import "github.com/egoist/mygo/ui"

// Indicator is a colored dot animating to new states.
func Indicator(c *ui.Context, color ui.Color) ui.Element {
	dot := ui.Box(c).AlignSelf(ui.Start).Size(c.Theme().Space(2.5), c.Theme().Space(2.5)).Radius(999)
	return dot.Background(dot.AnimateColor("color", color, ui.Snappy))
}
