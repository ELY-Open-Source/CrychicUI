package form

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NumberRangeOptions tune a NumberRange.
type NumberRangeOptions struct {
	// From and To name the two inputs.
	From, To string
}

// NumberRange edits low and high, pushing one to keep order.
func NumberRange(c *ui.Context, label string, low, high *float64, lo, hi, step float64, o NumberRangeOptions) ui.Element {
	finite := func(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }
	if label == "" || low == nil || high == nil || !finite(lo) || !finite(hi) || !finite(step) || lo >= hi || step <= 0 {
		panic(fmt.Sprintf("form: NumberRange %q needs values and finite lo < hi, step > 0", label))
	}
	if !(lo <= *low && *low <= *high && *high <= hi) {
		panic(fmt.Sprintf("form: NumberRange %q needs %v <= low %v <= high %v <= %v", label, lo, *low, *high, hi))
	}
	def(&o.From, "From")
	def(&o.To, "To")
	t := c.Theme()
	root := kit.Arrive(ui.Column(c).FillWidth(), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	return root.Children(func() {
		FieldGroup(c, label, FieldGroupOptions{}, func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				from := ui.NumberInput(c, low, lo, hi, step).Label(o.From).Grow(1)
				kit.Muted(c, "–")
				to := ui.NumberInput(c, high, lo, hi, step).Label(o.To).Grow(1)
				if from.Changed() && *low > *high {
					*high = *low
					c.AnimationFrame()
				}
				if to.Changed() && *high < *low {
					*low = *high
					c.AnimationFrame()
				}
				if from.Changed() || to.Changed() {
					sound.Play(sound.Tick)
					root.SetChanged()
					kit.Log("NumberRange", "changed")
				}
			})
		})
	})
}
