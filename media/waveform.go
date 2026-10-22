package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// WaveformOptions tune a Waveform.
type WaveformOptions struct {
	// Height of the bars; zero gives the default.
	Height float32
}

// Waveform draws a clip's levels; the played part is tinted.
func Waveform(c *ui.Context, label string, levels []float64, played float64, o WaveformOptions) ui.Element {
	if label == "" || len(levels) == 0 || !(played >= 0 && played <= 1) {
		panic(fmt.Sprintf("media: Waveform %q needs a label, levels and played from 0 to 1", label))
	}
	for _, v := range levels {
		if !(v >= 0 && v <= 1) {
			panic(fmt.Sprintf("media: Waveform %q levels must be from 0 to 1", label))
		}
	}
	t := c.Theme()
	size("Waveform", label, &o.Height, t.Space(16))
	box := ui.Box(c).FillWidth().Height(o.Height).Role(ui.RoleMeter).Label(label).Range(0, 1, played)
	g := box.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	p := box.AnimateTo("played", float32(played), ui.Snappy)
	n := len(levels)
	return box.Draw(kit.Inside(0, func(pt *ui.Painter, r ui.Rect) {
		w := r.W / float32(n)
		for i, v := range levels {
			ink := t.TextMuted.Alpha(0.35)
			if (float32(i)+0.5)/float32(n) <= p {
				ink = t.Accent
			}
			bh := max(t.Space(0.5), float32(v)*r.H*g)
			pt.Fill(ui.Rect{X: r.X + float32(i)*w + w*0.2, Y: r.Y + (r.H-bh)/2, W: w * 0.6, H: bh}, ink, w*0.3)
		}
	}))
}
