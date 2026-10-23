package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// GiftCard shows a card face with amount and message.
func GiftCard(c *ui.Context, amount float64, currency, message string) ui.Element {
	if amount < 0 || !finite(amount) || currency == "" || message == "" {
		panic("shop: GiftCard has invalid data")
	}
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Width(t.Space(64)).Gap(t.Space(3)).Padding(t.Space(4)).Radius(t.Radius).
		Background(t.Accent).TextColor(t.AccentText)
	kit.Arrive(root, ui.Pose{Rotate: -2, Scale: 0.95, Fade: 1}, 0)
	root.Children(func() {
		money(c, amount, currency).Bold().FontSize(t.Rem(1.5))
		ui.Text(c, message).Wrap()
	})
	return root
}
