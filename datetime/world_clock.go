package datetime

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Zone is a named place and its location.
type Zone struct {
	Name     string
	Location *time.Location
}

// WorldClockOptions tune a WorldClock.
type WorldClockOptions struct {
	// Same, Ahead and Behind describe offsets; two take hours.
	Same, Ahead, Behind string
}

// WorldClock lists zones with their local times and offsets.
func WorldClock(c *ui.Context, zones []Zone, o WorldClockOptions) ui.Element {
	if len(zones) == 0 {
		panic("datetime: WorldClock needs zones")
	}
	for _, z := range zones {
		if z.Name == "" || z.Location == nil {
			panic(fmt.Sprintf("datetime: WorldClock zone %q needs a name and a location", z.Name))
		}
	}
	def(&o.Same, "Same time")
	def(&o.Ahead, "%s h ahead")
	def(&o.Behind, "%s h behind")
	t := c.Theme()
	now := c.Now()
	c.After(time.Minute - time.Duration(now.Second())*time.Second - time.Duration(now.Nanosecond()))
	_, here := now.Zone()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		for i, z := range zones {
			at := now.In(z.Location)
			_, there := at.Zone()
			hours := strconv.FormatFloat(float64(there-here)/3600, 'f', -1, 64)
			offset := o.Same
			if there > here {
				offset = fmt.Sprintf(o.Ahead, hours)
			} else if there < here {
				offset = fmt.Sprintf(o.Behind, hours[1:])
			}
			sky := "moon"
			if at.Hour() >= 6 && at.Hour() < 18 {
				sky = "sun"
			}
			row := ui.Row(c).Key(z.Name).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), 0)
			kit.Arrive(row, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(i)*40*time.Millisecond).Children(func() {
				ui.Icon(c, icon.Must(sky)).Size(16, 16).TextColor(t.TextMuted).ContentTransition(ui.ContentTransition{})
				ui.Column(c).Grow(1).Children(func() {
					ui.Text(c, z.Name).FontWeight(500)
					kit.Muted(c, offset)
				})
				kit.Numeric(c, i18n.FormatDate(locale(c), at, i18n.DateOptions{Skeleton: "jm"})).FontWeight(600)
			})
		}
	})
}
