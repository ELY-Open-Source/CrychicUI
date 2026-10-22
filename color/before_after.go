package color

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BeforeAfterOptions tune a BeforeAfter.
type BeforeAfterOptions struct {
	// Before and After caption the two swatches.
	Before, After string
}

// BeforeAfter sets two swatches apart with an arrow.
func BeforeAfter(c *ui.Context, before, after ui.Color, o BeforeAfterOptions) ui.Element {
	def(&o.Before, "Before")
	def(&o.After, "After")
	t := c.Theme()
	side := func(caption string, col ui.Color) {
		ui.Column(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			swatch(c, col, t.Space(16), t.Space(16))
			ui.Text(c, caption).FontSize(t.Rem(0.85)).Bold()
			kit.Muted(c, hexOf(col)).FontSize(t.Rem(0.75)).ContentTransition(ui.ContentTransition{})
		})
	}
	return ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
		side(o.Before, before)
		arrow := ui.Icon(c, icon.Must("arrow-right")).Size(16, 16).TextColor(t.TextMuted)
		kit.Arrive(arrow, ui.Pose{X: -t.Space(2), Fade: 1}, 150*time.Millisecond)
		side(o.After, after)
	})
}
