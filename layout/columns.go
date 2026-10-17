package layout

import "github.com/egoist/mygo/ui"

// Columns gives each item an equal column; widths animate.
func Columns(c *ui.Context, items []func()) ui.Element {
	if len(items) == 0 {
		panic("layout: Columns needs an item")
	}
	return ui.Row(c).FillWidth().Gap(c.Theme().Space(4)).AlignItems(ui.Start).Children(func() {
		for i, it := range items {
			ui.Column(c).Key(i).Grow(1).Basis(0).Transition(ui.ElementTransition{Size: true, Position: true}).
				Children(it)
		}
	})
}
