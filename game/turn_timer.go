package game

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TurnTimerOptions word the timer; Warn defaults to 10 seconds.
type TurnTimerOptions struct {
	Turn, Over string
	Warn       time.Duration
}

// TurnTimer drains a ring; nearing deadline plays Warning once.
func TurnTimer(c *ui.Context, player string, deadline time.Time, total time.Duration, o TurnTimerOptions) ui.Element {
	if strings.TrimSpace(player) == "" || total < time.Second || o.Warn < 0 || o.Warn > total {
		panic(fmt.Sprintf("game: TurnTimer %q needs a player, a turn of a second or more and a warning within it", player))
	}
	def(&o.Turn, "%s's turn")
	def(&o.Over, "Time's up")
	if o.Warn == 0 {
		o.Warn = min(10*time.Second, total)
	}
	t := c.Theme()
	left := countdown(c, deadline)
	share := float32(min(1, float64(left)/total.Seconds()))
	warn := left > 0 && left <= int64(o.Warn/time.Second)
	ink := t.Success
	switch {
	case left == 0:
		ink = t.Danger
	case warn:
		ink = t.Warning
	}
	root := card(c).Row().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(4)).Role(ui.RoleStatus).Label(fmt.Sprintf(o.Turn, player))
	warned := ui.Local(root, "warned", func() time.Time { return time.Time{} })
	if warn && !warned.Equal(deadline) {
		*warned = deadline
		sound.Play(sound.Warning)
		kit.Log("TurnTimer", "warned", "left", left)
	}
	root.Children(func() {
		size := t.Space(24)
		dial := ui.Box(c).Size(size, size).Shrink(0).Center()
		arc := dial.AnimateTo("left", share, ui.Smooth)
		col := dial.AnimateColor("ink", ink, ui.Snappy)
		if warn {
			s := 1 + 0.06*float32(dial.Loop("pulse", time.Second, ui.EaseInOut))
			dial.Scale(s, s)
		}
		dial.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			w := t.Space(2)
			cx, cy, rr := r.X+r.W/2, r.Y+r.H/2, r.W/2-w/2
			p.StrokePath(ring(cx, cy, rr), w, t.SurfaceHover)
			if arc > 0 {
				p.StrokePath(ring(cx, cy, rr).Trim(0, arc), w, col)
			}
		}))
		dial.Children(func() {
			n := kit.Numeric(c, clock(left)).Font("monospace").FontSize(t.Rem(1.05)).Bold().Shrink(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
			n.TextColor(n.AnimateColor("ink", t.Text.Mix(ink, unit(left == 0 || warn)), ui.Snappy))
		})
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
			ui.Text(c, fmt.Sprintf(o.Turn, player)).FontSize(t.Rem(0.95)).Bold().FillWidth().MinWidth(0).SingleLine()
			c.Key(left == 0)
			if left == 0 {
				kit.Arrive(kit.Chip(c, o.Over, t.Danger), ui.Pose{Scale: 0.7, Fade: 1}, 0)
			} else {
				kit.Muted(c, clock(int64(total/time.Second))).Font("monospace").FontSize(t.Rem(0.72)).Shrink(0)
			}
		})
	})
	return root
}
