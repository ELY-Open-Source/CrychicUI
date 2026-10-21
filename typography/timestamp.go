package typography

import (
	"time"

	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Timestamp is localized date and time with a full tooltip.
func Timestamp(c *ui.Context, when time.Time) ui.Element {
	locale := formatLocale(c)
	text := i18n.FormatDate(locale, when, i18n.DateOptions{Date: i18n.MediumStyle, Time: i18n.ShortStyle})
	full := i18n.FormatDate(locale, when, i18n.DateOptions{Date: i18n.FullStyle, Time: i18n.LongStyle})
	return ui.Text(c, text).Tooltip(full).ContentTransition(ui.ContentTransition{})
}
