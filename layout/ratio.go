package layout

import "github.com/egoist/mygo/ui"

// Ratio splits a row's width by weights; weights spring.
func Ratio(c *ui.Context, weights []float32, items []func()) ui.Element {
	if len(weights) == 0 || len(weights) != len(items) {
		panic("layout: Ratio needs one weight per item")
	}
	sum := float32(0)
	for _, w := range weights {
		if !finite(w) || w < 0 {
			panic("layout: Ratio weights must be finite and non-negative")
		}
		sum += w
	}
	if sum == 0 {
		panic("layout: Ratio weights sum to zero")
	}
	return ui.Row(c).FillWidth().Gap(c.Theme().Space(3)).AlignItems(ui.Stretch).Children(func() {
		for i, it := range items {
			b := ui.Column(c).Key(i).Basis(0).ClipX()
			b.Grow(b.AnimateTo("w", weights[i], ui.Snappy)).Children(it)
		}
	})
}
