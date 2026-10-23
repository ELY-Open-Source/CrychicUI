package schedule

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Zone is one timezone and hour offset.
type Zone struct {
	Name   string
	Offset int
}

// ZoneCompare shows the same hour across zones.
func ZoneCompare(c *ui.Context, hour int, zones []Zone) ui.Element {
	if hour < 0 || hour > 23 || len(zones) == 0 {
		panic("schedule: ZoneCompare has invalid data")
	}
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(2)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, zone := range zones {
			if zone.Name == "" || zone.Offset < -23 || zone.Offset > 23 {
				panic("schedule: ZoneCompare has an invalid zone")
			}
			c.Key(i)
			ui.Column(c).Gap(t.Space(0.5)).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border).Children(func() {
				ui.Text(c, zone.Name).FontSize(t.Rem(0.75))
				kit.Numeric(c, fmt.Sprintf("%02d:00", (hour+zone.Offset+24)%24)).Bold()
			})
		}
	})
	return root
}
