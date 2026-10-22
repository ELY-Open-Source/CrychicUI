package color

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Eyedropper picks a color from a strip under the pointer.
func Eyedropper(c *ui.Context, label string, colors []ui.Color, value *ui.Color) ui.Element {
	if label == "" || len(colors) == 0 || value == nil {
		panic("color: Eyedropper needs a label, colors and a color")
	}
	t := c.Theme()
	w, h := t.Space(64), t.Space(10)
	cell := w / float32(len(colors))
	back, forth := ui.KeyLeft, ui.KeyRight
	if c.Direction() == ui.RightToLeft {
		back, forth = forth, back
	}
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(3)).AlignItems(ui.Center)
	return root.Children(func() {
		strip := kit.Press(ui.ButtonBase(c).Label(label).Size(w, h).Radius(t.Radius).Border(1, t.Border).Cursor(ui.CursorCrosshair))
		at := ui.Local(strip, "at", func() int { return 0 })
		*at = min(*at, len(colors)-1)
		switch {
		case strip.Shortcut(0, back):
			*at = max(0, *at-1)
		case strip.Shortcut(0, forth):
			*at = min(len(colors)-1, *at+1)
		}
		if x, _, over := strip.PointerPosition(); over {
			*at = min(len(colors)-1, max(0, int(x/cell)))
		}
		if strip.Clicked() {
			*value = colors[*at]
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("Eyedropper", "picked", "index", *at)
		}
		i, focus := *at, strip.Focused() || strip.Hovered()
		strip.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
			for k, col := range colors {
				p.Fill(ui.Rect{X: r.X + float32(k)*cell, Y: r.Y, W: cell + 0.5, H: r.H}, col, 0)
			}
			if focus {
				p.Stroke(ui.Rect{X: r.X + float32(i)*cell + 1, Y: r.Y + 1, W: cell - 2, H: r.H - 2}, t.Background, t.Radius/2, 2)
			}
		}))
		loupe := swatch(c, colors[i], t.Space(12), t.Space(12)).Radius(999)
		k := 1 + 0.1*loupe.AnimateTo("hot", unit(focus), ui.Bouncy)
		loupe.Scale(k, k)
		swatch(c, *value, t.Space(6), t.Space(6))
	})
}
