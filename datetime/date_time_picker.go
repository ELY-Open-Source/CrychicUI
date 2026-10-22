package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DateTimePickerOptions tune a DateTimePicker.
type DateTimePickerOptions struct {
	// Date and Time name the two inputs.
	Date, Time string
}

// DateTimePicker edits one moment's date and time; it rises in.
func DateTimePicker(c *ui.Context, label string, at *time.Time, o DateTimePickerOptions) ui.Element {
	if label == "" || at == nil {
		panic("datetime: DateTimePicker needs a label and a time")
	}
	def(&o.Date, "Date")
	def(&o.Time, "Time")
	t := c.Theme()
	root := kit.Arrive(ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(1)), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			day := ui.DateInput(c, at).Label(o.Date)
			clock := ui.TimeInput(c, at).Label(o.Time)
			if kit.Changed(day, sound.Select) || kit.Changed(clock, sound.Select) {
				root.SetChanged()
				kit.Log("DateTimePicker", "changed")
			}
		})
	})
}
