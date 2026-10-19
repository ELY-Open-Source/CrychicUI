package layout

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Layers stacks layers on one box; each fades in turn.
func Layers(c *ui.Context, layers []func()) ui.Element {
	if len(layers) == 0 {
		panic("layout: Layers needs a layer")
	}
	return ui.Box(c).Children(func() {
		for i, l := range layers {
			b := ui.Box(c).Key(i)
			if i > 0 {
				b.Absolute().Top(0).Right(0).Bottom(0).Left(0)
			}
			kit.Arrive(b, ui.Pose{Fade: 1}, time.Duration(i)*80*time.Millisecond).Children(l)
		}
	})
}
