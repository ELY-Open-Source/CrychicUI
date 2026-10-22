package form

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ScheduleFieldOptions tune a ScheduleField.
type ScheduleFieldOptions struct {
	// Date and Time name the two inputs.
	Date, Time string
}

// ScheduleField edits the date and the time of one moment.
func ScheduleField(c *ui.Context, label string, at *time.Time, o ScheduleFieldOptions) ui.Element {
	if label == "" || at == nil {
		panic("form: ScheduleField needs a label and a time")
	}
	def(&o.Date, "Date")
	def(&o.Time, "Time")
	t := c.Theme()
	root := kit.Arrive(ui.Column(c).FillWidth(), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	return root.Children(func() {
		FieldGroup(c, label, FieldGroupOptions{}, func() {
			ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				day := ui.DateInput(c, at).Label(o.Date)
				clock := ui.TimeInput(c, at).Label(o.Time)
				if kit.Changed(day, sound.Select) || kit.Changed(clock, sound.Select) {
					root.SetChanged()
					kit.Log("ScheduleField", "changed")
				}
			})
		})
	})
}
