package input

import (
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// NumberFieldOptions tune a NumberField.
type NumberFieldOptions struct {
	// Min, Max and Step bound keyboard changes.
	Min, Max, Step float64
}

// NumberField has arrows, keys and rolling digits.
func NumberField(c *ui.Context, label string, value *float64, o NumberFieldOptions) ui.Element {
	if value == nil || !finite(*value) {
		panic("input: NumberField needs a finite value")
	}
	if o.Step == 0 {
		o.Step = 1
	}
	if o.Max == 0 {
		o.Max = 100
	}
	if !finite(o.Min) || !finite(o.Max) || !finite(o.Step) || o.Step <= 0 || o.Min > o.Max {
		panic("input: NumberField needs a valid range and step")
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).AlignItems(ui.Center).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
	root.Children(func() {
		decrease := ui.ButtonBase(c).Label("Decrease " + label).Padding(t.Space(2)).TextColor(t.TextMuted)
		kit.Press(decrease).Children(func() { ui.Text(c, "−") })
		kit.Numeric(c, i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), *value, i18n.NumberOptions{MaxFractionDigits: stepDigits(o.Step)})).Padding(t.Space(1), t.Space(2))
		increase := ui.ButtonBase(c).Label("Increase " + label).Padding(t.Space(2)).TextColor(t.TextMuted)
		kit.Press(increase).Children(func() { ui.Text(c, "+") })
		step := 0.0
		if decrease.Clicked() || root.Shortcut(0, ui.KeyDown) {
			step = -o.Step
		} else if increase.Clicked() || root.Shortcut(0, ui.KeyUp) {
			step = o.Step
		}
		if step != 0 {
			next := min(o.Max, max(o.Min, *value+step))
			if next != *value {
				*value = next
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("NumberField", "changed", "value", *value)
			}
		}
	})
	return root
}

// stepDigits returns the decimals needed to show a step.
func stepDigits(step float64) int {
	text := strconv.FormatFloat(step, 'f', -1, 64)
	_, fraction, found := strings.Cut(text, ".")
	if !found {
		return 0
	}
	return len(strings.TrimRight(fraction, "0"))
}
