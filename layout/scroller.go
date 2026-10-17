package layout

import (
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/icon"
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/kit"
	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Scroller scrolls sideways; arrows move it a page.
func Scroller(c *ui.Context, content func()) (root, back, forward ui.Element) {
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
	st := ui.Local(root, "scroll", func() ui.ScrollState { return ui.ScrollState{} })
	kit.Settle(c, root, []float32{st.X, st.MaxX}, func() []float32 { return []float32{st.X, st.MaxX} })
	root.Children(func() {
		arrow := func(glyph, label string, off bool) ui.Element {
			b := ui.Button(c, "").Label(label).Disabled(off).Children(func() {
				ui.Icon(c, icon.Must(glyph)).Size(16, 16)
			})
			return kit.Press(b)
		}
		back = arrow("chevron-left", "Scroll back", st.X <= 0)
		view := ui.ScrollHorizontal(c).Grow(1).Basis(0).TrackScroll(st)
		view.Children(func() { ui.Row(c).Gap(t.Space(3)).Children(content) })
		forward = arrow("chevron-right", "Scroll forward", st.X >= st.MaxX)
		page := view.LayoutBounds().W
		if kit.Clicked(back, sound.Tap) {
			st.ScrollTo(max(0, st.X-page), 0, ui.Smooth)
			kit.Log("Scroller", "back", "x", st.X)
		}
		if kit.Clicked(forward, sound.Tap) {
			st.ScrollTo(min(st.MaxX, st.X+page), 0, ui.Smooth)
			kit.Log("Scroller", "forward", "x", st.X)
		}
	})
	return root, back, forward
}
