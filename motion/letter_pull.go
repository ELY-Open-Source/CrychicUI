package motion

import (
	"time"

	"github.com/egoist/mygo/ui"
)

// LetterPull drops the letters of text in from above, staggered.
func LetterPull(c *ui.Context, text string) ui.Element {
	if text == "" {
		panic("motion: LetterPull needs text")
	}
	gs, w, h := letters(c, text)
	drop, rm := c.Theme().Space(4), still(c)
	root := ui.Box(c)
	return root.Children(func() {
		e := ui.Box(c.Key(text)).Size(w, h).Role(ui.RoleText).Label(text)
		start := ui.Local(e, "start", func() time.Time { return c.Now() })
		e.Draw(func(p *ui.Painter, r ui.Rect) {
			moving := false
			for i, g := range gs {
				k := float32(1)
				if !rm {
					el := p.Now().Sub(*start) - time.Duration(i)*30*time.Millisecond
					k = ui.EaseOut(min(max(float32(el)/float32(250*time.Millisecond), 0), 1))
				}
				moving = moving || k < 1
				if k == 0 {
					continue
				}
				g.span.Color = g.span.Color.Alpha(k)
				p.RichText(r.X+g.x, r.Y-(1-k)*drop, 0, g.span)
			}
			if moving {
				p.AnimationFrame()
			}
		})
	})
}
