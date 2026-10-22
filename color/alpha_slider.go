package color

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AlphaSlider sets a color's opacity over a checkered track.
func AlphaSlider(c *ui.Context, label string, value *ui.Color) ui.Element {
	if label == "" || value == nil {
		panic("color: AlphaSlider needs a label and a color")
	}
	t := c.Theme()
	pad, track := t.Space(3), t.Space(3)
	a := float64(value.A) / 255
	s := ui.SliderBase(c, &a, 0, 1).Step(0.01).Label(label).Width(t.Space(64)).Height(t.Space(6)).PaddingX(pad)
	step := ui.Local(s, "step", func() int { return int(value.A) / 13 })
	_, _, held := s.Dragged()
	if s.Changed() {
		value.A = uint8(a*255 + 0.5)
		if n := int(value.A) / 13; !held || n != *step {
			*step = n
			sound.Play(sound.Tick)
			kit.Log("AlphaSlider", "moved", "step", n)
		}
	}
	rtl := c.Direction() == ui.RightToLeft
	shown := s.AnimateTo("a", float32(a), ui.Snappy)
	opaque := *value
	opaque.A = 255
	grip := s.AnimateTo("grip", unit(held), ui.Snappy)
	return s.Draw(func(p *ui.Painter, r ui.Rect) {
		bar := ui.Rect{X: r.X + pad, Y: r.Y + (r.H-track)/2, W: r.W - 2*pad, H: track}
		from, to := opaque.Alpha(0), opaque
		if rtl {
			from, to = to, from
		}
		p.Clip(bar, track/2, func() {
			checker(p, bar, t)
			p.FillGradient(bar, ui.LinearGradient{From: from, To: to, Angle: 90}, 0)
		})
		x := bar.X + bar.W*shown
		if rtl {
			x = bar.X + bar.W*(1-shown)
		}
		var k ui.Path
		p.FillPath(k.Circle(x, r.Y+r.H/2, (r.H/2-1)*(0.85+0.15*grip)), t.Background)
		p.StrokePath(&k, 1, t.Border)
	})
}
