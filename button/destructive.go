package button

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Destructive is a danger button with a delete cue.
func Destructive(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	b := textBase(c, label).TextColor(t.AccentText)
	bg := t.Danger
	if b.Pressed() {
		bg = t.Danger.Mix(t.Text, 0.16)
	} else if b.Hovered() {
		bg = t.Danger.Mix(t.Text, 0.08)
	}
	b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Children(func() { ui.Text(c, label).SingleLine() })
	kit.Clicked(b, sound.Delete)
	return b
}
