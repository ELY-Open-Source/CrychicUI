package typography

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// formatLocale returns the viewer's regional format.
func formatLocale(c *ui.Context) string { return c.Region().FormatLocale(c.Locale()) }

// finiteNumber rejects NaN and infinities.
func finiteNumber(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

// Number is localized text with rolling digits.
func Number(c *ui.Context, value float64) ui.Element {
	if !finiteNumber(value) {
		panic("typography: Number value must be finite")
	}
	return kit.Numeric(c, i18n.FormatNumber(formatLocale(c), value, i18n.NumberOptions{}))
}
