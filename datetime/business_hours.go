package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Hours hold one day's opening and closing times.
type Hours struct {
	Open     bool
	From, To time.Duration
}

// BusinessHoursOptions tune a BusinessHours.
type BusinessHoursOptions struct {
	// Closed shows on closed days; From and To name inputs.
	Closed, From, To string
}

// BusinessHours sets open and close times for each weekday.
func BusinessHours(c *ui.Context, label string, week *[7]Hours, o BusinessHoursOptions) ui.Element {
	if label == "" || week == nil {
		panic("datetime: BusinessHours needs a label and a week")
	}
	for d, h := range week {
		if h.From < 0 || h.From >= 24*time.Hour || h.To < 0 || h.To >= 24*time.Hour {
			panic(fmt.Sprintf("datetime: BusinessHours day %d needs times within a day", d))
		}
	}
	def(&o.Closed, "Closed")
	def(&o.From, "%s opens")
	def(&o.To, "%s closes")
	t := c.Theme()
	names := i18n.Weekdays(locale(c), i18n.Wide, true)
	first := int(c.Region().FirstWeekday)
	base := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Label(label)
	edit := func(d int, field string) {
		root.SetChanged()
		kit.Log("BusinessHours", "edited", "day", d, "field", field)
	}
	return root.Children(func() {
		for i := range 7 {
			d := (first + i) % 7
			h := &week[d]
			c.Key(d)
			ui.Row(c).FillWidth().MinHeight(t.Space(10)).Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Width(t.Space(36)).Shrink(0).Children(func() {
					if choice.Toggle(c, names[d], &h.Open).Changed() {
						edit(d, "open")
					}
				})
				if !h.Open {
					kit.Muted(c, o.Closed).Key("closed").Transition(kit.Enter)
					return
				}
				ui.Row(c).Key("hours").Gap(t.Space(2)).AlignItems(ui.Center).Transition(kit.Enter).Children(func() {
					from, to := base.Add(h.From), base.Add(h.To)
					if kit.Changed(ui.TimeInput(c, &from).Label(fmt.Sprintf(o.From, names[d])), sound.Tick) {
						h.From = from.Sub(base)
						edit(d, "from")
					}
					kit.Muted(c, "–")
					if kit.Changed(ui.TimeInput(c, &to).Label(fmt.Sprintf(o.To, names[d])), sound.Tick) {
						h.To = to.Sub(base)
						edit(d, "to")
					}
				})
			})
		}
	})
}
