package motion

import "github.com/egoist/mygo/ui"

// Side is the edge content slides in from.
type Side int

// Sides; Start and End follow the text direction.
const (
	Bottom Side = iota
	Top
	Start
	End
)

// SlideInOptions tune a SlideIn.
type SlideInOptions struct {
	// From is the edge content comes from; Bottom when zero.
	From Side
}

// SlideIn slides content in from an edge on first build.
func SlideIn(c *ui.Context, o SlideInOptions, content func()) ui.Element {
	if content == nil {
		panic("motion: SlideIn needs content")
	}
	d := c.Theme().Space(6)
	var x, y float32
	switch o.From {
	case Bottom:
		y = d
	case Top:
		y = -d
	case Start:
		x = -d
	case End:
		x = d
	default:
		panic("motion: SlideIn has an unknown side")
	}
	if c.Direction() == ui.RightToLeft {
		x = -x
	}
	root := ui.Box(c)
	k := root.Keyframes("crychic.in", arrive)
	return root.Translate(x*(1-k), y*(1-k)).Opacity(k).Children(content)
}
