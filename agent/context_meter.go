package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ContextMeterOptions tune a ContextMeter.
type ContextMeterOptions struct {
	Label, Tokens, Warning string
}

// ContextMeter is a ring of context used; warns near full.
func ContextMeter(c *ui.Context, used, limit int, o ContextMeterOptions) ui.Element {
	if used < 0 || limit <= 0 || used > limit {
		panic(fmt.Sprintf("agent: ContextMeter needs 0 ≤ used ≤ limit, got %d of %d", used, limit))
	}
	def(&o.Label, "Context window")
	def(&o.Tokens, "%s of %s tokens")
	def(&o.Warning, "Older messages will be summarized")
	t := c.Theme()
	share := float64(used) / float64(limit)
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
		Radius(t.Radius).Background(t.Background).Border(1, t.Border).Role(ui.RoleGroup).Label(o.Label)
	return root.Children(func() {
		side, w := t.Space(12), t.Space(1.25)
		ring := ui.Box(c).Size(side, side).Shrink(0).Center().Role(ui.RoleProgress).Label(o.Label).Range(0, 1, share)
		arc := ring.AnimateTo("arc", float32(share), ui.Smooth)
		ink := ring.AnimateColor("ink", meterInk(c, share), ui.Snappy)
		ring.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			cx, cy, rad := r.X+r.W/2, r.Y+r.H/2, min(r.W, r.H)/2-w
			var track, fill ui.Path
			p.StrokePath(track.Circle(cx, cy, rad), w, t.Border)
			if arc > 0 {
				p.Transform(ui.Rotation(-90), cx, cy, func() { p.StrokePath(fill.Circle(cx, cy, rad).Trim(0, arc), w, ink) })
			}
		}))
		ring.Children(func() { kit.Numeric(c, fmt.Sprintf("%.0f%%", share*100)).FontSize(t.Rem(0.75)).Bold() })
		ui.Column(c).Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, o.Label).Bold().SingleLine()
			kit.Numeric(c, fmt.Sprintf(o.Tokens, short(used), short(limit))).FontSize(t.Rem(0.8)).TextColor(t.TextMuted)
			if share >= 0.8 {
				warn := kit.Muted(c, o.Warning).FontSize(t.Rem(0.75)).TextColor(meterInk(c, share))
				kit.Arrive(warn, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
			}
		})
	})
}

// short writes n as 950, 12.4k or 1.2M.
func short(n int) string {
	switch {
	case n >= 1_000_000:
		return trim(fmt.Sprintf("%.1f", float64(n)/1e6)) + "M"
	case n >= 1000:
		return trim(fmt.Sprintf("%.1f", float64(n)/1e3)) + "k"
	}
	return fmt.Sprint(n)
}

func trim(s string) string {
	if len(s) > 2 && s[len(s)-2:] == ".0" {
		return s[:len(s)-2]
	}
	return s
}
