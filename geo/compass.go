package geo

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CompassOptions word the compass; Points run from N clockwise.
type CompassOptions struct {
	Label, Degrees string
	Points         [8]string
}

// Compass is a dial; its needle turns the short way.
func Compass(c *ui.Context, heading float64, o CompassOptions) ui.Element {
	def(&o.Label, "Heading")
	def(&o.Degrees, "%d°")
	for i, d := range [8]string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"} {
		def(&o.Points[i], d)
	}
	if math.IsNaN(heading) || math.IsInf(heading, 0) {
		panic(fmt.Sprintf("geo: Compass %q needs a finite heading, got %v", o.Label, heading))
	}
	h := math.Mod(math.Mod(heading, 360)+360, 360)
	t := c.Theme()
	root := ui.Column(c).Shrink(0).Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup).Label(o.Label)
	kit.Arrive(root, ui.Pose{Scale: 0.92, Fade: 1}, 0)
	// shown is unwrapped, so 350 to 10 turns 20 degrees.
	shown := ui.Local(root, "shown", func() float64 { return h })
	d := math.Mod(h-*shown, 360)
	if d > 180 {
		d -= 360
	} else if d < -180 {
		d += 360
	}
	*shown += d
	angle := root.AnimateTo("needle", float32(*shown), ui.Smooth)
	size := t.Space(40)
	return root.Children(func() {
		dial := ui.Box(c).Size(size, size).Shrink(0).Radius(999).Background(t.Background).Border(1, t.Border).Shadow(0, 1, 3, 0, kit.Shade(0.08))
		dial.Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, r.W/2
			for i := range 36 {
				a := float64(i) * 10 * math.Pi / 180
				in := rad - t.Space(1.5)
				if i%9 == 0 {
					in = rad - t.Space(3)
				}
				var tick ui.Path
				tick.MoveTo(cx+in*float32(math.Sin(a)), cy-in*float32(math.Cos(a))).
					LineTo(cx+(rad-t.Space(0.5))*float32(math.Sin(a)), cy-(rad-t.Space(0.5))*float32(math.Cos(a)))
				p.StrokePath(&tick, 1, t.TextMuted.Alpha(0.5))
			}
			a := float64(angle) * math.Pi / 180
			sin, cos := float32(math.Sin(a)), float32(math.Cos(a))
			at := func(along, across float32) (float32, float32) {
				return cx + along*sin + across*cos, cy - along*cos + across*sin
			}
			needle := func(tip float32, ink ui.Color) {
				var half ui.Path
				x0, y0 := at(tip, 0)
				x1, y1 := at(0, t.Space(2))
				x2, y2 := at(0, -t.Space(2))
				half.MoveTo(x0, y0).LineTo(x1, y1).LineTo(x2, y2).Close()
				p.FillPath(&half, ink)
			}
			needle(rad*0.62, t.Danger)
			needle(-rad*0.62, t.TextMuted.Alpha(0.6))
			var hub ui.Path
			p.FillPath(hub.Circle(cx, cy, t.Space(1.5)), t.Text)
		})
		dial.Children(func() {
			for i, xy := range [4][2]float32{{0.5, 0.2}, {0.8, 0.5}, {0.5, 0.8}, {0.2, 0.5}} {
				w := t.Space(6)
				ink := t.TextMuted
				if i == 0 {
					ink = t.Danger
				}
				ui.Box(c).Absolute().Left(size*xy[0]-w/2).Top(size*xy[1]-w/2).Size(w, w).Center().Children(func() {
					ui.Text(c, o.Points[i*2]).FontSize(t.Rem(0.74)).Bold().TextColor(ink).Shrink(0).SingleLine()
				})
			}
		})
		ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Degrees, int(math.Round(h))%360)).FontSize(t.Rem(1.3)).Bold().Shrink(0).MaxWidth(t.Space(30)).SingleLine()
			ui.Text(c, o.Points[int(math.Round(h/45))%8]).FontSize(t.Rem(0.9)).Bold().TextColor(t.Danger).Shrink(0).MaxWidth(t.Space(16)).SingleLine().
				ContentTransition(ui.ContentTransition{})
		})
	})
}
