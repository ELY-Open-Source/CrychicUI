package layout

import "github.com/egoist/mygo/ui"

// Cluster wraps items aligned Start, Center or End; items glide.
func Cluster(c *ui.Context, align ui.Align, items []func()) ui.Element {
	if align != ui.Start && align != ui.Center && align != ui.End {
		panic("layout: Cluster align must be Start, Center or End")
	}
	return ui.Row(c).FillWidth().Wrap().Justify(align).Gap(c.Theme().Space(2)).Children(func() {
		for i, it := range items {
			ui.Box(c).Key(i).Transition(ui.ElementTransition{Position: true}).Children(it)
		}
	})
}
