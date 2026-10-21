package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Stage scales fixed-size content to fit, letterboxed; scale springs.
func Stage(c *ui.Context, w, h float32, content func()) ui.Element {
	if !finite(w) || !finite(h) || w <= 0 || h <= 0 {
		panic("layout: Stage size must be finite and positive")
	}
	t := c.Theme()
	root := ui.Box(c).Fill().Center().Clip().Radius(t.Radius).Background(t.Surface)
	r := root.LayoutBounds()
	kit.Settle(c, root, []float32{r.W, r.H}, func() []float32 { b := root.LayoutBounds(); return []float32{b.W, b.H} })
	k := float32(1)
	if r.W > 0 && r.H > 0 {
		k = min(r.W/w, r.H/h)
	}
	return root.Children(func() {
		b := ui.Box(c).Size(w, h).Shrink(0)
		s := b.AnimateTo("k", k, ui.Snappy)
		b.Scale(s, s).Children(content)
	})
}
