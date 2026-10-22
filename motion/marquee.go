package motion

import "github.com/egoist/mygo/ui"

// MarqueeOptions tune a Marquee.
type MarqueeOptions struct {
	// Speed is DIPs per second; t.Space(10) when zero.
	Speed float32
}

// Marquee loops a row of content; hover pauses it.
func Marquee(c *ui.Context, o MarqueeOptions, content func()) ui.Element {
	if content == nil {
		panic("motion: Marquee needs content")
	}
	return marquee(c, false, positive("Marquee", "Speed", o.Speed, c.Theme().Space(10)), content)
}
