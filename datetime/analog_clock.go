package datetime

import (
	"fmt"
	"math"
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// AnalogClockOptions tune an AnalogClock.
type AnalogClockOptions struct {
	// Location is the zone shown; local time when nil.
	Location *time.Location
	// Size is the face's width; 40 spaces when zero.
	Size float32
}

// AnalogClock shows the time with sweeping hands.
func AnalogClock(c *ui.Context, o AnalogClockOptions) ui.Element {
	t := c.Theme()
	if o.Size == 0 {
		o.Size = t.Space(40)
	}
	if !(o.Size > 0) || math.IsInf(float64(o.Size), 0) {
		panic(fmt.Sprintf("datetime: AnalogClock size %v must be finite and positive", o.Size))
	}
	if o.Location == nil {
		o.Location = time.Local
	}
	now := c.Now().In(o.Location)
	still := c.Preferences().ReduceMotion
	if still {
		c.After(time.Second - time.Duration(now.Nanosecond()))
	} else {
		c.After(time.Minute - time.Duration(now.Second())*time.Second - time.Duration(now.Nanosecond()))
	}
	label := i18n.FormatDate(locale(c), now, i18n.DateOptions{Skeleton: "jm"})
	return ui.Box(c).Size(o.Size, o.Size).Shrink(0).Role(ui.RoleImage).Label(label).Draw(func(p *ui.Painter, r ui.Rect) {
		at := p.Now().In(o.Location)
		sec := float64(at.Second())
		if !still {
			sec += float64(at.Nanosecond()) / 1e9
			p.AnimationFrame()
		}
		cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2-1
		var face ui.Path
		face.Circle(cx, cy, rad)
		p.FillPath(&face, t.Surface)
		p.StrokePath(&face, 1, t.Border)
		hand := func(turn float64, length, width float32, col ui.Color) {
			a := 2*math.Pi*turn - math.Pi/2
			segment(p, cx, cy, cx+length*float32(math.Cos(a)), cy+length*float32(math.Sin(a)), width, col)
		}
		for i := range 12 {
			a := 2*math.Pi*float64(i)/12 - math.Pi/2
			in := rad * 0.88
			if i%3 == 0 {
				in = rad * 0.8
			}
			cs, sn := float32(math.Cos(a)), float32(math.Sin(a))
			segment(p, cx+in*cs, cy+in*sn, cx+rad*0.95*cs, cy+rad*0.95*sn, t.Space(0.5), t.TextMuted)
		}
		minute := float64(at.Minute()) + sec/60
		hand((float64(at.Hour()%12)+minute/60)/12, rad*0.5, t.Space(1), t.Text)
		hand(minute/60, rad*0.75, t.Space(0.75), t.Text)
		hand(sec/60, rad*0.85, t.Space(0.25), t.Accent)
		var hub ui.Path
		p.FillPath(hub.Circle(cx, cy, t.Space(1)), t.Accent)
	})
}
