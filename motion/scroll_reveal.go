package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ScrollReveal raises each row in as it first shows.
func ScrollReveal(c *ui.Context, n int, row func(i int)) ui.Element {
	if n < 0 || row == nil {
		panic("motion: ScrollReveal needs n ≥ 0 and a row")
	}
	t := c.Theme()
	return ui.Column(c).Gap(t.Space(2)).Children(func() {
		for i := range n {
			c.Key(i)
			w := ui.Box(c)
			seen := ui.Local(w, "seen", func() bool { return false })
			if w.InView() {
				*seen = true
			}
			if !*seen {
				kit.Settle(c, w, []float32{0}, func() []float32 { return []float32{unit(w.InView())} })
			}
			k := w.AnimateTo("in", unit(*seen), ui.Smooth)
			w.Children(func() { ui.Box(c).Opacity(k).Translate(0, (1-k)*t.Space(4)).Children(func() { row(i) }) })
		}
	})
}
