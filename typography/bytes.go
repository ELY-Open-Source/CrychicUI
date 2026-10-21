package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Bytes is a human byte size with rolling digits.
func Bytes(c *ui.Context, n int64) ui.Element {
	if n < 0 {
		panic("typography: Bytes count must be non-negative")
	}
	units := [...]string{"B", "KB", "MB", "GB", "TB"}
	value, unit := float64(n), 0
	for value >= 1000 && unit < len(units)-1 {
		value /= 1000
		unit++
	}
	o := i18n.NumberOptions{MaxFractionDigits: 1}
	if unit == 0 {
		o.Integer = true
	}
	text := i18n.FormatNumber(formatLocale(c), value, o)
	return kit.Numeric(c, text+" "+units[unit])
}
