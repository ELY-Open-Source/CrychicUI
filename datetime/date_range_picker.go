package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// DateRangePickerOptions tune a DateRangePicker.
type DateRangePickerOptions struct {
	// Calendar names its month buttons.
	Calendar RangeCalendarOptions
}

// DateRangePicker is a field opening a range calendar popover.
func DateRangePicker(c *ui.Context, label string, start, end *time.Time, o DateRangePickerOptions) ui.Element {
	if label == "" || start == nil || end == nil {
		panic("datetime: DateRangePicker needs a label and dates")
	}
	t := c.Theme()
	loc := locale(c)
	shown := i18n.FormatDate(loc, *start, i18n.DateOptions{Skeleton: "yMMMd"}) + " – " +
		i18n.FormatDate(loc, *end, i18n.DateOptions{Skeleton: "yMMMd"})
	field := kit.Press(ui.ButtonBase(c).Label(label).Gap(t.Space(2)).Padding(t.Space(2), t.Space(3)).Radius(t.Radius))
	open := ui.Local(field, "open", func() bool { return false })
	border := t.Border
	if *open || field.Hovered() {
		border = t.Accent
	}
	field.Background(t.Background).Border(1, field.AnimateColor("border", border, ui.Snappy)).Expanded(*open).
		Children(func() {
			ui.Icon(c, icon.Must("calendar")).Size(16, 16).TextColor(t.TextMuted)
			ui.Text(c, shown).ContentTransition(ui.ContentTransition{})
		})
	if field.Clicked() {
		if *open = !*open; *open {
			sound.Play(sound.Open)
			kit.Log("DateRangePicker", "open")
		}
	}
	ui.Popover(c, field, open, func() {
		if RangeCalendar(c, start, end, o.Calendar).Changed() {
			field.SetChanged()
		}
	})
	return field
}
