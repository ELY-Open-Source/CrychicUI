package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Placeholder is a dashed box whose dash color animates.
func Placeholder(c *ui.Context, label string, w, h float32) ui.Element {
	if !finite(float64(w)) || !finite(float64(h)) || w <= 0 || h <= 0 {
		panic("display: Placeholder size must be finite and positive")
	}
	t := c.Theme()
	root := ui.Box(c).Size(w, h).Center().Radius(t.Radius).BorderWidth(1).BorderStyle(ui.BorderDashed).TextColor(t.TextMuted)
	border := t.Border
	if root.Hovered() {
		border = t.Accent
	}
	root.BorderColor(root.AnimateColor("border", border, ui.Snappy)).Children(func() { kit.Muted(c, label) })
	return root
}
