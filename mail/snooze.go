package mail

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SnoozeOptions are the snooze choices.
type SnoozeOptions struct{ Later, Tomorrow, Pick string }

// Snooze selects a later reading time; cue Select.
func Snooze(c *ui.Context, at *time.Time, o SnoozeOptions) ui.Element {
	if at == nil || at.IsZero() {
		panic("mail: Snooze needs a time")
	}
	def(&o.Later, "Later today")
	def(&o.Tomorrow, "Tomorrow")
	def(&o.Pick, "Pick date")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		now := c.Now()
		later, available := laterToday(now)
		tomorrow := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
		choices := []struct {
			label     string
			at        time.Time
			available bool
		}{{o.Later, later, available}, {o.Tomorrow, tomorrow, true}}
		for i, choice := range choices {
			c.Key(i)
			b := button(c, choice.label).FillWidth().Disabled(!choice.available)
			b.Children(func() { ui.Text(c, choice.label) })
			if b.Clicked() && choice.available {
				*at = choice.at
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Snooze", "selected", "choice", i)
			}
		}
		if kit.Changed(ui.DateInput(c, at).Label(o.Pick), sound.Select) {
			root.SetChanged()
			kit.Log("Snooze", "picked")
		}
	})
	return root
}

func laterToday(now time.Time) (time.Time, bool) {
	later := now.Add(3 * time.Hour)
	return later, later.Year() == now.Year() && later.YearDay() == now.YearDay()
}
