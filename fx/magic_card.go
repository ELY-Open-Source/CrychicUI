package fx

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MagicCard is a card whose glow follows the pointer.
func MagicCard(c *ui.Context, content func()) ui.Element {
	if content == nil {
		panic("fx: MagicCard needs content")
	}
	t := c.Theme()
	root := kit.Card(c).Padding(t.Space(4)).Clip()
	x, y, a := follow(root)
	root.Border(1, root.AnimateColor("edge", t.Border.Mix(t.Accent, unit(root.Hovered())*0.6), ui.Smooth))
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		if a > 0 {
			glow(p, r.X+x, r.Y+y, t.Space(24), t.Accent, 0.3*a)
		}
	})).Children(content)
}
