package choice

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RangeSliderOptions tune a RangeSlider.
type RangeSliderOptions struct {
	// Format writes both values; "%.0f" when empty.
	Format string
}

// RangeSlider is two knobs with rolling values above.
func RangeSlider(c *ui.Context, label string, low, high *float64, lo, hi, step float64, o RangeSliderOptions) ui.Element {
	if low == nil || high == nil || !(step > 0) {
		panic("choice: " + who("RangeSlider", label) + " needs both values and a step")
	}
	sliderArgs("RangeSlider", label, low, lo, hi, step)
	sliderArgs("RangeSlider", label, high, lo, hi, step)
	if *low > *high || !onStep(*low, lo, step) || !onStep(*high, lo, step) {
		panic("choice: " + who("RangeSlider", label) + " needs low <= high, both on the step")
	}
	if o.Format == "" {
		o.Format = "%.0f"
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Row(c).Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Format, *low)).TextColor(t.TextMuted)
			ui.Spacer(c)
			kit.Numeric(c, fmt.Sprintf(o.Format, *high)).TextColor(t.TextMuted)
		})
		if kit.Changed(ui.RangeSlider(c, low, high, lo, reach(lo, hi, step), step).Label(label).FillWidth(), sound.Tick) {
			root.SetChanged()
			kit.Log("RangeSlider", "changed", "low", *low, "high", *high)
		}
	})
}
