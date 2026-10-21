package display

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Swatch is a color sample with its hex value.
func Swatch(c *ui.Context, color ui.Color) ui.Element {
	if color.A == 0 {
		panic("display: Swatch needs a color")
	}
	label := fmt.Sprintf("#%02X%02X%02X", color.R, color.G, color.B)
	if color.A != 255 {
		label = fmt.Sprintf("#%02X%02X%02X%02X", color.R, color.G, color.B, color.A)
	}
	t := c.Theme()
	return ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		box := ui.Box(c).Size(t.Space(8), t.Space(8)).Radius(t.Radius).Border(1, t.Border)
		box.Background(box.AnimateColor("color", color, ui.Snappy))
		kit.Muted(c, label).Font("monospace").FontSize(t.Rem(0.85))
	})
}
