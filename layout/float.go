package layout

import "github.com/egoist/mygo/ui"

// Float offsets content over its parent; offset springs.
func Float(c *ui.Context, x, y float32, content func()) ui.Element {
	if !finite(x) || !finite(y) {
		panic("layout: Float offset must be finite")
	}
	b := ui.Box(c).Absolute()
	return b.Left(b.AnimateTo("x", x, ui.Snappy)).Top(b.AnimateTo("y", y, ui.Snappy)).Children(content)
}
