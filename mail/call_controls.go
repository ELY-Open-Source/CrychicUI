package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CallState is the state of call controls.
type CallState struct {
	Muted, Video, Sharing, Ended bool
}

// CallControlsOptions are the call control labels.
type CallControlsOptions struct{ Mute, Video, Share, End string }

// CallControls toggles call controls; cue Toggle.
func CallControls(c *ui.Context, state *CallState, o CallControlsOptions) ui.Element {
	if state == nil {
		panic("mail: CallControls needs state")
	}
	def(&o.Mute, "Mute")
	def(&o.Video, "Video")
	def(&o.Share, "Share")
	def(&o.End, "End")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		controls := []struct {
			label string
			value *bool
		}{{o.Mute, &state.Muted}, {o.Video, &state.Video}, {o.Share, &state.Sharing}}
		for i, control := range controls {
			c.Key(i)
			b := button(c, control.label).Checked(*control.value).Disabled(state.Ended)
			b.Children(func() { ui.Text(c, control.label) })
			if b.Clicked() && !state.Ended {
				*control.value = !*control.value
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("CallControls", "toggled", "control", i, "on", *control.value)
			}
		}
		end := button(c, o.End).Background(t.Danger).TextColor(t.AccentText).Disabled(state.Ended)
		end.Children(func() { ui.Text(c, o.End) })
		if end.Clicked() && !state.Ended {
			state.Ended = true
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("CallControls", "ended")
		}
	})
	return root
}
