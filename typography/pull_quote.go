package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PullQuote is large centered quoted text fading up.
func PullQuote(c *ui.Context, quote, attribution string) ui.Element {
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
		ui.Text(c, quote).FontSize(t.Rem(1.5)).TextAlign(ui.Center)
		kit.Muted(c, attribution)
	})
	return kit.Arrive(root, ui.Pose{Y: 8, Fade: 1}, 0)
}
