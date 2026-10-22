package datetime

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DurationPickerOptions tune a DurationPicker.
type DurationPickerOptions struct {
	// MaxHours caps hours, 23 when zero; Step steps minutes, 5.
	MaxHours, Step int
	// Hours and Minutes name the steppers; the formats write them.
	Hours, Minutes, HourFormat, MinuteFormat string
}

// DurationPicker sets hours and minutes with two steppers.
func DurationPicker(c *ui.Context, label string, d *time.Duration, o DurationPickerOptions) ui.Element {
	if o.MaxHours == 0 {
		o.MaxHours = 23
	}
	if o.Step == 0 {
		o.Step = 5
	}
	def(&o.Hours, "Hours")
	def(&o.Minutes, "Minutes")
	def(&o.HourFormat, "%.0f h")
	def(&o.MinuteFormat, "%.0f min")
	step := time.Duration(o.Step) * time.Minute
	if label == "" || d == nil || o.MaxHours < 0 || o.Step < 1 || 60%o.Step != 0 || *d < 0 ||
		*d > time.Duration(o.MaxHours)*time.Hour+time.Hour-step || *d%step != 0 {
		panic(fmt.Sprintf("datetime: DurationPicker %q needs a duration on a %v step", label, step))
	}
	t := c.Theme()
	h, m := float64(*d/time.Hour), float64(*d%time.Hour/time.Minute)
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(3)).Label(label)
	return root.Children(func() {
		a := choice.NumberStepper(c, o.Hours, &h, 0, float64(o.MaxHours), 1, choice.NumberStepperOptions{Format: o.HourFormat}).Changed()
		b := choice.NumberStepper(c, o.Minutes, &m, 0, float64(60-o.Step), float64(o.Step),
			choice.NumberStepperOptions{Format: o.MinuteFormat}).Changed()
		if a || b {
			*d = time.Duration(h)*time.Hour + time.Duration(m)*time.Minute
			root.SetChanged()
			kit.Log("DurationPicker", "changed", "minutes", int(d.Minutes()))
		}
	})
}
