package schedule

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Calendar is one calendar and color index.
type Calendar struct {
	Name string
	Tint int
}

// CalendarList toggles calendar colors; cue Toggle.
func CalendarList(c *ui.Context, calendars []Calendar, visible *[]bool) ui.Element {
	if len(calendars) == 0 || visible == nil || len(*visible) != len(calendars) {
		panic("schedule: CalendarList needs calendars and visibility")
	}
	t := c.Theme()
	inks := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, calendar := range calendars {
			if calendar.Name == "" || calendar.Tint < 0 || calendar.Tint >= len(inks) {
				panic("schedule: CalendarList has an invalid calendar")
			}
			c.Key(i)
			b := kit.Press(ui.ButtonBase(c).Label(calendar.Name).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Checked((*visible)[i]))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.1), unit((*visible)[i])), ui.Snappy))
			b.Children(func() {
				ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(inks[calendar.Tint])
				ui.Text(c, calendar.Name).Grow(1)
			})
			if b.Clicked() {
				(*visible)[i] = !(*visible)[i]
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("CalendarList", "toggled", "index", i, "visible", (*visible)[i])
			}
		}
	})
	return root
}
