package typography

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Money shows a currency code and two decimals.
func Money(c *ui.Context, amount float64, currency string) ui.Element {
	if !finiteNumber(amount) || currency == "" {
		panic("typography: Money needs a finite amount and currency")
	}
	text := i18n.FormatNumber(formatLocale(c), amount, i18n.NumberOptions{MinFractionDigits: 2, MaxFractionDigits: 2})
	return kit.Numeric(c, currency+" "+text)
}
