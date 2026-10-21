package choice

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// ValueSliderOptions tune a ValueSlider.
type ValueSliderOptions struct {
	// Step snaps the value; free when zero.
	Step float64
	// Format writes the bubble's value; "%.0f" when empty.
	Format string
}

// ValueSlider is a slider whose value bubble follows the knob.
func ValueSlider(c *ui.Context, label string, value *float64, lo, hi float64, o ValueSliderOptions) ui.Element {
	sliderArgs("ValueSlider", label, value, lo, hi, o.Step)
	if o.Format == "" {
		o.Format = "%.0f"
	}
	t := c.Theme()
	s := ui.SliderBase(c, value, lo, hi).Label(label).Step(o.Step)
	show := s.Hovered() || s.Pressed() || s.Focused()
	text := fmt.Sprintf(o.Format, *value)
	rail(c, s, fraction(*value, lo, hi), false, func(ui.Element) {
		bubble := ui.Text(c, text).SingleLine().Attach(ui.AnchorTop, ui.AnchorBottom).Top(-t.Space(1.5)).
			Padding(t.Space(0.5), t.Space(2)).Radius(t.Radius).Background(t.Text).TextColor(t.Background).FontSize(t.Rem(0.85))
		k := bubble.AnimateTo("shown", unit(show), ui.Snappy)
		bubble.Scale(max(k, 0), max(k, 0)).Opacity(min(max(k, 0), 1))
	})
	released(s, "ValueSlider", *value)
	return s
}
