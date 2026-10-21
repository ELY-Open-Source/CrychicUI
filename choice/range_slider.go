package choice

import (
	"fmt"
	"math"
	"strconv"
	"strings"

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
		panic("choice: RangeSlider needs both values and a step")
	}
	sliderArgs("RangeSlider", label, low, lo, hi, step)
	sliderArgs("RangeSlider", label, high, lo, hi, step)
	if *low > *high || !onStep(*low, lo, step) || !onStep(*high, lo, step) {
		panic("choice: RangeSlider needs low <= high, both on the step")
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
		if kit.Changed(ui.RangeSlider(c, low, high, lo, hi, step).Label(label).FillWidth(), sound.Tick) {
			root.SetChanged()
			kit.Log("RangeSlider", "changed", "low", *low, "high", *high)
		}
	})
}

// onStep reports whether MyGo's slider snap keeps v.
func onStep(v, lo, step float64) bool {
	s := lo + math.Round((v-lo)/step)*step
	if d := max(decimals(step), decimals(lo)); d < 15 {
		r, err := strconv.ParseFloat(strconv.FormatFloat(s, 'f', d, 64), 64)
		if err != nil {
			panic(err)
		}
		s = r
	}
	return s == v
}

// decimals counts the digits after the point, as MyGo does.
func decimals(v float64) int {
	mant, exp, _ := strings.Cut(strconv.FormatFloat(math.Abs(v), 'e', -1, 64), "e")
	e, err := strconv.Atoi(exp)
	if err != nil {
		panic(err)
	}
	digits := 0
	if _, frac, ok := strings.Cut(mant, "."); ok {
		digits = len(frac)
	}
	return max(digits-e, 0)
}
