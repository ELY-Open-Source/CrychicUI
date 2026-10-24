package schedule

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SlotDay is one date and available count.
type SlotDay struct {
	Date      string
	Available int
}

// SlotCalendarOptions word the counts.
type SlotCalendarOptions struct{ Full, Left, One string }

// SlotCalendar picks a date with free slots; cue Select.
func SlotCalendar(c *ui.Context, days []SlotDay, picked *int, o SlotCalendarOptions) ui.Element {
	if len(days) == 0 || picked == nil || *picked < -1 || *picked >= len(days) {
		panic("schedule: SlotCalendar needs days and a pick")
	}
	def(&o.Full, "Full")
	def(&o.Left, "%d left")
	def(&o.One, "1 left")
	seen := map[string]bool{}
	for _, day := range days {
		if day.Date == "" || day.Available < 0 || seen[day.Date] {
			panic(fmt.Sprintf("schedule: SlotCalendar has an invalid or repeated day %q", day.Date))
		}
		seen[day.Date] = true
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(7).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, day := range days {
			text, ink := fmt.Sprintf(o.Left, day.Available), t.Success
			switch {
			case day.Available == 0:
				text, ink = o.Full, t.TextMuted
			case day.Available == 1:
				text, ink = o.One, t.Warning
			case day.Available == 2:
				ink = t.Warning
			}
			on := *picked == i
			c.Key(day.Date)
			b := kit.Press(ui.ButtonBase(c).Label(day.Date).Height(t.Space(15)).AlignItems(ui.Stretch).Padding(t.Space(1.5)).
				Radius(t.Space(1.5)).Checked(on).Disabled(day.Available == 0))
			b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
			b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
			kit.Stagger(c, root, b, i/7).Children(func() {
				ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
					ui.Text(c, day.Date).FontSize(t.Rem(0.7)).Bold().SingleLine().MinWidth(0)
					ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(t.Space(1.5), t.Space(1.5)).Shrink(0).Radius(999).Background(ink)
						ui.Text(c, text).FontSize(t.Rem(0.6)).TextColor(ink).SingleLine().MinWidth(0)
					})
				})
			})
			if b.Clicked() && !on && day.Available > 0 {
				*picked = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("SlotCalendar", "picked", "index", i)
			}
		}
	})
	return root
}
