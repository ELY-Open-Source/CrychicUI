// Package kit holds helpers three or more components share.
package kit

import (
	"log/slog"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Press shrinks e while it is pressed.
func Press(e ui.Element) ui.Element {
	p := ui.Pose{}
	if e.Pressed() {
		p.Scale = 0.97
	}
	return e.AnimatePose(p, ui.Snappy)
}

// Lift raises e 2 DIP while it is hovered.
func Lift(e ui.Element) ui.Element {
	p := ui.Pose{}
	if e.Hovered() {
		p.Y = -2
	}
	return e.AnimatePose(p, ui.Snappy)
}

// Arrive eases from a pose to rest on first build.
func Arrive(e ui.Element, from ui.Pose, delay time.Duration) ui.Element {
	p := e.Keyframes("crychic.arrive", []ui.Keyframe{{Value: 1}, {At: 200 * time.Millisecond, Value: 0, Ease: ui.EaseOut}},
		ui.KeyframeOptions{Delay: delay})
	s := from.Scale
	if s == 0 {
		s = 1
	}
	k := 1 + (s-1)*p
	return e.Translate(from.X*p, from.Y*p).Scale(k, k).Opacity(1 - from.Fade*p)
}

var popSpring = ui.Bouncy

var (
	// Enter fades an element in from 4 DIP below.
	Enter = ui.ElementTransition{
		Duration: 160 * time.Millisecond,
		Enter:    &ui.Motion{Y: 4},
		Exit:     &ui.Motion{Y: 4},
	}
	// Pop scales an element in and out with a bounce.
	Pop = ui.ElementTransition{
		Animation: &popSpring,
		Enter:     &ui.Motion{Scale: 0.9},
		Exit:      &ui.Motion{Scale: 0.9},
	}
)

// Clicked reports e.Clicked and plays cue when true.
func Clicked(e ui.Element, cue sound.Cue) bool {
	if !e.Clicked() {
		return false
	}
	sound.Play(cue)
	return true
}

// Changed reports e.Changed and plays cue when true.
func Changed(e ui.Element, cue sound.Cue) bool {
	if !e.Changed() {
		return false
	}
	sound.Play(cue)
	return true
}

// Card returns the shared card surface.
func Card(c *ui.Context) ui.Element {
	t := c.Theme()
	return ui.Column(c).Background(t.Background).Border(1, t.Border).Radius(t.Radius).
		Shadow(0, 1, 3, 0, Shade(0.06))
}

// Settle rebuilds when e's laid-out read differs from used.
func Settle(c *ui.Context, e ui.Element, used []float32, read func() []float32) {
	e.DrawOver(func(*ui.Painter, ui.Rect) {
		if !slices.Equal(read(), used) {
			c.Invalidate()
		}
	})
}

// Resting reports e let go and still.
func Resting(e ui.Element) bool {
	_, _, held := e.Dragged()
	return !held && !e.Dragging() && !e.Transitioning()
}

// Slide keys content, sliding in toward the side moved to.
func Slide(c *ui.Context, key int, content func()) {
	t := c.Theme()
	box := ui.Box(c).Clip()
	last := ui.Local(box, "key", func() int { return key })
	dir := ui.Local(box, "dir", func() float32 { return 1 })
	if key != *last {
		*dir = 1
		if key < *last {
			*dir = -1
		}
		*last = key
	}
	shift := *dir * t.Space(8)
	if c.Direction() == ui.RightToLeft {
		shift = -shift
	}
	box.Children(func() {
		ui.Column(c).Key(key).Transition(ui.ElementTransition{Enter: &ui.Motion{X: shift}, Exit: &ui.Motion{X: -shift}}).
			Children(content)
	})
}

// Inside clips fn to its box rounded by radius.
func Inside(radius float32, fn func(p *ui.Painter, r ui.Rect)) func(p *ui.Painter, r ui.Rect) {
	return func(p *ui.Painter, r ui.Rect) { p.Clip(r, min(radius, r.W/2, r.H/2), func() { fn(p, r) }) }
}

// Shade is black at alpha a, for shadows.
func Shade(a float32) ui.Color { return ui.RGBA(0, 0, 0, a) }

// Muted returns secondary text.
func Muted(c *ui.Context, s string) ui.Element {
	return ui.Text(c, s).TextColor(c.Theme().TextMuted)
}

// Numeric returns text whose digits roll on change.
func Numeric(c *ui.Context, s string) ui.Element {
	return ui.Text(c, s).ContentTransition(ui.ContentTransition{Numeric: true})
}

// Log records a component state change at debug level.
func Log(component, event string, args ...any) {
	slog.Debug("crychic: "+component+" "+event, args...)
}
