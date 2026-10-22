package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PaywallOptions tune a Paywall.
type PaywallOptions struct {
	Upgrade string
}

// Paywall blurs content under an upgrade card.
func Paywall(c *ui.Context, title, text string, o PaywallOptions, content func()) (root, upgrade ui.Element) {
	if title == "" || content == nil {
		panic(fmt.Sprintf("overlay: Paywall %q needs content", title))
	}
	def(&o.Upgrade, "Upgrade")
	t := c.Theme()
	root = ui.Box(c)
	cover(c, root, true, t.Space(1), content, func() {
		card(c, title, text, "sparkles", func() {
			upgrade = primary(c, o.Upgrade, t.Accent, t.AccentText)
			if upgrade.Clicked() {
				sound.Play(sound.Open)
				kit.Log("Paywall", "upgrade")
			}
		})
	})
	return root, upgrade
}
