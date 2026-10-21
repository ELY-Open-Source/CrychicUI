package button

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Solid is an accent button with animated hover and press.
func Solid(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	b := textBase(c, label).TextColor(t.AccentText)
	bg := t.Accent
	if b.Pressed() {
		bg = t.AccentPressed
	} else if b.Hovered() {
		bg = t.AccentHover
	}
	b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() { ui.Text(c, label).SingleLine() })
	kit.Clicked(b, sound.Tap)
	return b
}
