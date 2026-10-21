package button

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UndoButtonOptions tune an UndoButton.
type UndoButtonOptions struct {
	// Duration is the time before committing.
	Duration time.Duration
}

// UndoButton drains a ring before committing or undoing.
func UndoButton(c *ui.Context, label string, o UndoButtonOptions) ui.Element {
	if o.Duration < 0 {
		panic("button: UndoButton duration must be positive")
	}
	if o.Duration == 0 {
		o.Duration = 3 * time.Second
	}
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Background(t.Surface)
	done := ui.Local(root, "done", func() bool { return false })
	started := ui.Local(root, "started", func() time.Time { return c.Now() })
	root.Disabled(*done)
	progress := float32(0)
	if !*done {
		progress = root.Keyframes("drain", []ui.Keyframe{{Value: 1}, {At: o.Duration, Value: 0, Ease: ui.Linear}})
		if remaining := o.Duration - c.Now().Sub(*started); remaining > 0 {
			c.After(remaining)
		}
	}
	if root.Clicked() {
		*done = true
		sound.Play(sound.Close)
		kit.Log("UndoButton", "undone")
	}
	if !*done && c.Now().Sub(*started) >= o.Duration {
		*done = true
		root.SetChanged()
		kit.Log("UndoButton", "committed")
	}
	return root.Children(func() {
		ring := ui.Box(c).Size(16, 16)
		ring.Draw(func(p *ui.Painter, r ui.Rect) {
			if progress > 0 {
				circle := (&ui.Path{}).Circle(r.X+r.W/2, r.Y+r.H/2, 6)
				p.StrokePath(circle.Trim(0, progress), 1.5, t.Accent)
			}
		})
		ui.Text(c, label).SingleLine()
	})
}
