package schedule

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Calendar is one calendar and color index.
type Calendar struct {
	Name string
	Tint int
}

// CalendarListOptions word the list; Shown takes shown and total.
type CalendarListOptions struct{ Title, Shown string }

// CalendarList toggles calendar colors; cue Toggle.
func CalendarList(c *ui.Context, calendars []Calendar, visible *[]bool, o CalendarListOptions) ui.Element {
	if len(calendars) == 0 || visible == nil || len(*visible) != len(calendars) {
		panic("schedule: CalendarList needs calendars and visibility")
	}
	def(&o.Title, "My calendars")
	def(&o.Shown, "%d of %d")
	t := c.Theme()
	shown := 0
	seen := map[string]bool{}
	for i, calendar := range calendars {
		if calendar.Name == "" || calendar.Tint < 0 || calendar.Tint > 3 || seen[calendar.Name] {
			panic(fmt.Sprintf("schedule: CalendarList has an invalid or repeated calendar %q", calendar.Name))
		}
		seen[calendar.Name] = true
		if (*visible)[i] {
			shown++
		}
	}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "calendars", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Shown, shown, len(calendars)), t.TextMuted) })
		list := ui.Column(c).FillWidth().Padding(t.Space(1))
		list.Children(func() {
			for i, calendar := range calendars {
				ink, on := inks(t)[calendar.Tint], (*visible)[i]
				c.Key(calendar.Name)
				b := kit.Press(ui.ButtonBase(c).Label(calendar.Name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.5), t.Space(2)).
					Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(0.7*unit(b.Hovered())), ui.Snappy))
				kit.Stagger(c, list, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
						box := ui.Box(c).Size(t.Space(4.5), t.Space(4.5)).Shrink(0).Center().Radius(t.Space(1)).Border(1, ink).
							TextColor(t.AccentText)
						box.Background(box.AnimateColor("fill", ink.Alpha(unit(on)), ui.Snappy)).Children(func() {
							if on {
								c.Key("check")
								kit.Arrive(glyph(c, "check", 12), ui.Pose{Scale: 0.4, Fade: 1}, 0)
							}
						})
						ui.Text(c, calendar.Name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine().
							TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Text, unit(on)), ui.Snappy))
					})
				})
				if b.Clicked() {
					(*visible)[i] = !on
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("CalendarList", "toggled", "index", i, "visible", !on)
				}
			}
		})
	})
	return root
}
