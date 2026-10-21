package layout

import "github.com/egoist/mygo/ui"

// Gutter reserves n theme units; size changes transition.
func Gutter(c *ui.Context, n float32) ui.Element {
	if !finite(n) || n < 0 {
		panic("layout: Gutter n must be finite and non-negative")
	}
	return ui.Box(c).Basis(c.Theme().Space(n)).Shrink(0).Transition(ui.ElementTransition{Size: true})
}
