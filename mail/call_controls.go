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
	root := card(c).Row().Gap(t.Space(5)).Justify(ui.Center).Padding(t.Space(3), t.Space(4))
	root.Children(func() {
		quiet := t.SurfaceHover
		controls := []struct {
			label, on, off string
			value          *bool
			ink            ui.Color
			alarm          bool
		}{
			{o.Mute, "mic-off", "mic", &state.Muted, t.Danger, state.Muted},
			{o.Video, "video", "video-off", &state.Video, t.Danger, !state.Video},
			{o.Share, "monitor-up", "monitor-up", &state.Sharing, t.Accent, state.Sharing},
		}
		for i, control := range controls {
			c.Key(i)
			mark, bg, ink := control.off, quiet, t.Text
			if *control.value {
				mark = control.on
			}
			if control.alarm {
				bg, ink = control.ink.Alpha(0.14), control.ink
			}
			b := roundAction(c, control.label, mark, bg, ink).Checked(*control.value).Disabled(state.Ended)
			if b.Clicked() && !state.Ended {
				*control.value = !*control.value
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("CallControls", "toggled", "control", i, "on", *control.value)
			}
		}
		end := roundAction(c, o.End, "phone-off", t.Danger, t.AccentText).Disabled(state.Ended)
		if end.Clicked() && !state.Ended {
			state.Ended = true
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("CallControls", "ended")
		}
	})
	return root
}
