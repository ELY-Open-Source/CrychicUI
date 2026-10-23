package shop

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ThankYouOptions are the confirmation labels.
type ThankYouOptions struct{ Title, Order string }

// ThankYou confirms an order with a drawn check; cue Success.
func ThankYou(c *ui.Context, number string, o ThankYouOptions) ui.Element {
	if number == "" {
		panic("shop: ThankYou needs an order number")
	}
	def(&o.Title, "Thank you")
	def(&o.Order, "Order %s")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(4)).Radius(t.Radius).Border(1, t.Border)
	played := ui.Local(root, "played", func() bool { return false })
	if !*played {
		*played = true
		sound.Play(sound.Success)
		kit.Log("ThankYou", "confirmed")
	}
	root.Children(func() {
		check := ui.Box(c).Size(t.Space(10), t.Space(10)).Radius(999).Background(t.Success.Alpha(0.15))
		progress := check.Keyframes("check", []ui.Keyframe{{Value: 0}, {At: 400 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		check.Draw(func(p *ui.Painter, r ui.Rect) {
			var path ui.Path
			path.MoveTo(r.X+r.W*0.25, r.Y+r.H*0.52).LineTo(r.X+r.W*0.43, r.Y+r.H*0.7).LineTo(r.X+r.W*0.75, r.Y+r.H*0.3)
			p.StrokePath(path.Trim(0, progress), t.Space(0.75), t.Success)
		})
		ui.Text(c, o.Title).Bold().FontSize(t.Rem(1.4))
		kit.Muted(c, fmt.Sprintf(o.Order, number))
	})
	return root
}
