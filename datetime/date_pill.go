package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// DatePill shows a weekday and day; changes cross-fade.
func DatePill(c *ui.Context, day time.Time) ui.Element {
	if day.IsZero() {
		panic("datetime: DatePill needs a day")
	}
	t := c.Theme()
	loc := locale(c)
	root := ui.Row(c).AlignSelf(ui.Start).AlignItems(ui.Center).Gap(t.Space(1.5)).Padding(t.Space(1), t.Space(3)).
		Radius(999).Background(t.Surface).Border(1, t.Border).Role(ui.RoleText).
		Label(i18n.FormatDate(loc, day, i18n.DateOptions{Skeleton: "yMMMEd"}))
	return root.Children(func() {
		ui.Text(c, i18n.Weekdays(loc, i18n.Abbreviated, true)[day.Weekday()]).TextColor(t.TextMuted).
			ContentTransition(ui.ContentTransition{})
		kit.Numeric(c, i18n.FormatDate(loc, day, i18n.DateOptions{Skeleton: "d"})).Bold().
			ContentTransition(ui.ContentTransition{})
	})
}
