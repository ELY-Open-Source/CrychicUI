package display

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StarsOptions tune Stars.
type StarsOptions struct {
	// Max is the rating scale; zero is five.
	Max int
}

// Stars is a read-only rating with halves and staggered fills.
func Stars(c *ui.Context, rating float64, o StarsOptions) ui.Element {
	if o.Max == 0 {
		o.Max = 5
	}
	if o.Max <= 0 || !finite(rating) || rating < 0 || rating > float64(o.Max) {
		panic("display: Stars rating must fit the scale")
	}
	t := c.Theme()
	return ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).Children(func() {
		for i := range o.Max {
			fill := min(1, max(0, float32(rating)-float32(i)))
			star := ui.Box(c).Key(i).Size(t.Space(4), t.Space(4))
			kit.Arrive(star, ui.Pose{Scale: 0.5, Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				ui.Icon(c, icon.Must("star")).Fill().TextColor(t.Border)
				if fill > 0 {
					ui.Box(c).Absolute().Left(0).Top(0).Width(t.Space(4) * fill).Height(t.Space(4)).ClipX().Children(func() {
						ui.Icon(c, icon.Filled("star")).Size(t.Space(4), t.Space(4)).TextColor(t.Warning)
					})
				}
			})
		}
	})
}
