package feedback

import (
	"math"
	"testing"
	"time"

	"github.com/egoist/mygo/ui"
)

func TestProgressBar(t *testing.T) {
	v := 0.25
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(400).Children(func() { ProgressBar(c, "Upload", v) })
	}, 400, 100)
	v = 0.75
	tt.Frame()
	settles(t, tt, "ProgressBar fill")
	th := ui.LightTheme()
	img := tt.Image()
	at := func(x int) bool {
		p := img.RGBAAt(x, 3)
		return p.R == th.Accent.R && p.G == th.Accent.G && p.B == th.Accent.B
	}
	if !at(290) || at(310) {
		t.Errorf("fill does not end at 75%%: 290=%v 310=%v", at(290), at(310))
	}
	for _, bad := range []float64{-0.1, 1.5, math.NaN()} {
		panics(t, "bad value", func() { render(func(c *ui.Context) { ProgressBar(c, "x", bad) }) })
	}
}

func TestStripedAndSegments(t *testing.T) {
	done := 1
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(8).Children(func() {
			StripedProgress(c, "Sync", 0.5)
			SegmentProgress(c, "Steps", 4, done)
		})
	}, 400, 100)
	loops(t, tt, "stripes")
	panics(t, "done past n", func() { render(func(c *ui.Context) { SegmentProgress(c, "x", 3, 4) }) })
}

func TestSegmentProgressStaggers(t *testing.T) {
	done := 0
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Width(400).Children(func() { SegmentProgress(c, "Steps", 4, done) })
	}, 400, 40)
	done = 4
	tt.Frame()
	settles(t, tt, "segments")
	th := ui.LightTheme()
	p := tt.Image().RGBAAt(390, 3)
	if p.R != th.Accent.R || p.B != th.Accent.B {
		t.Errorf("last segment not filled: %v", p)
	}
}

func TestRingAndGauge(t *testing.T) {
	v := 0.4
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Gap(8).Children(func() {
			ProgressRing(c, "Done", v, ProgressRingOptions{})
			Gauge(c, "Speed", v*100, 0, 100, GaugeOptions{Format: "%.0f km/h"})
		})
	}, 400, 200)
	find(t, tt, "40%")
	find(t, tt, "40 km/h")
	v = 0.9
	tt.Frame()
	settles(t, tt, "ring and needle")
	find(t, tt, "90%")
	panics(t, "gauge out of range", func() {
		render(func(c *ui.Context) { Gauge(c, "x", 120, 0, 100, GaugeOptions{}) })
	})
	panics(t, "gauge infinite", func() {
		render(func(c *ui.Context) { Gauge(c, "x", 0, math.Inf(-1), 1, GaugeOptions{}) })
	})
}

func TestLoaders(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Row(c).Gap(8).Children(func() {
			IndeterminateBar(c, "Connecting").Width(80)
			Spin(c, "Small", SpinOptions{Size: SpinSmall})
			Spin(c, "Large", SpinOptions{Size: SpinLarge})
			DotsLoader(c, "Dots")
			PulseLoader(c, "Pulse")
			BarsLoader(c, "Bars")
			OrbitLoader(c, "Orbit")
		})
	}, 600, 100)
	for _, s := range []string{"Connecting", "Small", "Large", "Dots", "Pulse", "Bars", "Orbit"} {
		find(t, tt, s)
	}
	a := tt.Image()
	before := append(a.Pix[:0:0], a.Pix...)
	tt.Advance(170 * time.Millisecond)
	if string(before) == string(tt.Image().Pix) {
		t.Error("loaders drew the same frame twice")
	}
	loops(t, tt, "loaders")
	if s, ok := tt.Find("Large"); !ok || s.W != ui.LightTheme().Space(10) {
		t.Errorf("large spin %v", s)
	}
	panics(t, "unknown size", func() { render(func(c *ui.Context) { Spin(c, "x", SpinOptions{Size: 7}) }) })
	panics(t, "no label", func() { render(func(c *ui.Context) { DotsLoader(c, "") }) })
}

func TestSkeletons(t *testing.T) {
	tt := ui.NewTester(func(c *ui.Context) {
		ui.Column(c).Gap(8).Width(320).Children(func() {
			Skeleton(c)
			SkeletonText(c, 3)
			SkeletonCard(c)
			SkeletonList(c, 2)
		})
	}, 400, 700)
	moved := false
	for range 20 {
		a := tt.Image()
		before := append(a.Pix[:0:0], a.Pix...)
		tt.Advance(100 * time.Millisecond)
		moved = moved || string(before) != string(tt.Image().Pix)
	}
	if !moved {
		t.Error("sheen did not move over two seconds")
	}
	loops(t, tt, "skeletons")
	panics(t, "no lines", func() { render(func(c *ui.Context) { SkeletonText(c, 0) }) })
	panics(t, "no rows", func() { render(func(c *ui.Context) { SkeletonList(c, 0) }) })
}
