package mail

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AutoReply is one automatic reply configuration.
type AutoReply struct {
	Enabled    bool
	Start, End time.Time
	Message    string
}

// OutOfOfficeOptions are the automatic reply labels.
type OutOfOfficeOptions struct{ Enabled, Start, End, Message string }

// OutOfOffice edits automatic reply dates and message; cue Toggle.
func OutOfOffice(c *ui.Context, value *AutoReply, o OutOfOfficeOptions) ui.Element {
	if value == nil || value.Start.IsZero() || value.End.IsZero() || value.End.Before(value.Start) {
		panic("mail: OutOfOffice needs a valid range")
	}
	def(&o.Enabled, "Out of office")
	def(&o.Start, "Start")
	def(&o.End, "End")
	def(&o.Message, "Message")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		toggle := button(c, o.Enabled).Checked(value.Enabled)
		toggle.Children(func() { ui.Text(c, o.Enabled) })
		if toggle.Clicked() {
			value.Enabled = !value.Enabled
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("OutOfOffice", "toggled", "enabled", value.Enabled)
		}
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).Children(func() {
			start := ui.DateInput(c, &value.Start).Label(o.Start)
			end := ui.DateInput(c, &value.End).Label(o.End)
			if start.Changed() || end.Changed() {
				if value.End.Before(value.Start) {
					if start.Changed() {
						value.End = value.Start
					} else {
						value.Start = value.End
					}
				}
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("OutOfOffice", "dates")
			}
		})
		in := ui.TextArea(c, &value.Message).Label(o.Message).Placeholder(o.Message).FillWidth().Lines(3, 6)
		if in.Changed() {
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("OutOfOffice", "edited", "runes", len([]rune(value.Message)))
		}
	})
	return root
}
