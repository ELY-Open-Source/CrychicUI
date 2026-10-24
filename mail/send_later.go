package mail

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SendLaterOptions are the scheduled send choices.
type SendLaterOptions struct{ Tomorrow, Monday, Pick string }

// SendLater selects a scheduled send time; cue Select.
func SendLater(c *ui.Context, at *time.Time, o SendLaterOptions) ui.Element {
	if at == nil || at.IsZero() {
		panic("mail: SendLater needs a time")
	}
	def(&o.Tomorrow, "Tomorrow morning")
	def(&o.Monday, "Monday morning")
	def(&o.Pick, "Pick date")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		now := c.Now()
		morning := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
		days := (8 - int(now.Weekday())) % 7
		if days == 0 {
			days = 7
		}
		choices := []struct {
			label string
			at    time.Time
		}{{o.Tomorrow, morning.AddDate(0, 0, 1)}, {o.Monday, morning.AddDate(0, 0, days)}}
		for i, choice := range choices {
			c.Key(i)
			b := button(c, choice.label).FillWidth()
			b.Children(func() { ui.Text(c, choice.label) })
			if b.Clicked() {
				*at = choice.at
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("SendLater", "selected", "choice", i)
			}
		}
		if kit.Changed(ui.DateInput(c, at).Label(o.Pick), sound.Select) {
			root.SetChanged()
			kit.Log("SendLater", "picked")
		}
	})
	return root
}
