package layout

import (
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AutoGrid fits as many columns as minWidth allows.
func AutoGrid(c *ui.Context, minWidth float32, items []func()) ui.Element {
	if !finite(minWidth) || minWidth <= 0 {
		panic("layout: AutoGrid minWidth must be finite and positive")
	}
	gap := c.Theme().Space(3)
	root := ui.Grid(c).FillWidth().Gap(gap)
	w := root.LayoutBounds().W
	kit.Settle(c, root, []float32{w}, func() []float32 { return []float32{root.LayoutBounds().W} })
	return root.Columns(max(1, int((w+gap)/(minWidth+gap)))).Children(func() {
		for i, it := range items {
			ui.Box(c).Key(i).Transition(ui.ElementTransition{Position: true, Size: true}).Children(it)
		}
	})
}
