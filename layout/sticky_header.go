package layout

import (
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StickyHeader tops a scroll; its shadow fades in once scrolled.
func StickyHeader(c *ui.Context, header, content func()) ui.Element {
	t := c.Theme()
	root := ui.Column(c).Clip()
	st := ui.Local(root, "scroll", func() ui.ScrollState { return ui.ScrollState{} })
	kit.Settle(c, root, []float32{st.Y}, func() []float32 { return []float32{st.Y} })
	return root.Children(func() {
		head := ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2), t.Space(3)).AlignItems(ui.Center).
			Background(t.Background)
		target := float32(0)
		if st.Y > 0 {
			target = 1
		}
		a := head.AnimateTo("shadow", target, ui.Snappy)
		head.Shadow(0, 1, 4, 0, kit.Shade(0.1*a)).Children(header)
		ui.Scroll(c).Grow(1).Basis(0).TrackScroll(st).Children(content)
	})
}
