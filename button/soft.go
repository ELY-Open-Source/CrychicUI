package button

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Soft is an accent-tinted button with animated hover and press.
func Soft(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	b := textBase(c, label).TextColor(t.Accent)
	bg := t.Accent.Alpha(0.12)
	if b.Pressed() {
		bg = t.Accent.Alpha(0.24)
	} else if b.Hovered() {
		bg = t.Accent.Alpha(0.18)
	}
	b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() { ui.Text(c, label).SingleLine() })
	kit.Clicked(b, sound.Tap)
	return b
}
