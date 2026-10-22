package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

type weekThumb struct{ id uint64 }

// WeekStrip picks a day in its week; the mark slides.
func WeekStrip(c *ui.Context, label string, day *time.Time) ui.Element {
	if label == "" || day == nil || day.IsZero() {
		panic("datetime: WeekStrip needs a label and a day")
	}
	t := c.Theme()
	loc := locale(c)
	start := day.AddDate(0, 0, -((int(day.Weekday()) - int(c.Region().FirstWeekday) + 7) % 7))
	sel := days(start, *day)
	shown := sel
	names := i18n.Weekdays(loc, i18n.Narrow, true)
	today := c.Now()
	root := ui.Column(c).AlignSelf(ui.Start)
	id, changed := root.ID(), false
	root.Children(func() {
		ui.RadioGroup(c, func() {
			for i := range 7 {
				d := start.AddDate(0, 0, i)
				c.Key(i)
				b := kit.Press(ui.RadioBase(c, &sel, i).Label(i18n.FormatDate(loc, d, i18n.DateOptions{Skeleton: "yMMMEd"})).
					Radius(t.Radius))
				bg, fg := t.Accent.Alpha(0), t.Text
				if i == shown {
					fg = t.AccentText
				} else if b.Hovered() {
					bg = t.SurfaceHover
				}
				b.Background(b.AnimateColor("bg", bg, ui.Snappy)).TextColor(b.AnimateColor("fg", fg, ui.Snappy))
				// One cell stacks the mark under the day.
				b.Children(func() {
					ui.Grid(c).Columns(1).Children(func() {
						if i == shown {
							ui.Box(c).Key("mark").ColumnStart(1).RowStart(1).Radius(t.Radius).Background(t.Accent).
								Transition(ui.ElementTransition{Animation: &ui.Snappy, Match: weekThumb{id}})
						}
						ui.Column(c).Key("face").ColumnStart(1).RowStart(1).Size(t.Space(10), t.Space(14)).
							AlignItems(ui.Center).Justify(ui.Center).Gap(t.Space(0.5)).Children(func() {
							ui.Text(c, names[d.Weekday()]).FontSize(t.Rem(0.75))
							kit.Numeric(c, i18n.FormatDate(loc, d, i18n.DateOptions{Skeleton: "d"})).Bold()
							dot := t.Accent.Alpha(0)
							if sameDay(d, today) {
								dot = t.Accent
								if i == shown {
									dot = t.AccentText
								}
							}
							ui.Box(c).Size(t.Space(1), t.Space(1)).Radius(999).Background(dot)
						})
					})
				})
				changed = kit.Changed(b, sound.Select) || changed
			}
		}).Label(label).Row().Gap(t.Space(1))
	})
	if changed {
		*day = onDay(start.AddDate(0, 0, sel), *day)
		root.SetChanged()
		kit.Log("WeekStrip", "picked", "index", sel)
	}
	return root
}
