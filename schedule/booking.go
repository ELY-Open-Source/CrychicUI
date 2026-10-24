package schedule

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BookingOptions are the booking labels; Done takes day and slot.
type BookingOptions struct{ Title, Confirm, DayFormat, Done string }

// Booking chooses day and slot then confirms; cue Success.
func Booking(c *ui.Context, days []time.Time, slots []string, day *int, slot *string, confirmed *bool, o BookingOptions) (root, confirm ui.Element) {
	if len(days) == 0 || len(slots) == 0 || day == nil || slot == nil || *day < -1 || *day >= len(days) || confirmed == nil {
		panic("schedule: Booking has invalid data")
	}
	for i, value := range slots {
		if value == "" || slices.Contains(slots[:i], value) {
			panic(fmt.Sprintf("schedule: Booking slot %q is empty or repeated", value))
		}
	}
	for i, date := range days {
		for _, prior := range days[:i] {
			if prior.Equal(date) {
				panic(fmt.Sprintf("schedule: Booking day %v repeats", date))
			}
		}
	}
	def(&o.Title, "Book a time")
	def(&o.Confirm, "Book")
	def(&o.DayFormat, "Mon 02")
	def(&o.Done, "Booked for %s at %s")
	t := c.Theme()
	ready := *day >= 0 && slices.Contains(slots, *slot)
	root = card(c)
	root.Children(func() {
		kit.CardHead(c, "calendar-check", o.Title, nil)
		list := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3))
		list.Children(func() {
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).Children(func() {
				for i, date := range days {
					text := date.Format(o.DayFormat)
					c.Key(instant(date))
					b := pill(c, text, *day == i, t.Accent, func() { ui.Text(c, text).FontSize(t.Rem(0.75)).Bold().Shrink(0) })
					if b.Clicked() && *day != i {
						*day, *confirmed = i, false
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("Booking", "day", "index", i)
					}
				}
			})
			grid := ui.Grid(c).Columns(3).FillWidth().Gap(t.Space(1.5))
			grid.Children(func() {
				for i, value := range slots {
					c.Key(value)
					b := pill(c, value, *slot == value, t.Accent, func() {
						glyph(c, "clock", 12)
						ui.Text(c, value).FontSize(t.Rem(0.75)).Shrink(0)
					})
					kit.Stagger(c, grid, b, i)
					if b.Clicked() && *slot != value {
						*slot, *confirmed = value, false
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("Booking", "slot", "index", i)
					}
				}
			})
			c.Key("confirm")
			confirm = kit.Press(ui.ButtonBase(c).Label(o.Confirm).FillWidth().Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
				TextColor(t.AccentText).Disabled(!ready || *confirmed))
			confirm.Background(confirm.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(confirm.Hovered())), ui.Snappy))
			confirm.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					glyph(c, "calendar-check", 16)
					ui.Text(c, o.Confirm).Bold().Shrink(0)
				})
			})
			if *confirmed && ready {
				c.Key("done")
				banner := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
					Background(t.Success.Alpha(0.12)).TextColor(t.Success)
				kit.Arrive(banner, ui.Pose{Y: t.Space(1), Scale: 0.97, Fade: 1}, 0).Children(func() {
					glyph(c, "circle-check", 16)
					ui.Text(c, fmt.Sprintf(o.Done, days[*day].Format(o.DayFormat), *slot)).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).Wrap()
				})
			}
		})
	})
	if confirm.Clicked() && ready && !*confirmed {
		*confirmed = true
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("Booking", "confirmed")
	}
	return root, confirm
}
