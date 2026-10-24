package schedule

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RSVPResponse is one RSVP choice.
type RSVPResponse uint8

const (
	RSVPNone RSVPResponse = iota
	Yes
	Maybe
	No
)

// RSVPOptions are the RSVP labels.
type RSVPOptions struct{ Prompt, Yes, Maybe, No string }

// RSVP selects yes, maybe, or no; cue Select.
func RSVP(c *ui.Context, response *RSVPResponse, o RSVPOptions) ui.Element {
	if response == nil || *response > No {
		panic("schedule: RSVP needs a valid response")
	}
	def(&o.Prompt, "Going?")
	def(&o.Yes, "Yes")
	def(&o.Maybe, "Maybe")
	def(&o.No, "No")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Text(c, o.Prompt).FontSize(t.Rem(0.85)).Bold().Grow(1).Shrink(0)
		for i, item := range []struct {
			value       RSVPResponse
			name, glyph string
			ink         ui.Color
		}{{Yes, o.Yes, "check", t.Success}, {Maybe, o.Maybe, "circle-help", t.Warning}, {No, o.No, "x", t.Danger}} {
			on := *response == item.value
			c.Key(i)
			b := pill(c, item.name, on, item.ink, func() {
				c.Key(on)
				kit.Arrive(glyph(c, item.glyph, 14), ui.Pose{Scale: 0.5, Fade: 1}, 0)
				ui.Text(c, item.name).Bold().Shrink(0)
			})
			if b.Clicked() && !on {
				*response = item.value
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("RSVP", "selected", "response", *response)
			}
		}
	})
	return root
}
