package game

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ComboOptions word the counter; Times takes the count.
type ComboOptions struct{ Label, Times, Best, Lost string }

// Combo shows a bouncing count; losing it shows a note.
func Combo(c *ui.Context, count int, o ComboOptions) ui.Element {
	if count < 0 {
		panic(fmt.Sprintf("game: Combo %q count %d is negative", o.Label, count))
	}
	def(&o.Label, "Combo")
	def(&o.Times, "×%d")
	def(&o.Best, "Best %d")
	def(&o.Lost, "Combo lost")
	t := c.Theme()
	ink := t.Accent
	switch {
	case count >= 10:
		ink = t.Danger
	case count >= 5:
		ink = t.Warning
	}
	root := card(c).Row().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(4)).Role(ui.RoleStatus).Label(o.Label)
	last := ui.Local(root, "last", func() int { return count })
	best := ui.Local(root, "best", func() int { return count })
	bump := ui.Local(root, "bump", func() time.Time { return time.Time{} })
	lost := ui.Local(root, "lost", func() time.Time { return time.Time{} })
	now := c.Now()
	if count > *last {
		*bump = now
	}
	if count == 0 && *last > 0 {
		*lost = now
	}
	*last = count
	*best = max(*best, count)
	fresh := now.Sub(*bump) < 120*time.Millisecond
	gone := !lost.IsZero() && count == 0 && now.Sub(*lost) < 1500*time.Millisecond
	if fresh {
		c.After(120*time.Millisecond - now.Sub(*bump))
	}
	if gone {
		c.After(1500*time.Millisecond - now.Sub(*lost))
	}
	root.Children(func() {
		mark := "zap"
		if count >= 10 {
			mark = "flame"
		}
		c.Key(mark)
		kit.Arrive(tile(c, mark, ink, t.Space(12)), ui.Pose{Scale: 0.6, Rotate: -20, Fade: 1}, 0)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			kit.Muted(c, o.Label).FontSize(t.Rem(0.75)).Bold()
			n := kit.Numeric(c, fmt.Sprintf(o.Times, count)).FontSize(t.Rem(2)).Bold().Shrink(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
			s := n.AnimateTo("bump", 1+0.3*unit(fresh), ui.Bouncy)
			n.Scale(s, s).TextColor(n.AnimateColor("ink", ink, ui.Snappy))
		})
		ui.Column(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.End).Children(func() {
			kit.Chip(c, fmt.Sprintf(o.Best, *best), t.Warning)
			c.Key("lost")
			note := kit.Chip(c, o.Lost, t.Danger)
			// A brief push off center, then the spring shakes back.
			push := float32(0)
			if gone && now.Sub(*lost) < 80*time.Millisecond {
				push = t.Space(2)
				c.After(81*time.Millisecond - now.Sub(*lost))
			}
			shake := note.AnimateTo("shake", push, ui.Bouncy)
			note.Translate(shake, 0).Opacity(note.AnimateTo("shown", unit(gone), ui.Smooth))
		})
	})
	return root
}
