package button

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Outline is a bordered button with animated border color.
func Outline(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	b := textBase(c, label).TextColor(t.Text).Background(t.Background)
	border := t.Border
	if b.Pressed() {
		border = t.AccentPressed
	} else if b.Hovered() {
		border = t.Accent
	}
	b.Border(1, b.AnimateColor("border", border, ui.Snappy)).Children(func() { ui.Text(c, label).SingleLine() })
	kit.Clicked(b, sound.Tap)
	return b
}
