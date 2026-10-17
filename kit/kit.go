// Package kit holds helpers three or more components share.
package kit

import (
	"log/slog"
	"time"

	"git.zacharyzhang.com/ZacharyZhang-NY/CrychicUI/sound"
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
		Shadow(0, 1, 3, 0, ui.RGBA(0, 0, 0, 0.06))
}

// Muted returns secondary text.
func Muted(c *ui.Context, s string) ui.Element {
	return ui.Text(c, s).TextColor(c.Theme().TextMuted)
}

// Log records a component state change at debug level.
func Log(component, event string, args ...any) {
	slog.Debug("crychic: "+component+" "+event, args...)
}
