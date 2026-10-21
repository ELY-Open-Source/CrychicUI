package choice

import (
	"math"

	"github.com/egoist/mygo/ui"
)

// VerticalSliderOptions tune a VerticalSlider.
type VerticalSliderOptions struct {
	// Height is the slider's length; 40 spacing units when zero.
	Height float32
	// Step snaps the value; free when zero.
	Step float64
}

// VerticalSlider is a slider rising from lo at the bottom.
func VerticalSlider(c *ui.Context, label string, value *float64, lo, hi float64, o VerticalSliderOptions) ui.Element {
	sliderArgs("VerticalSlider", label, value, lo, hi, o.Step)
	if o.Height == 0 {
		o.Height = c.Theme().Space(40)
	}
	if !(o.Height > 0) || math.IsInf(float64(o.Height), 0) {
		panic("choice: VerticalSlider needs a positive height")
	}
	s := ui.SliderBase(c, value, lo, hi).Vertical().Height(o.Height).Label(label).Step(o.Step)
	rail(c, s, fraction(*value, lo, hi), true, nil)
	released(s, "VerticalSlider", *value)
	return s
}
