package choice

import (
	"math"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// rail draws slider track, fill and knob at frac.
func rail(c *ui.Context, s ui.Element, frac float32, vertical bool, extra func(knob ui.Element)) ui.Element {
	t := c.Theme()
	kw, kh, h, side := t.Space(5), t.Space(4), t.Space(1.5), t.Space(5)
	if vertical {
		kw, kh = kh, kw
		s.Width(side).PaddingY(kh / 2)
	} else {
		s.Height(side).MinWidth(t.Space(20)).PaddingX(kw / 2)
	}
	_, _, held := s.Dragged()
	rtl := c.Direction() == ui.RightToLeft
	return s.Children(func() {
		ui.Box(c).Fill().Children(func() {
			tr := ui.Box(c).Absolute().Radius(h / 2).Background(t.Border)
			fill := ui.Box(c).Absolute().Radius(h / 2).Background(t.Accent)
			knob := ui.Box(c).Absolute().Size(kw, kh).Radius(min(kw, kh)/2).Background(t.AccentText).Border(1, t.Border).
				Shadow(0, 1.5, 7, 0, kit.Shade(0.1))
			if vertical {
				x := (side - h) / 2
				tr.Top(0).Bottom(0).Left(x).Width(h)
				fill.Bottom(0).Left(x).Width(h).HeightPercent(frac * 100)
				knob.BottomPercent(frac*100).Left((side-kw)/2).Translate(0, kh/2)
			} else {
				y := (side - h) / 2
				tr.Start(0).End(0).Top(y).Height(h)
				fill.Start(0).Top(y).Height(h).WidthPercent(frac * 100)
				dx := -kw / 2
				if rtl {
					dx = kw / 2
				}
				knob.StartPercent(frac*100).Top((side-kh)/2).Translate(dx, 0)
			}
			pose := ui.Pose{}
			if held {
				pose.Scale = 1.2
			}
			knob.AnimatePose(pose, ui.Snappy)
			if extra != nil {
				knob.Children(func() { extra(knob) })
			}
		})
	})
}

// reach is hi as MyGo's step grid reaches it.
func reach(lo, hi, step float64) float64 {
	if step <= 0 {
		return hi
	}
	// MyGo drops a top above hi before rounding.
	if top := lo + math.Round((hi-lo)/step)*step; top > hi && snap(top, lo, step) == hi {
		return top
	}
	return hi
}

// fraction is where v sits from lo to hi.
func fraction(v, lo, hi float64) float32 {
	return float32(max(0, min(1, (v-lo)/(hi-lo))))
}

// released plays Tick after a drag or key move.
func released(s ui.Element, component string, value float64) {
	moved := ui.Local(s, "moved", func() bool { return false })
	*moved = *moved || s.Changed()
	if _, _, held := s.Dragged(); *moved && !held {
		*moved = false
		sound.Play(sound.Tick)
		kit.Log(component, "changed", "value", value)
	}
}

// sliderArgs panics on a bad slider range or step.
func sliderArgs(component, label string, value *float64, lo, hi, step float64) {
	if label == "" || value == nil || !(lo < hi) || math.IsInf(lo, 0) || math.IsInf(hi, 0) ||
		!(*value >= lo && *value <= hi) || !(step >= 0) || math.IsInf(step, 0) {
		panic("choice: " + who(component, label) + " needs a label, finite lo < hi, a value within and a finite step")
	}
}

// onStep reports whether MyGo's slider snap keeps v.
func onStep(v, lo, step float64) bool { return snap(v, lo, step) == v }

// snap rounds v to whole steps, as MyGo does.
func snap(v, lo, step float64) float64 {
	s := lo + math.Round((v-lo)/step)*step
	if d := max(decimals(step), decimals(lo)); d < 15 {
		r, err := strconv.ParseFloat(strconv.FormatFloat(s, 'f', d, 64), 64)
		if err != nil {
			panic(err)
		}
		s = r
	}
	return s
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
