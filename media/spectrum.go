package media

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SpectrumOptions tune a Spectrum.
type SpectrumOptions struct {
	// Bars counts the bands; zero gives 24.
	Bars int
	// Height of the bars; zero gives the default.
	Height float32
}

// Spectrum moves its bars like a spectrum while playing.
func Spectrum(c *ui.Context, label string, playing bool, o SpectrumOptions) ui.Element {
	if label == "" || o.Bars < 0 {
		panic(fmt.Sprintf("media: Spectrum %q needs a label and bars not negative", label))
	}
	if o.Bars == 0 {
		o.Bars = 24
	}
	t := c.Theme()
	size("Spectrum", label, &o.Height, t.Space(20))
	box := ui.Box(c).FillWidth().Height(o.Height).Role(ui.RoleImage).Label(label)
	hs := bands(c, box, o.Bars, playing)
	return box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		w := r.W / float32(len(hs))
		for i, h := range hs {
			bh := max(t.Space(0.5), h*r.H)
			p.Fill(ui.Rect{X: r.X + float32(i)*w + w*0.15, Y: r.Y + r.H - bh, W: w * 0.7, H: bh}, t.Accent.Mix(t.Success, float32(i)/float32(len(hs))), w*0.2)
		}
	}))
}
