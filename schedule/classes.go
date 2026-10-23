package schedule

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ClassSession is one weekly class.
type ClassSession struct {
	Name, Day, Start, End, Room string
}

// Classes shows a weekly class timetable.
func Classes(c *ui.Context, classes []ClassSession) ui.Element {
	if len(classes) == 0 {
		panic("schedule: Classes needs classes")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, class := range classes {
			if class.Name == "" || class.Day == "" || class.Start == "" || class.End == "" || class.Room == "" {
				panic("schedule: Classes has an invalid class")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				ui.Text(c, class.Name).Grow(1)
				kit.Muted(c, class.Day+" "+class.Start+"–"+class.End+" "+class.Room).FontSize(t.Rem(0.75))
			})
		}
	})
	return root
}
