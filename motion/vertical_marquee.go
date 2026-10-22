package motion

import "github.com/egoist/mygo/ui"

// VerticalMarqueeOptions tune a VerticalMarquee.
type VerticalMarqueeOptions struct {
	// Height is the window shown; t.Space(40) when zero.
	Height float32
	// Speed is DIPs per second; t.Space(6) when zero.
	Speed float32
}

// VerticalMarquee loops a column of content; hover pauses it.
func VerticalMarquee(c *ui.Context, o VerticalMarqueeOptions, content func()) ui.Element {
	if content == nil {
		panic("motion: VerticalMarquee needs content")
	}
	t := c.Theme()
	h := positive("VerticalMarquee", "Height", o.Height, t.Space(40))
	return marquee(c, true, positive("VerticalMarquee", "Speed", o.Speed, t.Space(6)), content).Height(h)
}
