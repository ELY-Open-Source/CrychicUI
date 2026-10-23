package finance

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AccountCardOptions are the account labels.
type AccountCardOptions struct{ Mask string }

// AccountCard shows a bank card with masked number.
func AccountCard(c *ui.Context, bank, last4 string, o AccountCardOptions) ui.Element {
	if bank == "" || len(last4) != 4 {
		panic("finance: AccountCard needs bank and four digits")
	}
	for _, r := range last4 {
		if r < '0' || r > '9' {
			panic("finance: AccountCard needs four digits")
		}
	}
	def(&o.Mask, "•••• %s")
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Width(t.Space(64)).Gap(t.Space(4)).Padding(t.Space(4)).Radius(t.Radius).
		Background(t.Accent).TextColor(t.AccentText)
	kit.Arrive(root, ui.Pose{Y: t.Space(3), Fade: 1}, 0)
	root.Children(func() {
		ui.Text(c, bank).Bold()
		ui.Text(c, fmt.Sprintf(o.Mask, last4)).Font("monospace").FontSize(t.Rem(1.1))
	})
	return root
}
