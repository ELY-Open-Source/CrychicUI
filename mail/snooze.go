package mail

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SnoozeOptions word the card and its choices; Format shows times.
type SnoozeOptions struct{ Title, Later, Tomorrow, Pick, Format string }

// Snooze selects a later reading time; cue Select.
func Snooze(c *ui.Context, at *time.Time, o SnoozeOptions) ui.Element {
	if at == nil || at.IsZero() {
		panic("mail: Snooze needs a time")
	}
	def(&o.Title, "Snooze until")
	def(&o.Later, "Later today")
	def(&o.Tomorrow, "Tomorrow")
	def(&o.Pick, "Pick date")
	def(&o.Format, "Mon 15:04")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "clock", o.Title, nil)
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			now := c.Now()
			later, available := laterToday(now)
			tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 9, 0, 0, 0, now.Location())
			choices := []struct {
				label, mark string
				at          time.Time
				available   bool
			}{{o.Later, "sunset", later, available}, {o.Tomorrow, "sunrise", tomorrow, true}}
			for i, choice := range choices {
				c.Key(i)
				b := timeChoice(c, choice.label, choice.mark, choice.at.Format(o.Format), at.Equal(choice.at), !choice.available)
				kit.Stagger(c, list, b, i)
				if b.Clicked() && choice.available && !at.Equal(choice.at) {
					*at = choice.at
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Snooze", "selected", "choice", i)
				}
			}
			ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2.5)).Children(func() {
				glyph(c, "calendar-days", 16).TextColor(t.TextMuted)
				if kit.Changed(ui.DateInput(c, at).Label(o.Pick).Grow(1), sound.Select) {
					root.SetChanged()
					kit.Log("Snooze", "picked")
				}
			})
		})
	})
	return root
}

// laterToday is three hours on, on the local hour.
func laterToday(now time.Time) (time.Time, bool) {
	h := now.Add(3 * time.Hour)
	later := h.Add(-time.Duration(h.Minute())*time.Minute - time.Duration(h.Second())*time.Second - time.Duration(h.Nanosecond()))
	return later, later.Year() == now.Year() && later.YearDay() == now.YearDay()
}
