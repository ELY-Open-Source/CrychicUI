package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Clip holds one note per beat; zero rests.
type Clip struct {
	Notes []float64
	Beat  time.Duration
}

// AudioPlayerOptions tune an AudioPlayer.
type AudioPlayerOptions struct {
	Play, Seek string
}

// AudioPlayer plays a synthesized clip with play, seek and time.
func AudioPlayer(c *ui.Context, label string, clip Clip, pos *time.Duration, playing *bool, o AudioPlayerOptions) ui.Element {
	if label == "" || len(clip.Notes) == 0 || clip.Beat <= 0 || playing == nil {
		panic(fmt.Sprintf("media: AudioPlayer %q needs a label, notes, a beat and state", label))
	}
	for _, hz := range clip.Notes {
		if hz != 0 && !(hz >= 20 && hz <= 20000) {
			panic(fmt.Sprintf("media: AudioPlayer %q notes must be 0 or 20–20000 Hz", label))
		}
	}
	total := clip.Beat * time.Duration(len(clip.Notes))
	if pos == nil || *pos < 0 || *pos > total {
		panic(fmt.Sprintf("media: AudioPlayer %q needs a position within the clip", label))
	}
	def(&o.Play, "Play")
	def(&o.Seek, "Position")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3)).Radius(t.Radius).
		Border(1, t.Border).Label(label)
	now := c.Now()
	last := ui.Local(root, "last", func() time.Time { return time.Time{} })
	beat := ui.Local(root, "beat", func() int { return -1 })
	sing := func() {
		if i := min(int(*pos/clip.Beat), len(clip.Notes)-1); i != *beat {
			*beat = i
			if hz := clip.Notes[i]; hz > 0 {
				sound.PlayNote(hz)
			}
		}
	}
	// track follows *playing whichever control changed it.
	track := func() {
		switch {
		case !*playing:
			*last = time.Time{}
		case last.IsZero():
			if *pos >= total {
				*pos = 0
			}
			*last, *beat = now, -1
		default:
			*pos = min(total, *pos+now.Sub(*last))
			*last = now
			if *pos >= total {
				*playing, *last = false, time.Time{}
				root.SetChanged()
				kit.Log("AudioPlayer", "ended")
			}
		}
	}
	track()
	return kit.Arrive(root, ui.Pose{Fade: 1}, 0).Children(func() {
		if button.PlayButton(c, playing, button.PlayButtonOptions{Label: o.Play}).Changed() {
			root.SetChanged()
			track()
		}
		if *playing {
			sing()
			c.AnimationFrame()
		}
		if SeekBar(c, o.Seek, pos, total, total, SeekBarOptions{Step: clip.Beat}).Grow(1).Changed() {
			root.SetChanged()
			*beat = -1
			if *playing {
				sing()
			}
		}
	})
}
