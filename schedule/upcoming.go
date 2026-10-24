package schedule

import (
	"fmt"
	"sort"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// UpcomingOptions word the list.
type UpcomingOptions struct{ Title, DayFormat, Empty string }

// Upcoming shows the next events, grouped by day.
func Upcoming(c *ui.Context, events []Event, o UpcomingOptions) ui.Element {
	validEvents("Upcoming", events)
	def(&o.Title, "Upcoming")
	def(&o.DayFormat, "Monday, Jan 2")
	def(&o.Empty, "Nothing scheduled")
	t := c.Theme()
	list := append([]Event(nil), events...)
	sort.SliceStable(list, func(i, j int) bool { return list[i].Start.Before(list[j].Start) })
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "calendar-clock", o.Title, func() { kit.Chip(c, fmt.Sprint(len(list)), t.Accent) })
		body := ui.Column(c).FillWidth().Padding(t.Space(1), 0, t.Space(2), 0)
		body.Children(func() {
			if len(list) == 0 {
				c.Key("empty")
				empty := ui.Column(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(2)).Padding(t.Space(6)).TextColor(t.TextMuted)
				kit.Arrive(empty, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
					ui.Icon(c, icon.Must("calendar-x")).Size(24, 24)
					ui.Text(c, o.Empty).FontSize(t.Rem(0.85))
				})
				return
			}
			// One zone, the earliest event's, for days and times.
			loc := list[0].Start.Location()
			for i := 0; i < len(list); {
				first := list[i].Start.In(loc)
				end := i + 1
				for end < len(list) && sameDate(list[end].Start, first) {
					end++
				}
				lo, hi := i, end
				c.Key(first.Format("2006-01-02"))
				ui.Column(c).FillWidth().Children(func() {
					kit.Muted(c, first.Format(o.DayFormat)).FontSize(t.Rem(0.7)).Bold().Padding(t.Space(2), t.Space(3), t.Space(1), t.Space(3))
					// Rows share a parent only with rows.
					ui.Column(c).FillWidth().Children(func() {
						for k, e := range list[lo:hi] {
							c.Key(e.ID)
							upcomingRow(c, body, e, loc, lo+k)
						}
					})
				})
				i = end
			}
		})
	})
	return root
}

// upcomingRow is one timed event line.
func upcomingRow(c *ui.Context, list ui.Element, e Event, loc *time.Location, i int) {
	t := c.Theme()
	ink := inks(t)[e.Tint]
	row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Stretch).Padding(t.Space(1.5), t.Space(3))
	row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(0.6*unit(row.Hovered())), ui.Snappy))
	kit.Stagger(c, list, row, i).Children(func() {
		ui.Column(c).Width(t.Space(12)).Shrink(0).Children(func() {
			kit.Numeric(c, e.Start.In(loc).Format("15:04")).FontSize(t.Rem(0.8)).Bold()
			kit.Muted(c, e.End.In(loc).Format("15:04")).FontSize(t.Rem(0.7))
		})
		ui.Box(c).Width(t.Space(0.75)).Shrink(0).Radius(999).Background(ink)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, e.Title).FontSize(t.Rem(0.85)).Bold().SingleLine().MinWidth(0)
			if e.Place != "" {
				meta(c, "map-pin", e.Place)
			}
		})
	})
}
