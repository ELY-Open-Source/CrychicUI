package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CornerSpot is the parent corner a Corner sits in.
type CornerSpot int

// Corner spots.
const (
	CornerTopLeft CornerSpot = iota
	CornerTopRight
	CornerBottomLeft
	CornerBottomRight
)

// Corner places content at a parent's corner; content pops in.
func Corner(c *ui.Context, spot CornerSpot, content func()) ui.Element {
	in := c.Theme().Space(2)
	b := ui.Box(c).Absolute()
	switch spot {
	case CornerTopLeft:
		b.Top(in).Left(in)
	case CornerTopRight:
		b.Top(in).Right(in)
	case CornerBottomLeft:
		b.Bottom(in).Left(in)
	case CornerBottomRight:
		b.Bottom(in).Right(in)
	default:
		panic("layout: Corner spot out of range")
	}
	return kit.Arrive(b, ui.Pose{Scale: 0.9, Fade: 1}, 0).Children(content)
}
