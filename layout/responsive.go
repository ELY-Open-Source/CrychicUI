package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Responsive rows items at minWidth or wider, else stacks them.
func Responsive(c *ui.Context, minWidth float32, items []func()) ui.Element {
	if !finite(minWidth) || minWidth <= 0 {
		panic("layout: Responsive minWidth must be finite and positive")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3))
	w := root.LayoutBounds().W
	kit.Settle(c, root, []float32{w}, func() []float32 { return []float32{root.LayoutBounds().W} })
	wide := w >= minWidth
	if wide {
		root.Row().AlignItems(ui.Start)
	}
	return root.Children(func() {
		for i, it := range items {
			b := ui.Column(c).Key(i).Transition(ui.ElementTransition{Position: true, Size: true})
			if wide {
				b.Grow(1).Basis(0)
			}
			b.Children(it)
		}
	})
}
