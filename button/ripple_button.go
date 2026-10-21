package button

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RippleButton grows a ripple from the press point.
func RippleButton(c *ui.Context, label string) ui.Element {
	t := c.Theme()
	root := textBase(c, label).Border(1, t.Border).Background(t.Surface).Clip()
	x := ui.Local(root, "x", func() float32 { return 0 })
	y := ui.Local(root, "y", func() float32 { return 0 })
	armed := ui.Local(root, "armed", func() bool { return false })
	activated := ui.Local(root, "activated", func() bool { return false })
	run := ui.Local(root, "run", func() uint64 { return 0 })
	bounds := root.LayoutBounds()
	kit.Settle(c, root, []float32{bounds.W, bounds.H}, func() []float32 {
		b := root.LayoutBounds()
		return []float32{b.W, b.H}
	})
	root.HandleInput(func(ev ui.InputEvent) bool {
		switch ev.Kind {
		case ui.InputPointerDown:
			if ev.Button == 0 && !*armed {
				*armed = true
				*x, *y = ev.X, ev.Y
				return true
			}
		case ui.InputPointerMove:
			return *armed
		case ui.InputPointerUp:
			if *armed {
				*armed = false
				if ev.X >= 0 && ev.Y >= 0 && ev.X <= bounds.W && ev.Y <= bounds.H {
					*activated = true
				}
				return true
			}
		}
		return false
	})
	if root.Clicked() || *activated {
		if !*activated {
			*x, *y = bounds.W/2, bounds.H/2
		}
		*activated = false
		*armed = false
		(*run)++
		sound.Play(sound.Tap)
	}
	if *run > 0 {
		p := root.Keyframes("ripple", []ui.Keyframe{{Value: 0}, {At: 400 * time.Millisecond, Value: 1, Ease: ui.EaseOut}},
			ui.KeyframeOptions{Run: *run})
		root.Draw(func(painter *ui.Painter, r ui.Rect) {
			if p < 1 {
				radius := max(r.W, r.H) * p
				circle := (&ui.Path{}).Circle(r.X+*x, r.Y+*y, radius)
				painter.FillPath(circle, t.Accent.Alpha(0.25*(1-p)))
			}
		})
	}
	return root.Children(func() { ui.Text(c, label).SingleLine() })
}
