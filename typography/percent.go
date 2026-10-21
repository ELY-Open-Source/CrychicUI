package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Percent is localized percent text with sign coloring.
func Percent(c *ui.Context, value float64) ui.Element {
	if !finiteNumber(value) {
		panic("typography: Percent value must be finite")
	}
	t := c.Theme()
	color := t.Text
	if value > 0 {
		color = t.Success
	} else if value < 0 {
		color = t.Danger
	}
	text := i18n.FormatNumber(formatLocale(c), value, i18n.NumberOptions{Percent: true})
	return kit.Numeric(c, text).TextColor(color)
}
