package devtools

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RateLimitOptions word the card; Used takes used and limit.
type RateLimitOptions struct{ Title, Used, Left, Limited, Resets, Reset string }

// RateLimit shows requests against a limit; zero reset hides it.
func RateLimit(c *ui.Context, used, limit int64, reset time.Time, o RateLimitOptions) ui.Element {
	def(&o.Title, "Rate limit")
	if used < 0 || limit <= 0 {
		panic(fmt.Sprintf("devtools: RateLimit %q needs used of at least 0 and a positive limit", o.Title))
	}
	def(&o.Used, "%s of %s requests")
	def(&o.Left, "%s left")
	def(&o.Limited, "Limited")
	def(&o.Resets, "Resets in %s")
	def(&o.Reset, "Window reset")
	t := c.Theme()
	share := float64(used) / float64(limit)
	pct := min(100, math.Round(share*100))
	if used < limit {
		pct = min(pct, 99)
	}
	ink := t.Accent
	switch {
	case used >= limit:
		ink = t.Danger
	case share >= 0.8:
		ink = t.Warning
	}
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleStatus)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "gauge", ink, t.Space(9))
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				kit.Muted(c, o.Title).FontSize(t.Rem(0.75)).Bold()
				n := kit.Numeric(c, fmt.Sprintf("%.0f%%", pct)).FontSize(t.Rem(1.4)).Bold().Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
				n.TextColor(n.AnimateColor("ink", ink, ui.Snappy))
			})
			c.Key(used >= limit)
			if used >= limit {
				kit.Arrive(kit.Chip(c, o.Limited, t.Danger), ui.Pose{Scale: 0.8, Fade: 1}, 0)
			} else {
				kit.Chip(c, fmt.Sprintf(o.Left, number(c, limit-used)), ink)
			}
		})
		meter(c, float32(min(1, share)), ink).Height(t.Space(2))
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			kit.Muted(c, fmt.Sprintf(o.Used, number(c, used), number(c, limit))).FontSize(t.Rem(0.75)).Grow(1).MinWidth(0).Wrap()
			if reset.IsZero() {
				return
			}
			left := countdown(c, reset)
			text := o.Reset
			if left > 0 {
				text = fmt.Sprintf(o.Resets, span(left))
			}
			ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Shrink(0).Children(func() {
				glyph(c, "timer", 14).TextColor(t.TextMuted)
				kit.Numeric(c, text).FontSize(t.Rem(0.75)).TextColor(t.TextMuted).Shrink(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
			})
		})
	})
	return root
}
