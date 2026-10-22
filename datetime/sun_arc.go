package datetime

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// SunArc moves the sun along an arc, sunrise to sunset.
func SunArc(c *ui.Context, sunrise, sunset time.Time) ui.Element {
	if !sunset.After(sunrise) {
		panic("datetime: SunArc needs a sunset after its sunrise")
	}
	t := c.Theme()
	loc := locale(c)
	now := c.Now()
	span := sunset.Sub(sunrise)
	f := float32(min(max(now.Sub(sunrise), 0), span)) / float32(span)
	if now.Before(sunset) {
		c.After(time.Minute)
	}
	at := func(x time.Time) string { return i18n.FormatDate(loc, x, i18n.DateOptions{Skeleton: "jm"}) }
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(1)).Role(ui.RoleMeter).Label(at(sunrise)+" – "+at(sunset)).
		Range(0, 1, float64(f))
	return root.Children(func() {
		arc := ui.Box(c).Size(t.Space(42), t.Space(24))
		shown := arc.AnimateTo("f", f, ui.Smooth)
		up := f > 0 && f < 1
		arc.Draw(func(p *ui.Painter, r ui.Rect) {
			pad := t.Space(3)
			rad := min(r.W/2-pad, r.H-2*pad)
			cx, cy := r.X+r.W/2, r.Y+r.H-pad
			path := func() *ui.Path {
				var l ui.Path
				for i := 0; i <= 48; i++ {
					a := math.Pi * (1 - float64(i)/48)
					x, y := cx+rad*float32(math.Cos(a)), cy-rad*float32(math.Sin(a))
					if i == 0 {
						l.MoveTo(x, y)
					} else {
						l.LineTo(x, y)
					}
				}
				return &l
			}
			segment(p, r.X+pad, cy, r.X+r.W-pad, cy, 1, t.Border)
			p.StrokePath(path(), 1, t.Border)
			if shown > 0 {
				p.StrokePath(path().Trim(0, shown), t.Space(0.75), t.Warning)
			}
			sun := t.TextMuted
			if up {
				sun = t.Warning
			}
			a := math.Pi * (1 - float64(shown))
			var dot ui.Path
			p.FillPath(dot.Circle(cx+rad*float32(math.Cos(a)), cy-rad*float32(math.Sin(a)), t.Space(2)), sun)
		})
		ui.Row(c).FillWidth().Justify(ui.SpaceBetween).Children(func() {
			kit.Muted(c, at(sunrise))
			kit.Muted(c, at(sunset))
		})
	})
}
