package motion

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TextShimmer loops a bright band over muted text.
func TextShimmer(c *ui.Context, text string) ui.Element {
	if text == "" {
		panic("motion: TextShimmer needs text")
	}
	t := c.Theme()
	band := t.Space(10)
	root := ui.Box(c).AlignSelf(ui.Start)
	width := func() []float32 { return []float32{root.LayoutBounds().W} }
	used := width()
	kit.Settle(c, root, used, width)
	x := -band + (used[0]+band)*root.Loop("crychic.shimmer", 2*time.Second, ui.Linear)
	return root.Children(func() {
		ui.Text(c, text).NoWrap().TextColor(t.TextMuted)
		ui.Box(c).Absolute().Top(0).Bottom(0).Left(0).Width(band).ClipX().Translate(x, 0).Children(func() {
			ui.Text(c, text).NoWrap().Width(used[0]).Shrink(0).Translate(-x, 0).Role(ui.RoleNone)
		})
	})
}
