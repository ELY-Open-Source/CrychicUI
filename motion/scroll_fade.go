package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ScrollFadeOptions tune a ScrollFade.
type ScrollFadeOptions struct {
	// Distance is how far the scroll goes until gone; t.Space(20).
	Distance float32
}

// ScrollFade fades content as state scrolls past it.
func ScrollFade(c *ui.Context, state *ui.ScrollState, o ScrollFadeOptions, content func()) ui.Element {
	if state == nil || content == nil {
		panic("motion: ScrollFade needs a scroll state and content")
	}
	d := positive("ScrollFade", "Distance", o.Distance, c.Theme().Space(20))
	k := min(max(state.Y/d, 0), 1)
	root := ui.Box(c).Opacity(1-k).Translate(0, -k*d/4)
	kit.Settle(c, root, []float32{state.Y}, func() []float32 { return []float32{state.Y} })
	return root.Children(content)
}
