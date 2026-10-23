package schedule

import (
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BookingOptions are the booking labels.
type BookingOptions struct{ Confirm, DayFormat string }

// Booking chooses day and slot then confirms; cue Success.
func Booking(c *ui.Context, days []time.Time, slots []string, day *int, slot *string, confirmed *bool, o BookingOptions) (root, confirm ui.Element) {
	if len(days) == 0 || len(slots) == 0 || day == nil || slot == nil || *day < -1 || *day >= len(days) || confirmed == nil {
		panic("schedule: Booking has invalid data")
	}
	def(&o.Confirm, "Book")
	def(&o.DayFormat, "Mon 02")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).Wrap().Gap(t.Space(1)).Children(func() {
			for i, date := range days {
				b := kit.Press(ui.ButtonBase(c).Label(date.Format(o.DayFormat)).Padding(t.Space(1), t.Space(2)).Radius(999).
					Border(1, t.Border).Checked(*day == i))
				b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*day == i)), ui.Snappy))
				b.Children(func() { ui.Text(c, date.Format(o.DayFormat)).FontSize(t.Rem(0.75)) })
				if b.Clicked() && *day != i {
					*day = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Booking", "day", "index", i)
				}
			}
		})
		ui.Row(c).Wrap().Gap(t.Space(1)).Children(func() {
			for _, value := range slots {
				if value == "" {
					panic("schedule: Booking slot must not be empty")
				}
				b := kit.Press(ui.ButtonBase(c).Label(value).Padding(t.Space(1), t.Space(2)).Radius(999).
					Border(1, t.Border).Checked(*slot == value))
				b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*slot == value)), ui.Snappy))
				b.Children(func() { ui.Text(c, value).FontSize(t.Rem(0.75)) })
				if b.Clicked() && *slot != value {
					*slot = value
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Booking", "slot", "slot", *slot)
				}
			}
		})
		confirm = kit.Press(ui.ButtonBase(c).Label(o.Confirm).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
			Background(t.Accent).TextColor(t.AccentText).Disabled(*day < 0 || !slices.Contains(slots, *slot)))
		confirm.Children(func() { ui.Text(c, o.Confirm).SingleLine() })
	})
	if confirm.Clicked() && *day >= 0 && slices.Contains(slots, *slot) {
		*confirmed = true
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("Booking", "confirmed")
	}
	return root, confirm
}
