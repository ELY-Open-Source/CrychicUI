package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LatencyOptions word the card.
type LatencyOptions struct{ Title, P50, P95, P99 string }

// Latency shows three percentiles on tiles and a scale.
func Latency(c *ui.Context, p50, p95, p99 float64, o LatencyOptions) ui.Element {
	if p50 < 0 || p95 < p50 || p99 < p95 || !finite(p50) || !finite(p95) || !finite(p99) {
		panic("dashboard: Latency needs finite p50 ≤ p95 ≤ p99 of zero or more")
	}
	def(&o.Title, "Latency")
	def(&o.P50, "p50")
	def(&o.P95, "p95")
	def(&o.P99, "p99")
	t := c.Theme()
	scale := p99
	if scale == 0 {
		scale = 1
	}
	rows := []struct {
		name  string
		value float64
		ink   ui.Color
	}{{o.P50, p50, t.Success}, {o.P95, p95, t.Warning}, {o.P99, p99, t.Danger}}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "timer", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			tiles := ui.Row(c).FillWidth().Gap(t.Space(2))
			tiles.Children(func() {
				for i, row := range rows {
					c.Key(i)
					box := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(0.5)).Padding(t.Space(2.5)).Radius(t.Radius).Background(t.Surface)
					kit.Stagger(c, tiles, box, i).Children(func() {
						ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							ui.Box(c).Size(t.Space(2), t.Space(2)).Shrink(0).Radius(999).Background(row.ink)
							kit.Muted(c, row.name).FontSize(t.Rem(0.75)).Bold().Shrink(0)
						})
						kit.Numeric(c, duration(row.value)).Bold().FontSize(t.Rem(1.1)).FillWidth().SingleLine().
							ContentTransition(ui.ContentTransition{Numeric: true})
					})
				}
			})
			track := ui.Box(c).FillWidth().Height(t.Space(5))
			xs := make([]float32, len(rows))
			for i, row := range rows {
				xs[i] = track.AnimateTo(fmt.Sprint("x", i), float32(row.value/scale), ui.Smooth)
			}
			track.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				y := r.Y + r.H/2
				p.Fill(ui.Rect{X: r.X, Y: y - t.Space(0.5), W: r.W, H: t.Space(1)}, t.SurfaceHover, 999)
				p.Fill(ui.Rect{X: r.X, Y: y - t.Space(0.5), W: r.W * xs[0], H: t.Space(1)}, t.Success.Alpha(0.5), 999)
				for i, row := range rows {
					x, d := r.X+t.Space(1.5)+(r.W-t.Space(3))*xs[i], t.Space(3)
					p.Fill(ui.Rect{X: x - d/2, Y: y - d/2, W: d, H: d}, t.Background, 999)
					p.Fill(ui.Rect{X: x - d/3, Y: y - d/3, W: 2 * d / 3, H: 2 * d / 3}, row.ink, 999)
				}
			}))
		})
	})
	return root
}

// duration writes milliseconds, switching to seconds from one second.
func duration(ms float64) string {
	if ms >= 1000 {
		return fmt.Sprintf("%.2f s", ms/1000)
	}
	return fmt.Sprintf("%.0f ms", ms)
}
