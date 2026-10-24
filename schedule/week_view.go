package schedule

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// WeekViewOptions are the week labels; Now overrides the clock.
type WeekViewOptions struct {
	DayFormat string
	Now       time.Time
}

// WeekView shows days with hours and event blocks; cue Select.
func WeekView(c *ui.Context, days []time.Time, events []Event, selected *string, o WeekViewOptions) ui.Element {
	if len(days) != 7 || selected == nil {
		panic("schedule: WeekView needs seven days and selection")
	}
	def(&o.DayFormat, "Mon 02")
	validEvents("WeekView", events)
	now := o.Now
	if now.IsZero() {
		now = c.Now()
	}
	t := c.Theme()
	hourH := t.Space(10)
	root := card(c)
	scroll := ui.Local(root, "scroll", func() *ui.ScrollState { return &ui.ScrollState{Y: firstHour(days, events) * hourH} })
	root.Children(func() {
		ui.Row(c).FillWidth().Padding(t.Space(2), 0).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			ui.Box(c).Width(t.Space(11)).Shrink(0)
			for _, day := range days {
				today := sameDate(now, day)
				label := day.Format(o.DayFormat)
				ui.Column(c).Grow(1).Basis(0).MinWidth(0).Padding(0, t.Space(0.5)).Children(func() {
					ui.Text(c, label).FontSize(t.Rem(0.75)).Bold().SingleLine().MinWidth(0).
						TextColor(t.TextMuted.Mix(t.Accent, unit(today)))
				})
			}
		})
		ui.Scroll(c).TrackScroll(*scroll).FillWidth().Height(t.Space(90)).Children(func() {
			ui.Row(c).FillWidth().Children(func() {
				gutter(c, hourH)
				for i, day := range days {
					c.Key(i)
					dayColumn(c, day, events, now, hourH, func(e Event, top, height float32) {
						on := *selected == e.ID
						b := ui.ButtonBase(c).Label(e.Title).Absolute().Left(t.Space(0.5)).Right(t.Space(0.5)).Top(top).Height(height).Checked(on)
						kit.Press(eventBox(c, e, day.Location(), b, on))
						if b.Clicked() && *selected != e.ID {
							*selected = e.ID
							root.SetChanged()
							sound.Play(sound.Select)
							kit.Log("WeekView", "selected", "event", e.ID)
						}
					})
				}
			})
		})
	})
	return root
}

// gutter is the hour label column.
func gutter(c *ui.Context, hourH float32) {
	t := c.Theme()
	ui.Column(c).Width(t.Space(11)).Shrink(0).Children(func() {
		for h := range 24 {
			ui.Box(c).Height(hourH).Padding(0, t.Space(1.5), 0, 0).Children(func() {
				kit.Muted(c, time.Date(2000, 1, 1, h, 0, 0, 0, time.UTC).Format("15:04")).FontSize(t.Rem(0.65)).AlignSelf(ui.End)
			})
		}
	})
}

// firstHour opens the grid at the earliest event, else 8.
func firstHour(days []time.Time, events []Event) float32 {
	first, found := float32(24), false
	for _, day := range days {
		for _, e := range events {
			if from, _, ok := span(e, day); ok {
				first, found = min(first, max(0, float32(int(hours(from, day)))-1)), true
			}
		}
	}
	if !found {
		return 8
	}
	return first
}
