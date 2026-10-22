package datetime

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DatePicker is a field opening MyGo's calendar; it rises in.
func DatePicker(c *ui.Context, label string, date *time.Time) ui.Element {
	if label == "" || date == nil {
		panic("datetime: DatePicker needs a label and a date")
	}
	t := c.Theme()
	root := kit.Arrive(ui.Column(c).AlignSelf(ui.Start), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	return root.Children(func() {
		ui.Field(c, label, func() {
			if kit.Changed(ui.DateInput(c, date), sound.Select) {
				root.SetChanged()
				kit.Log("DatePicker", "picked")
			}
		})
	})
}
