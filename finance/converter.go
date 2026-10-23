package finance

import (
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ConverterOptions are the converter labels.
type ConverterOptions struct{ Invalid string }

// Converter converts between currencies at a rate; cue Tick.
func Converter(c *ui.Context, amount *string, rate float64, from, to string, o ConverterOptions) ui.Element {
	if amount == nil || rate <= 0 || !finite(rate) || from == "" || to == "" {
		panic("finance: Converter has invalid data")
	}
	def(&o.Invalid, "Enter a number")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			in := ui.TextInput(c, amount).Label(from).Grow(1)
			if in.Changed() {
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("Converter", "edited")
			}
			ui.Text(c, from)
		})
		value, err := strconv.ParseFloat(strings.TrimSpace(*amount), 64)
		if err != nil || value < 0 || !finite(value) {
			ui.Text(c, o.Invalid).TextColor(t.Danger).FontSize(t.Rem(0.8))
			return
		}
		converted := value * rate
		if !finite(converted) {
			ui.Text(c, o.Invalid).TextColor(t.Danger).FontSize(t.Rem(0.8))
			return
		}
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			money(c, converted, to).Bold().ContentTransition(ui.ContentTransition{Numeric: true})
			ui.Text(c, to)
		})
	})
	return root
}
