package button

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MagneticButton leans toward the pointer and springs back.
func MagneticButton(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Background(t.Surface)
	bounds := root.LayoutBounds()
	kit.Settle(c, root, []float32{bounds.W, bounds.H}, func() []float32 {
		b := root.LayoutBounds()
		return []float32{b.W, b.H}
	})
	var tx, ty float32
	if x, y, over := root.PointerPosition(); over {
		if bounds.W > 0 && bounds.H > 0 {
			tx = (x/bounds.W - 0.5) * 8
			ty = (y/bounds.H - 0.5) * 8
		}
	}
	mx := root.AnimateTo("mx", tx, ui.Snappy)
	my := root.AnimateTo("my", ty, ui.Snappy)
	root.Translate(mx, my).Children(func() { ui.Text(c, label).SingleLine() })
	if root.Clicked() {
		sound.Play(sound.Tap)
	}
	return root
}
