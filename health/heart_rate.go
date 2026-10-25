package health

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// HeartRateOptions word the card; Zones name resting, fat burn, cardio.
type HeartRateOptions struct {
	Title, Unit, Range string
	Zones              [3]string
}

// HeartRate beats at bpm over a line of recent samples.
func HeartRate(c *ui.Context, bpm int, samples []int, o HeartRateOptions) ui.Element {
	def(&o.Title, "Heart rate")
	if bpm < 20 || bpm > 250 || len(samples) < 2 {
		panic(fmt.Sprintf("health: HeartRate %q needs 20 to 250 bpm and two or more samples, got %d", o.Title, bpm))
	}
	lo, hi := samples[0], samples[0]
	values := make([]float64, len(samples))
	for i, s := range samples {
		if s < 20 || s > 250 {
			panic(fmt.Sprintf("health: HeartRate %q sample %d is %d bpm", o.Title, i, s))
		}
		lo, hi = min(lo, s), max(hi, s)
		values[i] = float64(s)
	}
	def(&o.Unit, "bpm")
	def(&o.Range, "%d–%d bpm")
	for i, z := range [...]string{"Resting", "Fat burn", "Cardio"} {
		def(&o.Zones[i], z)
	}
	t := c.Theme()
	zone, ink := 0, t.Success
	switch {
	case bpm >= 140:
		zone, ink = 2, t.Danger
	case bpm >= 100:
		zone, ink = 1, t.Warning
	}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			beat := ui.Box(c).Size(t.Space(11), t.Space(11)).Shrink(0).Center().Radius(t.Radius).Background(t.Danger.Alpha(0.12)).TextColor(t.Danger)
			beat.Children(func() {
				g := glyph(c, "heart", 22)
				v := float64(g.Loop("beat", time.Duration(float64(time.Minute)/float64(bpm)), ui.Linear))
				k := 1 + 0.2*float32(math.Max(0, 1-v*5))
				g.Scale(k, k)
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				kit.Muted(c, o.Title).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.End).Children(func() {
					kit.Numeric(c, fmt.Sprint(bpm)).FontSize(t.Rem(1.6)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
					kit.Muted(c, o.Unit).FontSize(t.Rem(0.78)).Shrink(0)
				})
			})
			kit.Chip(c, o.Zones[zone], ink)
		})
		trend(c, values, float64(lo)-5, float64(hi)+5, t.Danger, t.Space(20), math.NaN())
		kit.Muted(c, fmt.Sprintf(o.Range, lo, hi)).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
	})
	return root
}
