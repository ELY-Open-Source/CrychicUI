package finance

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AccountCardOptions word the face; Mask takes the last four digits.
type AccountCardOptions struct{ Mask, Holder, Expires, Network string }

// AccountCard is a bank card face with a masked number.
func AccountCard(c *ui.Context, bank, last4 string, o AccountCardOptions) ui.Element {
	if bank == "" || len(last4) != 4 {
		panic(fmt.Sprintf("finance: AccountCard %q needs four digits, got %q", bank, last4))
	}
	for _, r := range last4 {
		if r < '0' || r > '9' {
			panic(fmt.Sprintf("finance: AccountCard %q needs four digits, got %q", bank, last4))
		}
	}
	def(&o.Mask, "•••• %s")
	t := c.Theme()
	root := ui.Column(c).FillWidth().MaxWidth(t.Space(90)).MinHeight(t.Space(50)).Justify(ui.SpaceBetween).Gap(t.Space(3)).
		Padding(t.Space(5)).Radius(t.Radius).Background(t.Accent).TextColor(t.AccentText).Shadow(0, 1, 3, 0, kit.Shade(0.12)).
		Role(ui.RoleGroup).Label(bank)
	kit.Arrive(root, ui.Pose{Y: t.Space(3), Rotate: -2, Fade: 1}, 0)
	sweep := root.Keyframes("crychic.sheen", []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		ring := t.AccentText.Alpha(0.1)
		big, small := r.H*1.2, r.H*0.7
		p.Fill(ui.Rect{X: r.X + r.W - big*0.6*sweep, Y: r.Y - big*0.5, W: big, H: big}, ring, 999)
		p.Fill(ui.Rect{X: r.X + r.W - small*0.8, Y: r.Y + r.H - small*0.5, W: small, H: small}, ring, 999)
	}))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, bank).Bold().Grow(1).MinWidth(0).SingleLine()
			if o.Network != "" {
				ui.Text(c, o.Network).Bold().FontSize(t.Rem(0.9)).Italic().Shrink(0)
			}
		})
		ui.Box(c).Size(t.Space(10), t.Space(7.5)).Radius(t.Space(1.5)).Background(t.Warning.Mix(t.AccentText, 0.35)).
			Border(1, t.AccentText.Alpha(0.35))
		ui.Text(c, fmt.Sprintf(o.Mask, last4)).Font("monospace").FontSize(t.Rem(1.15)).Shrink(0)
		if o.Holder != "" || o.Expires != "" {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, o.Holder).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
				kit.Numeric(c, o.Expires).FontSize(t.Rem(0.8)).Shrink(0)
			})
		}
	})
	return root
}
