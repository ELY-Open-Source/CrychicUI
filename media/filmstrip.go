package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FilmstripOptions tune a Filmstrip.
type FilmstripOptions struct {
	// Height of a frame; zero gives the default.
	Height float32
}

// Filmstrip shows frames on film under a playhead at pos.
func Filmstrip(c *ui.Context, label string, frames []*ui.Bitmap, pos, total time.Duration, o FilmstripOptions) ui.Element {
	pictures("Filmstrip", label, frames)
	if total <= 0 || pos < 0 || pos > total {
		panic(fmt.Sprintf("media: Filmstrip %q needs a total above 0 and pos within it", label))
	}
	t := c.Theme()
	size("Filmstrip", label, &o.Height, t.Space(16))
	hole := t.Space(1.5)
	band := hole + t.Space(1.5)
	ink := t.Background
	if t.Dark {
		ink = t.Text
	}
	strip := ui.Row(c).FillWidth().Height(o.Height+2*band).Gap(t.Space(1)).Padding(band, t.Space(1)).Radius(t.Radius).
		Background(kit.Shade(0.86)).Role(ui.RoleImage).Label(label)
	strip.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		for x := r.X + hole; x+hole <= r.X+r.W; x += hole * 2.5 {
			p.Fill(ui.Rect{X: x, Y: r.Y + (band-hole)/2, W: hole, H: hole}, ink.Alpha(0.6), hole*0.3)
			p.Fill(ui.Rect{X: x, Y: r.Y + r.H - band + (band-hole)/2, W: hole, H: hole}, ink.Alpha(0.6), hole*0.3)
		}
	}))
	f := strip.AnimateTo("head", float32(pos.Seconds()/total.Seconds()), ui.Snappy)
	return strip.Children(func() {
		for i, img := range frames {
			c.Key(i)
			fr := ui.Box(c).Grow(1).Basis(0).Height(o.Height).Radius(t.Space(0.5)).Clip().Children(func() {
				ui.Image(c, img).Fill().Fit(ui.Cover)
			})
			kit.Arrive(fr, ui.Pose{Fade: 1, X: t.Space(2)}, time.Duration(i)*30*time.Millisecond)
		}
		c.Key("head")
		w := t.Space(0.5)
		dx := -w / 2
		if c.Direction() == ui.RightToLeft {
			dx = w / 2
		}
		ui.Box(c).Absolute().Top(0).Bottom(0).Width(w).StartPercent(100*f).Translate(dx, 0).Background(t.Accent)
	})
}
