// Package chart draws data with MyGo's ui/chart.
package chart

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// tone is series i's color: theme hues, lighter each lap.
func tone(t *ui.Theme, i int) ui.Color {
	base := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	return base[i%len(base)].Mix(t.Background, 0.35*float32(i/len(base)%3))
}

// grow rises 0 to 1 once, as marks first show.
func grow(e ui.Element) float32 {
	return e.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
}

// spring moves marks to new values.
var spring = ui.ElementTransition{Animation: &ui.Snappy}

// springs eases each value to its target on e.
func springs(e ui.Element, key string, vs []float64) []float64 {
	out := make([]float64, len(vs))
	for i, v := range vs {
		out[i] = float64(e.AnimateTo(fmt.Sprintf("%s.%d", key, i), float32(v), ui.Snappy))
	}
	return out
}

// eased springs plotted fractions as marks' layout springs.
func eased(e ui.Element, key string, fs []float32) []float32 {
	out := make([]float32, len(fs))
	for i, f := range fs {
		out[i] = e.AnimateTo(fmt.Sprintf("%s.%d", key, i), f, ui.Snappy)
	}
	return out
}

// finite panics on NaN or infinite values.
func finite(component, label string, vs ...float64) {
	for _, v := range vs {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			panic(fmt.Sprintf("chart: %s %q values must be finite", component, label))
		}
	}
}

// chosen checks a selection is -1 or below n.
func chosen(component, label string, at *int, n int) {
	if at == nil || *at < -1 || *at >= n {
		panic(fmt.Sprintf("chart: %s %q needs a choice from -1 to %d", component, label, n-1))
	}
}

// pick selects i on a click.
func pick(root ui.Element, at *int, i int, component string) {
	*at = i
	root.SetChanged()
	sound.Play(sound.Select)
	kit.Log(component, "chose", "index", i)
}

// def fills an empty copy field.
func def(field *string, english string) {
	if *field == "" {
		*field = english
	}
}

// height is h, or the default for zero.
func height(c *ui.Context, h float32) float32 {
	if h < 0 || math.IsNaN(float64(h)) || math.IsInf(float64(h), 0) {
		panic("chart: Height must be finite and not negative")
	}
	if h == 0 {
		return c.Theme().Space(56)
	}
	return h
}

// side is a size v, or fallback for zero.
func side(comp, label string, v, fallback float32) float32 {
	if v < 0 || math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
		panic(fmt.Sprintf("chart: %s %q sizes must be finite and not negative", comp, label))
	}
	if v == 0 {
		return fallback
	}
	return v
}

// plot is a labeled plot with muted axis labels.
func plot(c *ui.Context, label string, axes mc.Axes, h float32, build func(area ui.Element)) ui.Element {
	t := c.Theme()
	axes.Label = func(_ mc.Side, tick mc.Tick) {
		ui.Text(c, tick.Label).NoWrap().FontSize(t.Rem(0.7)).TextColor(t.TextMuted)
	}
	return mc.Plot(c, axes, build).Label(label).FillWidth().Height(height(c, h)).Gap(t.Space(1.5))
}

// grid draws a hairline at each tick.
func grid(t *ui.Theme, ticks []mc.Tick) func(p *ui.Painter, r ui.Rect) {
	return func(p *ui.Painter, r ui.Rect) {
		for _, tick := range ticks {
			y := mc.Y(r, tick.At)
			p.Line(r.X, y, r.X+r.W, y, 1, t.Border)
		}
	}
}

// follow tips the hovered point; a click picks it.
func follow(c *ui.Context, root, area ui.Element, at, hover *int, xs []float32, top float32, describe func(i int) string, component string) {
	t := c.Theme()
	mc.Select(c, area, hover, xs, describe)
	if area.Clicked() && *hover >= 0 {
		pick(root, at, *hover, component)
	}
	if i := *hover; i >= 0 {
		mc.RuleX(c, xs[i], 1).Background(t.TextMuted.Alpha(0.5))
		mc.Anchor(c, xs[i], top, ui.Center, ui.End, func() {
			tip := ui.Text(c, describe(i)).NoWrap().FontSize(t.Rem(0.75)).Padding(t.Space(0.5), t.Space(2)).Radius(t.Radius).
				Background(t.Background).Border(1, t.Border).Margin(0, 0, t.Space(1), 0)
			kit.Arrive(tip, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
		})
	}
}

// legend lists series names by their colors.
func legend(c *ui.Context, names []string) {
	t := c.Theme()
	ui.Row(c).Wrap().Gap(t.Space(3)).Children(func() {
		for i, n := range names {
			c.Key(i)
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Radius(999).Background(tone(t, i))
				kit.Muted(c, n).NoWrap().FontSize(t.Rem(0.8))
			})
		}
	})
}

// named checks names are non-empty and distinct.
func named(component, label string, names []string) {
	if label == "" || len(names) == 0 {
		panic(fmt.Sprintf("chart: %s needs a label and names", component))
	}
	seen := map[string]bool{}
	for _, n := range names {
		if n == "" || seen[n] {
			panic(fmt.Sprintf("chart: %s %q names need text and no twins", component, label))
		}
		seen[n] = true
	}
}

// minTime is the earlier of a and b.
func minTime(a, b time.Time) time.Time {
	if b.Before(a) {
		return b
	}
	return a
}

// maxTime is the later of a and b.
func maxTime(a, b time.Time) time.Time {
	if b.After(a) {
		return b
	}
	return a
}
