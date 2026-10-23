package media

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ImageStackOptions tune an ImageStack.
type ImageStackOptions struct {
	// Size is a card's width; zero gives the default.
	Size float32
}

// ImageStack fans images; a click brings the next forward.
func ImageStack(c *ui.Context, label string, images []*ui.Bitmap, top *int, o ImageStackOptions) ui.Element {
	pictures("ImageStack", label, images)
	chosen("ImageStack", label, top, 0, len(images))
	t := c.Theme()
	size("ImageStack", label, &o.Size, t.Space(40))
	n, sz := len(images), o.Size
	root := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(8), t.Space(14)).Center())
	if root.Clicked() && n > 1 {
		*top = (*top + 1) % n
		root.SetChanged()
		sound.Play(sound.Swipe)
		kit.Log("ImageStack", "cycled", "index", *top)
	}
	shown := *top
	return root.Children(func() {
		ui.Grid(c).Columns(1).Children(func() {
			for k := n - 1; k >= 0; k-- {
				i := (shown + k) % n
				c.Key(i)
				side, fan := float32(1), float32((k+1)/2)
				if k%2 == 1 {
					side = -1
				}
				card := ui.Box(c).ColumnStart(1).RowStart(1).Size(sz, sz*0.75).Radius(t.Radius).Clip().Border(1, t.Border).
					Background(t.Surface).Shadow(0, t.Space(1), t.Space(3), 0, kit.Shade(0.12))
				card.AnimatePose(ui.Pose{Rotate: side * fan * 6, X: side * fan * t.Space(5), Scale: 1 - float32(k)*0.04}, ui.Snappy).
					Children(func() { ui.Image(c, images[i]).Fill().Fit(ui.Cover) })
			}
		})
	})
}
