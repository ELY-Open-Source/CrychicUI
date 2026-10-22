package color

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Well is a swatch button that opens MyGo's picker.
func Well(c *ui.Context, label string, value *ui.Color) ui.Element {
	if label == "" || value == nil {
		panic("color: Well needs a label and a color")
	}
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Tooltip(label).Size(t.Space(10), t.Space(10)).Shrink(0).Padding(t.Space(1)).
		Radius(t.Radius).Background(t.Background))
	open := ui.Local(b, "open", func() bool { return false })
	edge := t.Border
	if *open || b.Hovered() {
		edge = t.Accent
	}
	b.Border(1, b.AnimateColor("edge", edge, ui.Snappy)).Expanded(*open).Children(func() {
		swatch(c, *value, t.Space(8), t.Space(8))
	})
	if b.Clicked() {
		if *open = !*open; *open {
			sound.Play(sound.Open)
			kit.Log("Well", "open")
		}
	}
	if popover(c, b, open, value) {
		b.SetChanged()
	}
	return b
}
