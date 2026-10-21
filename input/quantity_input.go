package input

import (
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// QuantityInputOptions tune a QuantityInput.
type QuantityInputOptions struct {
	// Min and Max bound the quantity.
	Min, Max int
}

// QuantityInput is minus, a rolling value and plus.
func QuantityInput(c *ui.Context, label string, value *int, o QuantityInputOptions) ui.Element {
	if value == nil {
		panic("input: QuantityInput needs a value")
	}
	if o.Max == 0 {
		o.Max = 99
	}
	if o.Min > o.Max || *value < o.Min || *value > o.Max {
		panic("input: QuantityInput needs a valid range")
	}
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).AlignItems(ui.Center).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
	return root.Children(func() {
		decrease := ui.ButtonBase(c).Label("Decrease " + label).Padding(t.Space(2)).TextColor(t.TextMuted)
		kit.Press(decrease).Children(func() { ui.Text(c, "−") })
		kit.Numeric(c, strconv.Itoa(*value)).Padding(t.Space(1), t.Space(2))
		increase := ui.ButtonBase(c).Label("Increase " + label).Padding(t.Space(2)).TextColor(t.TextMuted)
		kit.Press(increase).Children(func() { ui.Text(c, "+") })
		step := 0
		if decrease.Clicked() || root.Shortcut(0, ui.KeyDown) {
			step = -1
		} else if increase.Clicked() || root.Shortcut(0, ui.KeyUp) {
			step = 1
		}
		if step != 0 {
			if step > 0 && *value == o.Max || step < 0 && *value == o.Min {
				return
			}
			next := *value + step
			if next != *value {
				*value = next
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("QuantityInput", "changed", "value", *value)
			}
		}
	})
}
