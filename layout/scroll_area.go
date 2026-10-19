package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ScrollArea scrolls content; edge shadows fade with the offset.
func ScrollArea(c *ui.Context, content func()) ui.Element {
	root := ui.Box(c).Clip()
	st := ui.Local(root, "scroll", func() ui.ScrollState { return ui.ScrollState{} })
	kit.Settle(c, root, []float32{st.Y, st.MaxY}, func() []float32 { return []float32{st.Y, st.MaxY} })
	return root.Children(func() {
		ui.Scroll(c).Fill().TrackScroll(st).Children(content)
		edge := func(key string, show bool, angle float32) {
			e := ui.Box(c).Key(key).PassThrough().Absolute().Left(0).Right(0).Height(c.Theme().Space(3))
			target := float32(0)
			if show {
				target = 1
			}
			e.Gradient(kit.Shade(0.12), kit.Shade(0), angle).Opacity(e.AnimateTo("o", target, ui.Snappy))
			if key == "top" {
				e.Top(0)
			} else {
				e.Bottom(0)
			}
		}
		edge("top", st.Y > 0, 180)
		edge("bottom", st.Y < st.MaxY, 0)
	})
}
