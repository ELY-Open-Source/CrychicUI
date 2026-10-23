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
type RSVPOptions struct{ Yes, Maybe, No string }

// RSVP selects yes, maybe, or no; cue Select.
func RSVP(c *ui.Context, response *RSVPResponse, o RSVPOptions) ui.Element {
	if response == nil || *response > No {
		panic("schedule: RSVP needs a valid response")
	}
	def(&o.Yes, "Yes")
	def(&o.Maybe, "Maybe")
	def(&o.No, "No")
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, item := range []struct {
			value RSVPResponse
			name  string
		}{{Yes, o.Yes}, {Maybe, o.Maybe}, {No, o.No}} {
			b := kit.Press(ui.ButtonBase(c).Label(item.name).Padding(t.Space(1.5), t.Space(3)).Radius(999).
				Border(1, t.Border).Checked(*response == item.value))
			b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*response == item.value)), ui.Snappy))
			b.Children(func() { ui.Text(c, item.name).SingleLine() })
			if b.Clicked() && *response != item.value {
				*response = item.value
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("RSVP", "selected", "response", *response)
			}
		}
	})
	return root
}
