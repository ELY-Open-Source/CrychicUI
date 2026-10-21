package button

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Ghost is transparent until hovered or pressed.
func Ghost(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	b := textBase(c, label).TextColor(t.Text)
	bg := t.Background.Alpha(0)
	if b.Pressed() {
		bg = t.SurfacePressed
	} else if b.Hovered() {
		bg = t.SurfaceHover
	}
	b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() { ui.Text(c, label).SingleLine() })
	kit.Clicked(b, sound.Tap)
	return b
}
