package color

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// popover shows MyGo's picker for value below anchor.
func popover(c *ui.Context, anchor ui.Element, open *bool, value *ui.Color) (changed bool) {
	ui.Popover(c, anchor, open, func() { changed = ui.ColorPicker(c, value).Changed() })
	return changed
}

// triple lays out three named number inputs.
func triple(c *ui.Context, names [3]string, vals *[3]float64, his, steps [3]float64) (changed bool) {
	for i := range vals {
		c.Key(i)
		ui.Field(c, names[i], func() {
			changed = kit.Changed(ui.NumberInput(c, &vals[i], 0, his[i], steps[i]), sound.Tick) || changed
		})
	}
	return changed
}

// kept holds numbers until the color changes from outside.
func kept(root ui.Element, value *ui.Color, split func(ui.Color) [3]float64) (*[3]float64, *ui.Color) {
	vals := ui.Local(root, "vals", func() [3]float64 { return split(*value) })
	made := ui.Local(root, "made", func() ui.Color { return *value })
	if *value != *made {
		*vals, *made = split(*value), *value
	}
	return vals, made
}
