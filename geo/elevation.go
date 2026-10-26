package geo

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ElevationOptions word the profile; High, Low and At format values.
type ElevationOptions struct {
	Title, High, Low, At string
	Imperial             bool
	Units                Units
}

// Elevation is a route's height profile; hover reads a point.
func Elevation(c *ui.Context, samples []float64, length float64, o ElevationOptions) ui.Element {
	def(&o.Title, "Elevation")
	def(&o.High, "High %s")
	def(&o.Low, "Low %s")
	def(&o.At, "%s at %s")
	o.Units.fill()
	if len(samples) < 2 || len(samples) > 200 {
		panic(fmt.Sprintf("geo: Elevation %q needs 2 to 200 samples, got %d", o.Title, len(samples)))
	}
	meters("Elevation", o.Title, length)
	lo, hi, gain, loss := samples[0], samples[0], 0.0, 0.0
	for i, h := range samples {
		if !(h >= -500 && h <= 9000) {
			panic(fmt.Sprintf("geo: Elevation %q needs heights from -500 to 9,000 m, got %v", o.Title, h))
		}
		lo, hi = min(lo, h), max(hi, h)
		if i > 0 {
			if d := h - samples[i-1]; d > 0 {
				gain += d
			} else {
				loss -= d
			}
		}
	}
	height := func(m float64) string {
		if o.Imperial {
			return fmt.Sprintf(o.Units.Feet, strconv.FormatFloat(math.Round(m*3.28084), 'f', 0, 64))
		}
		return fmt.Sprintf(o.Units.Meters, strconv.FormatFloat(math.Round(m), 'f', 0, 64))
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			tile(c, "mountain", t.Success, t.Space(9))
			ui.Text(c, o.Title).FontSize(t.Rem(0.9)).Bold().Grow(1).MinWidth(0).Wrap()
		})
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(1.5)).Children(func() {
			for i, v := range [2]float64{gain, loss} {
				mark, ink, sign := "trending-up", t.Success, "+"
				if i == 1 {
					mark, ink, sign = "trending-down", t.Danger, "−"
				}
				ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999).Background(ink.Alpha(0.12)).Children(func() {
					glyph(c, mark, 12).TextColor(ink)
					ui.Text(c, sign+height(v)).FontSize(t.Rem(0.72)).Bold().TextColor(ink).Shrink(0).MaxWidth(t.Space(56)).SingleLine()
				})
			}
		})
		span0 := max(hi-lo, 1)
		c.Key("chart")
		chart := ui.Row(c).FillWidth().Height(t.Space(36)).AlignItems(ui.Stretch)
		grow := chart.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		hover := -1
		chart.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
			n := len(samples)
			at := func(i int) (float32, float32) {
				x := r.X + r.W*float32(i)/float32(n-1)
				y := r.Y + r.H - r.H*0.85*float32((samples[i]-lo)/span0)*grow - r.H*0.08
				return x, y
			}
			var area, line ui.Path
			area.MoveTo(r.X, r.Y+r.H)
			for i := range n {
				x, y := at(i)
				area.LineTo(x, y)
				if i == 0 {
					line.MoveTo(x, y)
				} else {
					line.LineTo(x, y)
				}
			}
			area.LineTo(r.X+r.W, r.Y+r.H).Close()
			p.FillPath(&area, t.Success.Alpha(0.16))
			p.StrokePath(&line, t.Space(0.5), t.Success)
			if hover >= 0 {
				x, y := at(hover)
				p.Fill(ui.Rect{X: x - 0.5, Y: r.Y, W: 1, H: r.H}, t.Text.Alpha(0.35), 0)
				var dot ui.Path
				p.FillPath(dot.Circle(x, y, t.Space(1.25)), t.Success)
			}
		}))
		chart.Children(func() {
			for i := range samples {
				c.Key(i)
				if ui.Box(c).Grow(1).Basis(0).MinWidth(0).FillHeight().Hovered() {
					hover = i
				}
			}
		})
		c.Key("readout")
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			if i := hover; i >= 0 {
				where := span(length*float64(i)/float64(len(samples)-1), o.Imperial, o.Units)
				kit.Numeric(c, fmt.Sprintf(o.At, height(samples[i]), where)).FontSize(t.Rem(0.78)).Bold().Grow(1).MinWidth(0).SingleLine()
				return
			}
			kit.Muted(c, fmt.Sprintf(o.High, height(hi))).FontSize(t.Rem(0.76)).Shrink(0).MaxWidth(t.Space(56)).SingleLine()
			kit.Muted(c, fmt.Sprintf(o.Low, height(lo))).FontSize(t.Rem(0.76)).Shrink(0).MaxWidth(t.Space(56)).SingleLine()
			ui.Box(c).Grow(1)
			kit.Muted(c, span(length, o.Imperial, o.Units)).FontSize(t.Rem(0.76)).Shrink(0).MaxWidth(t.Space(56)).SingleLine()
		})
	})
}
