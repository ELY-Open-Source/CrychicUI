package schedule

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MonthViewOptions word the month; Now overrides the clock.
type MonthViewOptions struct {
	Format, More string
	Weekdays     [7]string
	Now          time.Time
}

// MonthView shows month cells with event chips; cue Select.
func MonthView(c *ui.Context, month time.Time, events []Event, selected *string, o MonthViewOptions) ui.Element {
	if selected == nil {
		panic("schedule: MonthView needs selection")
	}
	validEvents("MonthView", events)
	def(&o.Format, "January 2006")
	def(&o.More, "+%d more")
	if o.Weekdays == [7]string{} {
		o.Weekdays = [7]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	}
	now := o.Now
	if now.IsZero() {
		now = c.Now()
	}
	t := c.Theme()
	first := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	start := first.AddDate(0, 0, -int(first.Weekday()))
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "calendar", month.Format(o.Format), nil)
		ui.Grid(c).Columns(7).FillWidth().Padding(t.Space(1), t.Space(2), 0, t.Space(2)).Children(func() {
			for _, name := range o.Weekdays {
				kit.Muted(c, name).FontSize(t.Rem(0.65)).Bold().Padding(0, t.Space(1))
			}
		})
		grid := ui.Grid(c).Columns(7).FillWidth().Padding(t.Space(1), t.Space(2), t.Space(2), t.Space(2))
		grid.Children(func() {
			for i := range 42 {
				day := start.AddDate(0, 0, i)
				inside := day.Month() == month.Month()
				var list []Event
				for _, lane := range lanes(events, day, 0) {
					list = append(list, lane...)
				}
				c.Key(i)
				cell := ui.Column(c).MinWidth(0).Height(t.Space(19)).Gap(t.Space(0.5)).Padding(t.Space(1)).Clip().
					BorderWidth(1, 0, 0, 0).BorderColor(t.Border)
				cell.Background(cell.AnimateColor("bg", t.SurfaceHover.Alpha(0.6*unit(cell.Hovered())), ui.Snappy))
				kit.Stagger(c, grid, cell, i/7).Children(func() {
					today := sameDate(now, day)
					num := ui.Box(c).Size(t.Space(5), t.Space(5)).Center().Radius(999).Background(t.Accent.Alpha(unit(today)))
					num.Children(func() {
						ink := t.Text
						if !inside {
							ink = t.TextMuted.Alpha(0.6)
						}
						if today {
							ink = t.AccentText
						}
						ui.Text(c, fmt.Sprint(day.Day())).FontSize(t.Rem(0.7)).Bold().TextColor(ink)
					})
					for k, e := range list {
						if k == 2 && len(list) > 3 {
							kit.Muted(c, fmt.Sprintf(o.More, len(list)-2)).FontSize(t.Rem(0.6)).Padding(0, t.Space(1))
							break
						}
						monthChip(c, root, e, selected)
					}
				})
			}
		})
	})
	return root
}

// monthChip is one selectable event line in a cell.
func monthChip(c *ui.Context, root ui.Element, e Event, selected *string) {
	t := c.Theme()
	ink := inks(t)[e.Tint]
	on := *selected == e.ID
	c.Key(e.ID)
	b := kit.Press(ui.ButtonBase(c).Label(e.Title).FillWidth().MinWidth(0).Padding(t.Space(0.25), t.Space(1)).Radius(t.Space(1)).Checked(on))
	b.Background(b.AnimateColor("bg", ink.Alpha(0.1).Mix(ink.Alpha(0.24), unit(on || b.Hovered())), ui.Snappy))
	b.Children(func() {
		ui.Row(c).FillWidth().MinWidth(0).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(1.5), t.Space(1.5)).Shrink(0).Radius(999).Background(ink)
			ui.Text(c, e.Title).FontSize(t.Rem(0.6)).Bold().TextColor(ink.Mix(t.Text, 0.35)).SingleLine().Grow(1).MinWidth(0)
		})
	})
	if b.Clicked() && *selected != e.ID {
		*selected = e.ID
		root.SetChanged()
		sound.Play(sound.Select)
		kit.Log("MonthView", "selected", "event", e.ID)
	}
}
