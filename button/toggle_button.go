package button

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ToggleButton stays pressed with an animated tint.
func ToggleButton(c *ui.Context, label string, on *bool) ui.Element {
	if on == nil {
		panic("button: ToggleButton needs state")
	}
	t := c.Theme()
	b := textBase(c, label).Checked(*on).Border(1, t.Border)
	bg, fg := t.Surface, t.Text
	if *on {
		bg, fg = t.Accent.Alpha(0.16), t.Accent
	} else if b.Pressed() {
		bg = t.SurfacePressed
	} else if b.Hovered() {
		bg = t.SurfaceHover
	}
	b.Background(b.AnimateColor("bg", bg, ui.Snappy)).TextColor(b.AnimateColor("fg", fg, ui.Snappy)).
		Children(func() { ui.Text(c, label).SingleLine() })
	if kit.Clicked(b, sound.Toggle) {
		*on = !*on
		b.SetChanged()
		kit.Log("ToggleButton", "toggled", "on", *on)
	}
	return b
}
