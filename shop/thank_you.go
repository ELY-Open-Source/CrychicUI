package shop

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ThankYouOptions word the card; Order takes the number.
type ThankYouOptions struct{ Title, Order, Note string }

// ThankYou confirms an order with check and confetti; cue Success.
func ThankYou(c *ui.Context, number string, o ThankYouOptions) ui.Element {
	if number == "" {
		panic("shop: ThankYou needs an order number")
	}
	def(&o.Title, "Thank you for your order")
	def(&o.Order, "Order %s")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(5), t.Space(4)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Scale: 0.96, Fade: 1}, 0)
	played := ui.Local(root, "played", func() bool { return false })
	if !*played {
		*played = true
		sound.Play(sound.Success)
		kit.Log("ThankYou", "confirmed")
	}
	burst := root.Keyframes("crychic.burst", []ui.Keyframe{{Value: 0}, {At: 250 * time.Millisecond, Value: 0}, {At: 1100 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	inks := [...]ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	top, side := t.Space(5), t.Space(14)
	root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		if burst <= 0 || burst >= 1 {
			return
		}
		cx, cy := r.X+r.W/2, r.Y+top+side/2
		for i := range 16 {
			a := float64(i)*math.Pi/8 + 0.2
			d := side*0.6 + t.Space(18)*burst*float32(1+i%3)/3
			s := t.Space(1.5) * (1 - burst)
			x, y := cx+d*float32(math.Cos(a)), cy+d*float32(math.Sin(a))
			p.Fill(ui.Rect{X: x - s/2, Y: y - s/2, W: s, H: s}, inks[i%4].Alpha(1-burst), 999)
		}
	}))
	root.Children(func() {
		check := ui.Box(c).Size(side, side).Radius(999).Background(t.Success.Alpha(0.15))
		kit.Arrive(check, ui.Pose{Scale: 0.6, Fade: 1}, 0)
		progress := check.Keyframes("check", []ui.Keyframe{{Value: 0}, {At: 150 * time.Millisecond, Value: 0}, {At: 550 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		check.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
			var path ui.Path
			path.MoveTo(r.X+r.W*0.3, r.Y+r.H*0.52).LineTo(r.X+r.W*0.45, r.Y+r.H*0.66).LineTo(r.X+r.W*0.7, r.Y+r.H*0.36)
			p.StrokePath(path.Trim(0, progress), t.Space(1), t.Success)
		}))
		ui.Text(c, o.Title).Bold().FontSize(t.Rem(1.3)).TextAlign(ui.Center).FillWidth().Wrap().Padding(t.Space(1), 0, 0, 0)
		kit.Numeric(c, fmt.Sprintf(o.Order, number)).FontSize(t.Rem(0.85)).Bold().Padding(t.Space(1), t.Space(3)).
			Radius(999).Background(t.Surface)
		if o.Note != "" {
			kit.Muted(c, o.Note).FontSize(t.Rem(0.85)).TextAlign(ui.Center).FillWidth().Wrap()
		}
	})
	return root
}
