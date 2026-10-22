package color

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// HSLFieldsOptions tune HSLFields.
type HSLFieldsOptions struct {
	// Names label the hue, saturation and lightness inputs.
	Names [3]string
}

// HSLFields edits hue, saturation and lightness; the swatch follows.
func HSLFields(c *ui.Context, label string, value *ui.Color, o HSLFieldsOptions) ui.Element {
	if label == "" || value == nil {
		panic("color: HSLFields needs a label and a color")
	}
	if o.Names == ([3]string{}) {
		o.Names = [3]string{"H", "S", "L"}
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).AlignItems(ui.End).Label(label)
	vals, made := kept(root, value, func(v ui.Color) [3]float64 {
		h, s, l := hsl(v)
		return [3]float64{h, s * 100, l * 100}
	})
	return root.Children(func() {
		swatch(c, *value, t.Space(9), t.Space(9))
		if triple(c, o.Names, vals, [3]float64{360, 100, 100}, [3]float64{1, 1, 1}) {
			*value = fromHSL(vals[0], vals[1]/100, vals[2]/100, value.A)
			*made = *value
			root.SetChanged()
			kit.Log("HSLFields", "changed")
		}
	})
}
