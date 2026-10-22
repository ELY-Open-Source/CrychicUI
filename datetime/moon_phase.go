package datetime

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MoonPhaseOptions tune a MoonPhase.
type MoonPhaseOptions struct {
	// Names name the eight phases from new moon.
	Names [8]string
}

// synodic is the days from new moon to new moon.
const synodic = 29.530588853

// newMoon is a known new moon.
var newMoon = time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)

// MoonPhase draws the moon as lit on day.
func MoonPhase(c *ui.Context, day time.Time, o MoonPhaseOptions) ui.Element {
	if day.IsZero() {
		panic("datetime: MoonPhase needs a day")
	}
	if o.Names == ([8]string{}) {
		o.Names = [8]string{"New moon", "Waxing crescent", "First quarter", "Waxing gibbous", "Full moon",
			"Waning gibbous", "Last quarter", "Waning crescent"}
	}
	t := c.Theme()
	cycles := day.Sub(newMoon).Hours() / 24 / synodic
	name := o.Names[int((cycles-math.Floor(cycles))*8+0.5)%8]
	size := t.Space(16)
	root := ui.Box(c).Size(size, size).Shrink(0).Role(ui.RoleImage).Label(name).Tooltip(name)
	shown := float64(root.AnimateTo("cycles", float32(cycles), ui.Smooth))
	return kit.Arrive(root, ui.Pose{Scale: 0.8, Fade: 1}, 0).Draw(func(p *ui.Painter, r ui.Rect) {
		f := shown - math.Floor(shown)
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-1
		var disc ui.Path
		p.FillPath(disc.Circle(cx, cy, rad), t.Surface)
		p.StrokePath(&disc, 1, t.Border)
		side, k := float32(1), float32(math.Cos(2*math.Pi*f))
		if f > 0.5 {
			side = -1
		}
		const n = 32
		var lit ui.Path
		for i := 0; i <= n; i++ {
			a := math.Pi * float64(i) / n
			x, y := cx+side*rad*float32(math.Sin(a)), cy-rad*float32(math.Cos(a))
			if i == 0 {
				lit.MoveTo(x, y)
			} else {
				lit.LineTo(x, y)
			}
		}
		for i := n; i >= 0; i-- {
			a := math.Pi * float64(i) / n
			lit.LineTo(cx+side*k*rad*float32(math.Sin(a)), cy-rad*float32(math.Cos(a)))
		}
		p.FillPath(lit.Close(), t.Warning)
	})
}
