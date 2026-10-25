package mail

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SendLaterOptions word the card and choices; Format shows times.
type SendLaterOptions struct{ Title, Tomorrow, Monday, Pick, Format string }

// SendLater selects a scheduled send time; cue Select.
func SendLater(c *ui.Context, at *time.Time, o SendLaterOptions) ui.Element {
	if at == nil || at.IsZero() {
		panic("mail: SendLater needs a time")
	}
	def(&o.Title, "Schedule send")
	def(&o.Tomorrow, "Tomorrow morning")
	def(&o.Monday, "Monday morning")
	def(&o.Pick, "Pick date")
	def(&o.Format, "Mon, Jan 2 15:04")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "calendar-clock", o.Title, nil)
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			now := c.Now()
			choices := []struct {
				label, mark string
				at          time.Time
			}{
				{o.Tomorrow, "sunrise", time.Date(now.Year(), now.Month(), now.Day()+1, 9, 0, 0, 0, now.Location())},
				{o.Monday, "calendar-days", mondayAfter(now)},
			}
			for i, choice := range choices {
				c.Key(i)
				b := timeChoice(c, choice.label, choice.mark, choice.at.Format(o.Format), at.Equal(choice.at), false)
				kit.Stagger(c, list, b, i)
				if b.Clicked() && !at.Equal(choice.at) {
					*at = choice.at
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("SendLater", "selected", "choice", i)
				}
			}
			ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2.5)).Children(func() {
				glyph(c, "calendar", 16).TextColor(t.TextMuted)
				if kit.Changed(ui.DateInput(c, at).Label(o.Pick).Grow(1), sound.Select) {
					root.SetChanged()
					kit.Log("SendLater", "picked")
				}
			})
		})
	})
	return root
}

// mondayAfter is 9:00 on the first Monday after tomorrow.
func mondayAfter(now time.Time) time.Time {
	days := (8 - int(now.Weekday())) % 7
	if days <= 1 {
		days += 7
	}
	return time.Date(now.Year(), now.Month(), now.Day()+days, 9, 0, 0, 0, now.Location())
}
