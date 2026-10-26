package home

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LockOptions set the hold; Lock and Unlock take names.
type LockOptions struct {
	Hold                                               time.Duration
	Locked, Unlocked, Lock, Unlock, HoldLock, HoldOpen string
}

// Lock locks and unlocks after a held press; cue Success.
func Lock(c *ui.Context, name string, locked *bool, o LockOptions) ui.Element {
	if o.Hold == 0 {
		o.Hold = 600 * time.Millisecond
	}
	if strings.TrimSpace(name) == "" || locked == nil || o.Hold < 100*time.Millisecond || o.Hold > 5*time.Second {
		panic(fmt.Sprintf("home: Lock %q needs a locked state and a hold of 0.1 to 5 s", name))
	}
	def(&o.Locked, "Locked")
	def(&o.Unlocked, "Unlocked")
	def(&o.Lock, "Lock %s")
	def(&o.Unlock, "Unlock %s")
	def(&o.HoldLock, "Hold to lock")
	def(&o.HoldOpen, "Hold to unlock")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(5), t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	since := ui.Local(root, "since", func() time.Time { return time.Time{} })
	presses := ui.Local(root, "presses", func() int { return 0 })
	held := ui.Local(root, "held", func() bool { return false })
	done := ui.Local(root, "done", func() bool { return false })
	root.Children(func() {
		ui.Text(c, name).FontSize(t.Rem(0.92)).Bold().Shrink(0)
		label := o.Unlock
		if !*locked {
			label = o.Lock
		}
		size := t.Space(22)
		c.Key("lock")
		b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(label, name)).Size(size, size).Shrink(0).Center().Radius(999))
		now := c.Now()
		// Dragged's ok holds for the whole press, inside or out.
		_, _, down := b.Dragged()
		pressed := b.Pressed()
		if down && !*held {
			*since, *done = now, false
			*presses++
		}
		*held = down
		if down && !pressed {
			*done = true
		}
		if pressed && !*done {
			if left := o.Hold - now.Sub(*since); left > 0 {
				c.After(left + time.Millisecond)
			} else {
				*locked, *done = !*locked, true
				root.SetChanged()
				sound.Play(sound.Success)
				kit.Log("Lock", "turned", "locked", *locked)
			}
		}
		ink := t.Success
		mark := "lock"
		if !*locked {
			ink, mark = t.Warning, "lock-open"
		}
		frames := []ui.Keyframe{{Value: 0}, {At: 200 * time.Millisecond, Value: 0, Ease: ui.EaseOut}}
		run := ui.KeyframeOptions{Run: uint64(2**presses + 1)}
		if pressed && !*done {
			frames = []ui.Keyframe{{Value: 0}, {At: o.Hold, Value: 1}}
			run = ui.KeyframeOptions{Run: uint64(2 * *presses)}
		}
		k := b.Keyframes("hold", frames, run)
		b.Background(b.AnimateColor("bg", ink.Alpha(0.1).Mix(ink, 0.08*unit(b.Hovered())), ui.Snappy))
		b.TextColor(b.AnimateColor("ink", ink, ui.Snappy))
		b.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			w := t.Space(1)
			cx, cy, rr := r.X+r.W/2, r.Y+r.H/2, r.W/2-w/2
			p.StrokePath(ring(cx, cy, rr), w, ink.Alpha(0.25))
			if k > 0.01 {
				p.StrokePath(arc(cx, cy, rr, 0, 360*k), w, ink)
			}
		}))
		b.Children(func() {
			c.Key(mark)
			kit.Arrive(glyph(c, mark, size*0.36), ui.Pose{Scale: 0.5, Rotate: -20, Fade: 1}, 0)
		})
		state, hint := o.Locked, o.HoldOpen
		if !*locked {
			state, hint = o.Unlocked, o.HoldLock
		}
		ui.Column(c).AlignItems(ui.Center).Gap(t.Space(0.5)).Children(func() {
			c.Key(state)
			kit.Arrive(ui.Text(c, state).FontSize(t.Rem(1)).Bold().TextColor(ink).Shrink(0), ui.Pose{Y: t.Space(1), Fade: 1}, 0)
			kit.Muted(c, hint).FontSize(t.Rem(0.72)).Shrink(0)
		})
	})
	return root
}
