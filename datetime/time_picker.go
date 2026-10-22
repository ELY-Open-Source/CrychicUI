package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// TimePickerOptions tune a TimePicker.
type TimePickerOptions struct {
	// Step is minutes between choices, dividing 60; 5 when zero.
	Step int
	// Hours, Minutes and Period name the wheels.
	Hours, Minutes, Period string
}

// TimePicker turns hour and minute wheels in the region's cycle.
func TimePicker(c *ui.Context, label string, at *time.Time, o TimePickerOptions) ui.Element {
	if o.Step == 0 {
		o.Step = 5
	}
	if label == "" || at == nil || o.Step < 1 || 60%o.Step != 0 {
		panic(fmt.Sprintf("datetime: TimePicker needs a label, a time and a step dividing 60, not %d", o.Step))
	}
	def(&o.Hours, "Hours")
	def(&o.Minutes, "Minutes")
	def(&o.Period, "AM or PM")
	t := c.Theme()
	h12 := twelve(c)
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		hour, minute, period := at.Hour(), at.Minute()/o.Step, at.Hour()/12
		var hours []string
		if h12 {
			hour %= 12
			for h := range 12 {
				hours = append(hours, fmt.Sprint((h+11)%12+1))
			}
		} else {
			for h := range 24 {
				hours = append(hours, fmt.Sprintf("%02d", h))
			}
		}
		var mins []string
		for m := 0; m < 60; m += o.Step {
			mins = append(mins, fmt.Sprintf("%02d", m))
		}
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			moved := choice.Wheel(c, o.Hours, hours, &hour).Changed()
			ui.Text(c, ":").Bold()
			moved = choice.Wheel(c, o.Minutes, mins, &minute).Changed() || moved
			if h12 {
				am, pm := i18n.DayPeriods(locale(c), i18n.Abbreviated)
				moved = choice.Wheel(c, o.Period, []string{am, pm}, &period).Changed() || moved
				hour += 12 * period
			}
			if moved {
				*at = time.Date(at.Year(), at.Month(), at.Day(), hour, minute*o.Step, 0, 0, at.Location())
				root.SetChanged()
				kit.Log("TimePicker", "changed")
			}
		})
	})
}
