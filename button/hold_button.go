package button

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HoldButtonOptions tune a HoldButton.
type HoldButtonOptions struct {
	// Duration is the required hold time.
	Duration time.Duration
}

// HoldButton fills while held and fires when full.
func HoldButton(c *ui.Context, label string, o HoldButtonOptions) ui.Element {
	if o.Duration < 0 {
		panic("button: HoldButton duration must be positive")
	}
	if o.Duration == 0 {
		o.Duration = time.Second
	}
	t := c.Theme()
	b := textBase(c, label).Border(1, t.Border).Background(t.Surface).Clip()
	fired := ui.Local(b, "fired", func() bool { return false })
	wasPressed := ui.Local(b, "wasPressed", func() bool { return false })
	started := ui.Local(b, "started", func() time.Time { return time.Time{} })
	pressed := b.Pressed()
	if pressed && !*wasPressed {
		*started = c.Now()
		*fired = false
	}
	*wasPressed = pressed
	p := b.AnimateTo("hold", unit(pressed), ui.Animation{Duration: o.Duration, Ease: ui.Linear})
	if pressed && !*fired {
		elapsed := c.Now().Sub(*started)
		if elapsed >= o.Duration {
			*fired = true
			sound.Play(sound.Success)
			b.SetChanged()
			kit.Log("HoldButton", "completed")
		} else {
			c.After(o.Duration - elapsed)
		}
	}
	b.Draw(kit.Inside(t.Radius, func(painter *ui.Painter, r ui.Rect) {
		if p > 0 {
			painter.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * p, H: r.H}, t.Accent.Alpha(0.25), 0)
		}
	}))
	return b.Children(func() { ui.Text(c, label).SingleLine() })
}
