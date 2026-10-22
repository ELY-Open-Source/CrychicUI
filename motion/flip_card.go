package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FlipCard shows front or back; a click turns it over.
func FlipCard(c *ui.Context, label string, flipped *bool, front, back func()) ui.Element {
	if label == "" || flipped == nil || front == nil || back == nil {
		panic("motion: FlipCard needs a label, a flipped pointer and two faces")
	}
	t := c.Theme()
	root := ui.ButtonBase(c).Label(label)
	if kit.Clicked(root, sound.Swipe) {
		*flipped = !*flipped
		root.SetChanged()
		kit.Log("FlipCard", "flipped", "back", *flipped)
	}
	var target float32
	if *flipped {
		target = 180
	}
	turn := root.AnimateTo("turn", target, ui.Snappy)
	depth := t.Space(200)
	return kit.Press(root.Checked(*flipped)).Children(func() {
		ui.Grid(c).Columns(1).Children(func() {
			kit.Card(c).Padding(t.Space(4)).ColumnStart(1).RowStart(1).RotateY(turn).Perspective(depth).BackfaceHidden().Children(front)
			kit.Card(c).Padding(t.Space(4)).ColumnStart(1).RowStart(1).RotateY(turn + 180).Perspective(depth).BackfaceHidden().Children(back)
		})
	})
}
