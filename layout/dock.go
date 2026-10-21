package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DockSide is the parent edge a Dock clings to.
type DockSide int

// Dock sides.
const (
	DockLeft DockSide = iota
	DockRight
	DockTop
	DockBottom
)

var dockSpring = ui.Snappy

// Dock slides an open panel from its parent's side.
func Dock(c *ui.Context, side DockSide, open bool, content func()) ui.Element {
	t := c.Theme()
	root, from := ui.Box(c).Absolute().PassThrough(), ui.Motion{Opacity: 1}
	w, h := t.Rem(16), t.Rem(12)
	switch side {
	case DockLeft:
		root.Top(0).Bottom(0).Left(0).Width(w)
		from.X = -w
	case DockRight:
		root.Top(0).Bottom(0).Right(0).Width(w)
		from.X = w
	case DockTop:
		root.Top(0).Left(0).Right(0).Height(h)
		from.Y = -h
	case DockBottom:
		root.Bottom(0).Left(0).Right(0).Height(h)
		from.Y = h
	default:
		panic("layout: Dock side out of range")
	}
	return root.Children(func() {
		if !open {
			return
		}
		panel := ui.Column(c).Key("panel").Gap(t.Space(2)).Padding(t.Space(3)).Background(t.Background).
			Border(1, t.Border).Shadow(0, 4, 12, 0, kit.Shade(0.1))
		if from.X != 0 {
			panel.Width(w).FillHeight()
		} else {
			panel.Height(h).FillWidth()
		}
		panel.Transition(ui.ElementTransition{Animation: &dockSpring, Enter: &from, Exit: &from}).Children(content)
	})
}
