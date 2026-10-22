package motion

import (
	"time"
	"unicode"

	"github.com/egoist/mygo/ui"
)

// ScrambleOptions tune a Scramble.
type ScrambleOptions struct {
	// Glyphs are drawn while letters scramble; A–Z and 0–9.
	Glyphs string
}

var settle = []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.Linear}}

// Scramble settles random letters on text, left to right.
func Scramble(c *ui.Context, text string, o ScrambleOptions) ui.Element {
	if text == "" {
		panic("motion: Scramble needs text")
	}
	if o.Glyphs == "" {
		o.Glyphs = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}
	pool, runes := []rune(o.Glyphs), []rune(text)
	root := ui.Box(c)
	return root.Children(func() {
		w := ui.Box(c.Key(text))
		k := w.Keyframes("crychic.settle", settle)
		tick := int(k * 30)
		out := make([]rune, len(runes))
		for i, r := range runes {
			out[i] = r
			if !unicode.IsSpace(r) && k < float32(i+1)/float32(len(runes)) {
				out[i] = pool[(i*7+tick*13)%len(pool)]
			}
		}
		w.Children(func() { ui.Text(c, string(out)).Label(text) })
	})
}
