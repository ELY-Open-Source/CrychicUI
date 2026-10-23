package media

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CollageOptions tune a Collage.
type CollageOptions struct {
	// Gap between images; zero gives the default.
	Gap float32
	// Height of the collage; zero gives the default.
	Height float32
}

// Collage sets one large image beside the rest.
func Collage(c *ui.Context, label string, images []*ui.Bitmap, o CollageOptions) ui.Element {
	pictures("Collage", label, images)
	t := c.Theme()
	size("Collage", label, &o.Gap, t.Space(2))
	size("Collage", label, &o.Height, t.Space(64))
	root := ui.Row(c).FillWidth().Height(o.Height).Gap(o.Gap).Label(label)
	tile := func(i int) ui.Element {
		e := ui.Box(c).Radius(t.Radius).Clip().Background(t.Surface).Children(func() { ui.Image(c, images[i]).Fill().Fit(ui.Cover) })
		return kit.Arrive(e, ui.Pose{Fade: 1, Scale: 0.96}, time.Duration(i)*50*time.Millisecond)
	}
	return root.Children(func() {
		c.Key(0)
		tile(0).Grow(2).Basis(0).Height(o.Height)
		if len(images) == 1 {
			return
		}
		c.Key("rest")
		ui.Column(c).Grow(1).Basis(0).Height(o.Height).Gap(o.Gap).Children(func() {
			for i := 1; i < len(images); i++ {
				c.Key(i)
				tile(i).FillWidth().Grow(1).Basis(0)
			}
		})
	})
}
