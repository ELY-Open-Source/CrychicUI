package datetime

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// TimeRangeSliderOptions tune a TimeRangeSlider.
type TimeRangeSliderOptions struct {
	// Step divides a day; 30 minutes when zero.
	Step time.Duration
}

// TimeRangeSlider is two thumbs over one day.
func TimeRangeSlider(c *ui.Context, label string, from, to *time.Duration, o TimeRangeSliderOptions) ui.Element {
	const day = 24 * time.Hour
	if o.Step == 0 {
		o.Step = 30 * time.Minute
	}
	if label == "" || from == nil || to == nil || o.Step < time.Minute || day%o.Step != 0 ||
		*from < 0 || *from > *to || *to > day || *from%o.Step != 0 || *to%o.Step != 0 {
		panic(fmt.Sprintf("datetime: TimeRangeSlider %q needs 0 <= from <= to <= 24h on a %v step", label, o.Step))
	}
	t := c.Theme()
	loc := locale(c)
	write := func(d time.Duration) string {
		return i18n.FormatDate(loc, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Add(d), i18n.DateOptions{Skeleton: "jm"})
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			kit.Numeric(c, write(*from)).FontWeight(600)
			kit.Muted(c, "–")
			kit.Numeric(c, write(*to)).FontWeight(600)
		})
		lo, hi := from.Minutes(), to.Minutes()
		if kit.Changed(ui.RangeSlider(c, &lo, &hi, 0, day.Minutes(), o.Step.Minutes()).Label(label).FillWidth(), sound.Tick) {
			*from, *to = onStep(lo, o.Step), onStep(hi, o.Step)
			root.SetChanged()
			kit.Log("TimeRangeSlider", "changed")
		}
		ui.Row(c).FillWidth().Justify(ui.SpaceBetween).Children(func() {
			for h := 0; h <= 24; h += 6 {
				kit.Muted(c, write(time.Duration(h)*time.Hour)).FontSize(t.Rem(0.75))
			}
		})
	})
}

// onStep turns slider minutes into the nearest step.
func onStep(minutes float64, step time.Duration) time.Duration {
	return time.Duration(math.Round(minutes*float64(time.Minute)/float64(step))) * step
}
