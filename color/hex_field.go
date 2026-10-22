package color

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HexFieldOptions tune a HexField.
type HexFieldOptions struct {
	// Invalid shows under text that is no hex color.
	Invalid string
}

// HexField edits a color as hex beside its swatch.
func HexField(c *ui.Context, label string, value *ui.Color, o HexFieldOptions) ui.Element {
	if label == "" || value == nil {
		panic("color: HexField needs a label and a color")
	}
	def(&o.Invalid, "Use #RGB, #RRGGBB or #RRGGBBAA")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.End)
	text := ui.Local(root, "text", func() string { return hexOf(*value) })
	synced := ui.Local(root, "synced", func() ui.Color { return *value })
	if *value != *synced {
		*text, *synced = hexOf(*value), *value
	}
	return root.Children(func() {
		swatch(c, *value, t.Space(9), t.Space(9))
		bad := ""
		ui.Field(c, label, func() {
			in := ui.TextInput(c, text).Width(t.Space(36))
			col, ok := parseHex(*text)
			if !ok {
				bad = o.Invalid
			} else if in.Changed() && col != *value {
				*value, *synced = col, col
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("HexField", "changed")
			}
		}).Error(bad)
	})
}
