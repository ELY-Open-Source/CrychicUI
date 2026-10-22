package motion

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FollowerOptions tune a Follower.
type FollowerOptions struct {
	// Height is the area's height; t.Space(40) when zero.
	Height float32
}

var lag = ui.Animation{Spring: ui.SpringOptions{Duration: 600 * time.Millisecond, Bounce: 0.3}}

// Follower's dot trails the pointer on a spring.
func Follower(c *ui.Context, o FollowerOptions) ui.Element {
	t := c.Theme()
	root := ui.Box(c).FillWidth().Height(positive("Follower", "Height", o.Height, t.Space(40))).Radius(t.Radius).
		Background(t.Surface).Clip()
	size := func() []float32 { b := root.LayoutBounds(); return []float32{b.W, b.H} }
	used := size()
	kit.Settle(c, root, used, size)
	x, y, over := root.PointerPosition()
	if !over {
		x, y = used[0]/2, used[1]/2
	}
	d := t.Space(4)
	return root.Children(func() {
		dot := ui.Box(c).Absolute().Top(0).Left(0).Size(d, d).Radius(999).Background(t.Accent)
		dot.Translate(dot.AnimateTo("x", x-d/2, lag), dot.AnimateTo("y", y-d/2, lag))
	})
}
