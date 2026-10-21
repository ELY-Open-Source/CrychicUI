package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ArrowButton slides its arrow right on hover.
func ArrowButton(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Background(t.Surface)
	root.Children(func() {
		ui.Text(c, label).SingleLine()
		arrow := ui.Icon(c, icon.Must("arrow-right")).Size(16, 16)
		arrow.Translate(arrow.AnimateTo("slide", 4*unit(root.Hovered()), ui.Snappy), 0)
	})
	if root.Clicked() {
		sound.Play(sound.Tap)
	}
	return root
}
