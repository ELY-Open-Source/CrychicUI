package color

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// OklchFieldsOptions tune OklchFields.
type OklchFieldsOptions struct {
	// Names label the lightness, chroma and hue inputs.
	Names [3]string
	// Gamut notes a color sRGB screens cannot show.
	Gamut string
}

// OklchFields edits lightness, chroma and hue, noting sRGB's edge.
func OklchFields(c *ui.Context, label string, value *ui.Color, o OklchFieldsOptions) ui.Element {
	if label == "" || value == nil {
		panic("color: OklchFields needs a label and a color")
	}
	if o.Names == ([3]string{}) {
		o.Names = [3]string{"L", "C", "H"}
	}
	def(&o.Gamut, "Outside sRGB; most screens show the nearest color")
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(1.5)).Label(label)
	vals, made := kept(root, value, func(v ui.Color) [3]float64 {
		l, ch, h := oklch(v)
		return [3]float64{l * 100, ch, h}
	})
	return root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.End).Children(func() {
			swatch(c, *value, t.Space(9), t.Space(9))
			if triple(c, o.Names, vals, [3]float64{100, 0.4, 360}, [3]float64{1, 0.001, 1}) {
				*value = ui.Oklch(float32(vals[0]/100), float32(vals[1]), float32(vals[2])).Alpha(float32(value.A) / 255)
				*made = *value
				root.SetChanged()
				kit.Log("OklchFields", "changed")
			}
		})
		if !inSRGB(vals[0]/100, vals[1], vals[2]) {
			ui.Row(c).Key("gamut").Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.Warning).Transition(kit.Enter).Children(func() {
				ui.Icon(c, icon.Must("triangle-alert")).Size(16, 16)
				ui.Text(c, o.Gamut).FontSize(t.Rem(0.85))
			})
		}
	})
}
