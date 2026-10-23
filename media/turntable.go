package media

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TurntableOptions tune a Turntable.
type TurntableOptions struct {
	// Size is the record's width; zero gives the default.
	Size float32
}

// Turntable spins a record while playing; the arm swings on.
func Turntable(c *ui.Context, label string, art *ui.Bitmap, playing bool, o TurntableOptions) ui.Element {
	if label == "" {
		panic(fmt.Sprintf("media: Turntable %q needs a label", label))
	}
	t := c.Theme()
	size("Turntable", label, &o.Size, t.Space(48))
	d, pad := o.Size, t.Space(4)
	ink := t.Background
	if t.Dark {
		ink = t.Text
	}
	root := ui.Box(c).Size(d+2*pad+t.Space(12), d+2*pad).Radius(t.Radius).Background(t.Surface).Border(1, t.Border).
		Role(ui.RoleImage).Label(label)
	return root.Children(func() {
		disc := ui.Box(c).Absolute().Left(pad).Top(pad).Size(d, d).Radius(999).Background(kit.Shade(0.88)).Center()
		disc.Rotate(360 * spin(c, disc, playing, 0.55)).Draw(kit.Inside(d/2, func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			for k := range 6 {
				var ring ui.Path
				ring.Circle(cx, cy, r.W*(0.24+0.04*float32(k)))
				p.StrokePath(&ring, 1, ink.Alpha(0.1))
			}
			s := t.Space(1.5)
			p.Fill(ui.Rect{X: cx - s/2, Y: r.Y + r.H*0.06, W: s, H: s}, ink.Alpha(0.35), s/2)
		})).Children(func() {
			cover(c, art, d*0.36).Radius(999)
		})
		ui.Box(c).Absolute().Right(pad-t.Space(1.75)).Top(pad-t.Space(1.75)).Size(t.Space(5.5), t.Space(5.5)).Radius(999).
			Background(t.Border)
		arm := ui.Box(c).Absolute().Right(pad+t.Space(1)).Top(pad).Width(t.Space(1.5)).Height(d*0.72).Radius(999).
			Background(t.TextMuted).TransformOrigin(0.5, 0)
		arm.AnimatePose(ui.Pose{Rotate: 26 * unit(playing)}, ui.Bouncy).Children(func() {
			ui.Box(c).Absolute().Bottom(0).Left(-t.Space(1)).Size(t.Space(3.5), t.Space(5)).Radius(t.Space(0.75)).Background(t.Text)
		})
	})
}
