package chat

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/media"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VoiceMessageOptions name the actions.
type VoiceMessageOptions struct {
	Play, Pause, Waveform string
}

// VoiceMessage is play, a waveform and the time left.
func VoiceMessage(c *ui.Context, label string, levels []float64, length time.Duration, playing *bool, o VoiceMessageOptions) ui.Element {
	if label == "" || len(levels) == 0 || length <= 0 || playing == nil {
		panic(fmt.Sprintf("chat: VoiceMessage %q needs a label, levels, a length and playing", label))
	}
	def(&o.Play, "Play")
	def(&o.Pause, "Pause")
	def(&o.Waveform, "Waveform")
	t := c.Theme()
	root := shell(c, false).Row().Gap(t.Space(2)).AlignItems(ui.Center).Width(t.Space(70)).Label(label)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	now := c.Now()
	pos := ui.Local(root, "pos", func() time.Duration { return 0 })
	last := ui.Local(root, "last", func() time.Time { return time.Time{} })
	switch {
	case !*playing:
		*last = time.Time{}
	case last.IsZero():
		if *pos >= length {
			*pos = 0
		}
		*last = now
	default:
		*pos = min(length, *pos+now.Sub(*last))
		*last = now
		if *pos >= length {
			*playing, *last = false, time.Time{}
			root.SetChanged()
			kit.Log("VoiceMessage", "ended")
		}
	}
	if *playing {
		c.After(time.Second - *pos%time.Second)
	}
	name, glyph := o.Play, "play"
	if *playing {
		name, glyph = o.Pause, "pause"
	}
	left := (length - *pos + time.Second - 1) / time.Second * time.Second
	var toggle ui.Element
	root.Children(func() {
		toggle = round(c, name, glyph)
		media.Waveform(c, o.Waveform, levels, float64(*pos)/float64(length), media.WaveformOptions{Height: t.Space(8)}).Grow(1)
		kit.Numeric(c, clock(left)).FontSize(t.Rem(0.85)).TextColor(t.TextMuted)
	})
	if toggle.Clicked() {
		*playing = !*playing
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("VoiceMessage", "playing", *playing)
	}
	return root
}
