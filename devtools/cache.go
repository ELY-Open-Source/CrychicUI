package devtools

import (
	"fmt"
	"strconv"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// CacheStats describe a cache; History holds hit rates, 0..1.
type CacheStats struct {
	Hits, Misses, Keys, Evictions, Size, Capacity int64
	History                                       []float64
}

// CacheOptions word the card.
type CacheOptions struct{ Title, HitRate, Hits, Misses, Keys, Evictions, Of, Idle string }

// Cache shows a cache's hit rate and size.
func Cache(c *ui.Context, s CacheStats, o CacheOptions) ui.Element {
	def(&o.Title, "Cache")
	if s.Hits < 0 || s.Misses < 0 || s.Keys < 0 || s.Evictions < 0 || s.Size < 0 || s.Capacity <= 0 || s.Size > s.Capacity {
		panic(fmt.Sprintf("devtools: Cache %q needs non-negative counts and a size within a positive capacity", o.Title))
	}
	for _, v := range s.History {
		if !finite(v) || v < 0 || v > 1 {
			panic(fmt.Sprintf("devtools: Cache %q hit rate %v is outside 0..1", o.Title, v))
		}
	}
	def(&o.HitRate, "hit rate")
	def(&o.Hits, "Hits")
	def(&o.Misses, "Misses")
	def(&o.Keys, "Keys")
	def(&o.Evictions, "Evictions")
	def(&o.Of, "of")
	def(&o.Idle, "No requests")
	t := c.Theme()
	requests := float64(s.Hits) + float64(s.Misses)
	rate := 0.0
	text := o.Idle
	ink := t.TextMuted
	if requests > 0 {
		rate = float64(s.Hits) / requests
		pct := rate * 100
		// Never 100% with a miss, nor 0% with a hit.
		if s.Misses > 0 {
			pct = min(pct, 99.9)
		}
		if s.Hits > 0 {
			pct = max(pct, 0.1)
		}
		text = strconv.FormatFloat(pct, 'f', 1, 64) + "%"
		switch {
		case rate >= 0.9:
			ink = t.Success
		case rate >= 0.7:
			ink = t.Warning
		default:
			ink = t.Danger
		}
	}
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "zap", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3), t.Space(4), t.Space(4), t.Space(4)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center).Children(func() {
				gauge := ui.Box(c).Size(t.Space(26), t.Space(26)).Shrink(0).Center()
				sweep := gauge.Keyframes("sweep", []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
				share := gauge.AnimateTo("rate", float32(rate), ui.Smooth)
				col := gauge.AnimateColor("ink", ink, ui.Snappy)
				gauge.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
					w := t.Space(2)
					cx, cy, rr := r.X+r.W/2, r.Y+r.H/2, min(r.W, r.H)/2-w/2
					p.StrokePath(ring(cx, cy, rr), w, t.SurfaceHover)
					if share*sweep > 0 {
						p.StrokePath(ring(cx, cy, rr).Trim(0, share*sweep), w, col)
					}
				}))
				gauge.Children(func() {
					ui.Column(c).AlignItems(ui.Center).Children(func() {
						kit.Numeric(c, text).FontSize(t.Rem(1.05)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
						kit.Muted(c, o.HitRate).FontSize(t.Rem(0.62)).Shrink(0)
					})
				})
				grid := ui.Grid(c).Columns(2).Grow(1).MinWidth(0).Gap(t.Space(2))
				grid.Children(func() {
					for i, stat := range []struct {
						name string
						n    int64
					}{{o.Hits, s.Hits}, {o.Misses, s.Misses}, {o.Keys, s.Keys}, {o.Evictions, s.Evictions}} {
						cell := ui.Column(c).MinWidth(0).Gap(t.Space(0.5)).Padding(t.Space(2)).Radius(t.Radius).Background(t.SurfaceHover.Alpha(0.6))
						kit.Stagger(c, grid, cell, i).Children(func() {
							kit.Muted(c, stat.name).FontSize(t.Rem(0.65)).Bold().FillWidth().MinWidth(0).SingleLine()
							kit.Numeric(c, number(c, stat.n)).FontSize(t.Rem(0.9)).Bold().FillWidth().MinWidth(0).SingleLine().
								ContentTransition(ui.ContentTransition{Numeric: true})
						})
					}
				})
			})
			ui.Column(c).FillWidth().Gap(t.Space(1.5)).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
					glyph(c, "hard-drive", 14).TextColor(t.TextMuted)
					typography.Bytes(c, s.Size).FontSize(t.Rem(0.75)).Bold().Shrink(0)
					kit.Muted(c, o.Of).FontSize(t.Rem(0.75)).Shrink(0)
					typography.Bytes(c, s.Capacity).FontSize(t.Rem(0.75)).TextColor(t.TextMuted).Shrink(0)
				})
				fill := float64(s.Size) / float64(s.Capacity)
				tone := t.Accent
				if fill >= 0.9 {
					tone = t.Warning
				}
				meter(c, float32(fill), tone).Height(t.Space(1.5))
			})
			if len(s.History) > 1 {
				chart := ui.Box(c).FillWidth().Height(t.Space(12))
				draw := chart.Keyframes("draw", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
				// Scale from just under the lowest rate, so changes show.
				lo, _ := bounds("Cache", s.History)
				lo = max(0, lo-0.05)
				chart.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
					area := spark(r, s.History, lo, 1)
					area.LineTo(r.X+r.W, r.Y+r.H)
					area.LineTo(r.X, r.Y+r.H)
					p.Clip(ui.Rect{X: r.X, Y: r.Y, W: r.W * draw, H: r.H}, 0, func() { p.FillPath(area.Close(), ink.Alpha(0.12)) })
					p.StrokePath(spark(r, s.History, lo, 1).Trim(0, draw), t.Space(0.5), ink)
				}))
			}
		})
	})
	return root
}
