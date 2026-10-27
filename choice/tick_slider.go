package choice

import (
	"math"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TickSliderOptions tune a TickSlider.
type TickSliderOptions struct {
	// Labels name each tick; the values when empty.
	Labels []string
}

// TickSlider is a stepped slider with tick marks and labels.
func TickSlider(c *ui.Context, label string, value *float64, lo, hi, step float64, o TickSliderOptions) ui.Element {
	sliderArgs("TickSlider", label, value, lo, hi, step)
	steps := (hi - lo) / step
	n := int(math.Round(steps)) + 1
	if !(step > 0) || math.IsInf(step, 0) || math.Abs(steps-math.Round(steps)) > 1e-9 || n < 2 || n > 21 {
		panic("choice: " + who("TickSlider", label) + " needs a step dividing the range into at most 20")
	}
	if len(o.Labels) == 0 {
		for i := range n {
			o.Labels = append(o.Labels, strconv.FormatFloat(lo+step*float64(i), 'f', -1, 64))
		}
	}
	if len(o.Labels) != n {
		panic("choice: " + who("TickSlider", label) + " needs one label per tick")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		s := ui.SliderBase(c, value, lo, reach(lo, hi, step)).Step(step).Label(label).FillWidth()
		rail(c, s, fraction(*value, lo, hi), false, nil)
		if kit.Changed(s, sound.Tick) {
			root.SetChanged()
			kit.Log("TickSlider", "changed", "value", *value)
		}
		at := int(math.Round((*value - lo) / step))
		ui.Box(c).FillWidth().Height(t.Space(7)).PaddingX(t.Space(2.5)).Children(func() {
			ui.Box(c).Fill().Children(func() {
				for i, name := range o.Labels {
					c.Key(i)
					spot := ui.Box(c).Absolute().StartPercent(float32(i)*100/float32(n-1)).Top(0).Size(0, 0)
					spot.Children(func() {
						ui.Box(c).Attach(ui.AnchorTop, ui.AnchorTop).Size(1, t.Space(1.5)).Background(t.TextMuted.Alpha(0.6))
						text := ui.Text(c, name).SingleLine().FontSize(t.Rem(0.8)).Attach(ui.AnchorTop, ui.AnchorTop).Top(t.Space(2))
						fg := t.TextMuted
						if i == at {
							fg = t.Accent
						}
						text.TextColor(text.AnimateColor("fg", fg, quick))
					})
				}
			})
		})
	})
}
