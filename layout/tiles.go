package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Tiles lays square tiles in a wrapping grid; tiles lift.
func Tiles(c *ui.Context, size float32, items []func()) ui.Element {
	if !finite(size) || size <= 0 {
		panic("layout: Tiles size must be finite and positive")
	}
	t := c.Theme()
	return ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).Children(func() {
		for i, it := range items {
			b := ui.Column(c).Key(i).Size(size, size).Center().Radius(t.Radius).Background(t.Surface).Border(1, t.Border)
			kit.Lift(b).Children(it)
		}
	})
}
