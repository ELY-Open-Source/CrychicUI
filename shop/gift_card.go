package shop

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// GiftCardOptions word the face; Code shows only its last four.
type GiftCardOptions struct{ Title, From, Code string }

// GiftCard is a card face: amount, message, sender, code.
func GiftCard(c *ui.Context, amount float64, currency, message string, o GiftCardOptions) ui.Element {
	if amount < 0 || !finite(amount) || currency == "" || message == "" {
		panic(fmt.Sprintf("shop: GiftCard needs an amount, a currency and a message, got %v", amount))
	}
	def(&o.Title, "Gift card")
	t := c.Theme()
	root := ui.Column(c).FillWidth().MinHeight(t.Space(50)).Justify(ui.SpaceBetween).Gap(t.Space(3)).Padding(t.Space(4)).
		Radius(t.Radius).Background(t.Accent).TextColor(t.AccentText).Shadow(0, 1, 3, 0, kit.Shade(0.12)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Rotate: -3, Scale: 0.94, Fade: 1}, 0)
	glow := root.Keyframes("crychic.sheen", []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		ring := t.AccentText.Alpha(0.1)
		big, small := r.H*1.1, r.H*0.6
		p.Fill(ui.Rect{X: r.X + r.W - big*0.55*glow, Y: r.Y - big*0.45, W: big, H: big}, ring, 999)
		p.Fill(ui.Rect{X: r.X + r.W - small*0.9, Y: r.Y + r.H - small*0.55, W: small, H: small}, ring, 999)
	}))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon.Must("gift")).Size(18, 18)
			ui.Text(c, o.Title).Bold().FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
		})
		ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
			money(c, amount, currency).Bold().FontSize(t.Rem(1.8))
			ui.Text(c, message).FontSize(t.Rem(0.9)).FillWidth().Wrap()
		})
		if o.From != "" || o.Code != "" {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, o.From).FontSize(t.Rem(0.8)).TextColor(t.AccentText.Alpha(0.85)).Grow(1).MinWidth(0).SingleLine()
				if n := utf8.RuneCountInString(o.Code); n > 0 {
					last := string([]rune(o.Code)[max(0, n-4):])
					kit.Numeric(c, "•••• "+last).FontSize(t.Rem(0.8)).Bold().Shrink(0)
				}
			})
		}
	})
	return root
}
