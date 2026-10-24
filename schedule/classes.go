package schedule

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ClassSession is one weekly class.
type ClassSession struct {
	Name, Day, Start, End, Room string
}

// ClassesOptions name the timetable's days.
type ClassesOptions struct{ Days []string }

// Classes shows a weekly class timetable, one column per day.
func Classes(c *ui.Context, classes []ClassSession, o ClassesOptions) ui.Element {
	if len(classes) == 0 {
		panic("schedule: Classes needs classes")
	}
	if len(o.Days) == 0 {
		o.Days = []string{"Mon", "Tue", "Wed", "Thu", "Fri"}
	}
	t := c.Theme()
	tint := map[string]int{}
	for _, class := range classes {
		if class.Name == "" || class.Day == "" || class.Start == "" || class.End == "" || class.Room == "" {
			panic(fmt.Sprintf("schedule: Classes has an invalid class %q", class.Name))
		}
		found := false
		for _, day := range o.Days {
			found = found || day == class.Day
		}
		if !found {
			panic(fmt.Sprintf("schedule: Classes class %q is on unknown day %q", class.Name, class.Day))
		}
		if _, ok := tint[class.Name]; !ok {
			tint[class.Name] = len(tint) % 4
		}
	}
	root := ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Start).Role(ui.RoleGroup)
	root.Children(func() {
		n := 0
		for d, day := range o.Days {
			c.Key(d)
			ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(1.5)).Children(func() {
				ui.Text(c, day).FontSize(t.Rem(0.75)).Bold().Padding(0, 0, t.Space(1), 0).
					BorderWidth(0, 0, 1, 0).BorderColor(t.Border).FillWidth()
				for i, class := range classes {
					if class.Day != day {
						continue
					}
					ink := inks(t)[tint[class.Name]]
					c.Key(i)
					cell := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5)).Radius(t.Space(1.5))
					cell.Background(cell.AnimateColor("bg", ink.Alpha(0.12).Mix(ink.Alpha(0.22), unit(cell.Hovered())), ui.Snappy))
					cell.Draw(func(p *ui.Painter, r ui.Rect) { p.Fill(ui.Rect{X: r.X, Y: r.Y, W: t.Space(0.75), H: r.H}, ink, 0) })
					kit.Stagger(c, root, cell, n).Clip().TextColor(ink.Mix(t.Text, 0.4)).Children(func() {
						kit.Numeric(c, class.Start+"–"+class.End).FontSize(t.Rem(0.6))
						ui.Text(c, class.Name).FontSize(t.Rem(0.75)).Bold().Wrap()
						ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
							glyph(c, "map-pin", 10)
							ui.Text(c, class.Room).FontSize(t.Rem(0.6)).SingleLine().MinWidth(0)
						})
					})
					n++
				}
			})
		}
	})
	return root
}
