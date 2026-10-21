package button

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// IconLabel pairs an icon and text with motion.
func IconLabel(c *ui.Context, label, iconName string) ui.Element {
	t := c.Theme()
	b := textBase(c, label).TextColor(t.Text)
	bg := t.Surface
	if b.Pressed() {
		bg = t.SurfacePressed
	} else if b.Hovered() {
		bg = t.SurfaceHover
	}
	b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Border(1, t.Border).Children(func() {
		glyph := ui.Icon(c, icon.Must(iconName)).Size(16, 16)
		glyph.Translate(glyph.AnimateTo("nudge", 2*unit(b.Hovered()), ui.Snappy), 0)
		ui.Text(c, label).SingleLine()
	})
	kit.Clicked(b, sound.Tap)
	return b
}
