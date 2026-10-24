package finance

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ConverterOptions word errors; Rate takes from, rate and to.
type ConverterOptions struct{ Invalid, Rate string }

// Converter turns an amount into another currency; cue Tick.
func Converter(c *ui.Context, amount *string, rate float64, from, to string, o ConverterOptions) ui.Element {
	if amount == nil || rate <= 0 || !finite(rate) || from == "" || to == "" {
		panic(fmt.Sprintf("finance: Converter %s to %s needs a positive rate and bound text, got %v", from, to, rate))
	}
	def(&o.Invalid, "Enter a number")
	def(&o.Rate, "1 %s = %s %s")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	field := func() ui.Element {
		return ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Surface)
	}
	root.Children(func() {
		field().Children(func() {
			in := ui.TextInput(c, amount).Label(from).Grow(1).MinWidth(0)
			if in.Changed() {
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("Converter", "edited")
			}
			ui.Text(c, from).Bold().FontSize(t.Rem(0.85)).Shrink(0)
		})
		ui.Row(c).FillWidth().Justify(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(8), t.Space(8)).Center().Radius(999).Border(1, t.Border).Background(t.Background).TextColor(t.TextMuted).
				Children(func() { ui.Icon(c, icon.Must("arrow-down")).Size(14, 14) })
		})
		value, err := strconv.ParseFloat(strings.TrimSpace(*amount), 64)
		converted := value * rate
		field().Children(func() {
			if err != nil || value < 0 || !finite(value) || !finite(converted) {
				ui.Icon(c, icon.Must("triangle-alert")).Size(14, 14).TextColor(t.Danger)
				ui.Text(c, o.Invalid).FontSize(t.Rem(0.85)).TextColor(t.Danger).Grow(1).MinWidth(0)
				return
			}
			money(c, converted, to).Bold().FontSize(t.Rem(1.3)).Grow(1).MinWidth(0).ContentTransition(ui.ContentTransition{Numeric: true})
		})
		kit.Muted(c, fmt.Sprintf(o.Rate, from, strconv.FormatFloat(rate, 'f', 4, 64), to)).FontSize(t.Rem(0.75)).Shrink(0)
	})
	return root
}
