package color

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Picker is a color field opening MyGo's picker.
func Picker(c *ui.Context, label string, value *ui.Color) ui.Element {
	if label == "" || value == nil {
		panic("color: Picker needs a label and a color")
	}
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).AlignSelf(ui.Start).Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2.5)).
		Radius(t.Radius).Background(t.Background))
	open := ui.Local(b, "open", func() bool { return false })
	from := ui.Local(b, "from", func() ui.Color { return *value })
	opened := *open
	edge := t.Border
	if *open || b.Hovered() {
		edge = t.Accent
	}
	b.Border(1, b.AnimateColor("edge", edge, ui.Snappy)).Expanded(*open).Children(func() {
		swatch(c, *value, t.Space(5), t.Space(5))
		kit.Numeric(c, hexOf(*value)).ContentTransition(ui.ContentTransition{})
		chev := ui.Icon(c, icon.Must("chevron-down")).Size(16, 16)
		chev.Rotate(chev.AnimateTo("turn", 180*unit(*open), ui.Snappy))
	})
	if b.Clicked() {
		if *open = !*open; *open {
			*from = *value
		}
	}
	if popover(c, b, open, value) {
		b.SetChanged()
	}
	if opened && !*open && *value != *from {
		sound.Play(sound.Select)
		kit.Log("Picker", "picked")
	}
	return b
}
