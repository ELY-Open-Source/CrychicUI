package devtools

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MemoryOptions word the card; Value, Peak and Limit are formats.
type MemoryOptions struct{ Title, Value, Peak, Limit string }

// Memory charts heap use over time against a limit.
func Memory(c *ui.Context, samples []float64, limit float64, o MemoryOptions) ui.Element {
	lo, hi := bounds("Memory", samples)
	if lo < 0 || !finite(limit) || limit <= 0 || len(samples) < 2 {
		panic(fmt.Sprintf("devtools: Memory %q needs two or more non-negative samples and a positive limit", o.Title))
	}
	def(&o.Title, "Heap")
	def(&o.Value, "%s MB")
	def(&o.Peak, "Peak %s")
	def(&o.Limit, "Limit %s")
	t := c.Theme()
	now := samples[len(samples)-1]
	ink := t.Accent
	switch {
	case now >= limit*0.9:
		ink = t.Danger
	case now >= limit*0.7:
		ink = t.Warning
	}
	// Scale by the peak first so headroom cannot overflow.
	peak := max(hi, limit)
	scaled := make([]float64, len(samples))
	for i, v := range samples {
		scaled[i] = v / peak
	}
	num := func(v float64) string { return fmt.Sprintf(o.Value, strconv.FormatFloat(math.Round(v), 'g', 6, 64)) }
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "memory-stick", ink, t.Space(9))
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				kit.Muted(c, o.Title).FontSize(t.Rem(0.75)).Bold()
				kit.Numeric(c, num(now)).FontSize(t.Rem(1.3)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
			})
			ui.Column(c).Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
				kit.Chip(c, fmt.Sprintf(o.Peak, num(hi)), t.Warning)
				kit.Chip(c, fmt.Sprintf(o.Limit, num(limit)), t.Danger)
			})
		})
		chart := ui.Box(c).FillWidth().Height(t.Space(28))
		draw := chart.Keyframes("draw", []ui.Keyframe{{Value: 0}, {At: 800 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		col := chart.AnimateColor("ink", ink, ui.Snappy)
		chart.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			ly := r.Y + r.H - r.H*float32(limit/peak/1.1)
			for x := r.X; x < r.X+r.W; x += t.Space(3) {
				p.Line(x, ly, min(x+t.Space(1.5), r.X+r.W), ly, 1, t.Danger.Alpha(0.6))
			}
			area := spark(r, scaled, 0, 1.1)
			area.LineTo(r.X+r.W, r.Y+r.H)
			area.LineTo(r.X, r.Y+r.H)
			p.Clip(ui.Rect{X: r.X, Y: r.Y, W: r.W * draw, H: r.H}, 0, func() { p.FillPath(area.Close(), col.Alpha(0.15)) })
			p.StrokePath(spark(r, scaled, 0, 1.1).Trim(0, draw), t.Space(0.5), col)
		}))
	})
	return root
}
