package fx

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

const glyphs = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// MatrixRain drops columns of changing glyphs behind content.
func MatrixRain(c *ui.Context, content func()) ui.Element {
	root := scene(c, content)
	t := c.Theme()
	k := phase(c, root, 6*time.Second)
	tick := uint32(60 * k)
	col, size := t.Space(4), t.Rem(0.75)
	const trail = 8
	return root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		span := r.H + trail*size
		for i := range int(r.W / col) {
			head := r.Y + fract(rnd(i, 1)+k*float32(1+int(rnd(i, 2)*2)))*span
			for g := range trail {
				y := head - float32(g)*size
				if y < r.Y-size || y > r.Y+r.H {
					continue
				}
				ch := glyphs[int(rnd(i*trail+g, tick)*float32(len(glyphs)))]
				p.Text(r.X+float32(i)*col+col/4, y, string(ch), size, t.Success.Alpha(1-float32(g)/trail))
			}
		}
	})).Children(content)
}
