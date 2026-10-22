package color

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RGBFieldsOptions tune RGBFields.
type RGBFieldsOptions struct {
	// Names label the red, green and blue inputs.
	Names [3]string
}

// RGBFields edits red, green and blue; the swatch follows.
func RGBFields(c *ui.Context, label string, value *ui.Color, o RGBFieldsOptions) ui.Element {
	if label == "" || value == nil {
		panic("color: RGBFields needs a label and a color")
	}
	if o.Names == ([3]string{}) {
		o.Names = [3]string{"R", "G", "B"}
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.End).Label(label)
	vals, made := kept(root, value, func(v ui.Color) [3]float64 { return [3]float64{float64(v.R), float64(v.G), float64(v.B)} })
	return root.Children(func() {
		swatch(c, *value, t.Space(9), t.Space(9))
		if triple(c, o.Names, vals, [3]float64{255, 255, 255}, [3]float64{1, 1, 1}) {
			b := func(v float64) uint8 { return uint8(v + 0.5) }
			*value = ui.Color{R: b(vals[0]), G: b(vals[1]), B: b(vals[2]), A: value.A}
			*made = *value
			root.SetChanged()
			kit.Log("RGBFields", "changed")
		}
	})
}
