package game

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// HeartsOptions word the row; Label takes lives and max.
type HeartsOptions struct{ Label string }

// Hearts shows lives; a lost heart breaks apart.
func Hearts(c *ui.Context, lives, most int, o HeartsOptions) ui.Element {
	if most < 1 || most > 20 || lives < 0 || lives > most {
		panic(fmt.Sprintf("game: Hearts needs 1 to 20 hearts and 0 <= lives <= max, got %d of %d", lives, most))
	}
	def(&o.Label, "%d of %d lives")
	t := c.Theme()
	size := t.Space(9)
	root := ui.Row(c).Wrap().Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleMeter).Label(fmt.Sprintf(o.Label, lives, most))
	born := ui.Local(root, "born", func() time.Time { return c.Now() })
	// Hearts built in the first instant appear still.
	settled := c.Now().Sub(*born) > 50*time.Millisecond
	if !settled {
		c.After(51*time.Millisecond - c.Now().Sub(*born))
	}
	root.Children(func() {
		for i := range most {
			alive := i < lives
			c.Key([2]int{i, int(unit(alive))})
			h := ui.Box(c).Size(size, size).Shrink(0)
			if alive {
				if settled {
					kit.Arrive(h, ui.Pose{Scale: 0.3, Fade: 1}, 0)
				}
				h.Draw(func(p *ui.Painter, r ui.Rect) {
					b := body(r)
					p.FillPath(heart(b), t.Danger)
					p.FillPath(ring(b.X+b.W*0.32, b.Y+b.H*0.32, b.W*0.08), t.Background.Alpha(0.5))
				})
				continue
			}
			split := float32(1)
			if settled {
				split = h.Keyframes("break", []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			}
			h.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				b := body(r)
				p.StrokePath(heart(b), t.Space(0.375), t.TextMuted.Alpha(0.5))
				if split >= 1 {
					return
				}
				// Two halves part and drop within the box margin.
				fade := t.Danger.Alpha(1 - split)
				drop := (r.Y + r.H - b.Y - b.H) * split
				gap := (b.X - r.X) * split
				p.Clip(ui.Rect{X: b.X - gap, Y: b.Y + drop, W: b.W / 2, H: b.H}, 0, func() {
					p.FillPath(heart(ui.Rect{X: b.X - gap, Y: b.Y + drop, W: b.W, H: b.H}), fade)
				})
				p.Clip(ui.Rect{X: b.X + b.W/2 + gap, Y: b.Y + drop, W: b.W / 2, H: b.H}, 0, func() {
					p.FillPath(heart(ui.Rect{X: b.X + gap, Y: b.Y + drop, W: b.W, H: b.H}), fade)
				})
			}))
		}
	})
	return root
}

// body is the heart's rect, leaving room for a break.
func body(r ui.Rect) ui.Rect {
	return ui.Rect{X: r.X + r.W*0.14, Y: r.Y + r.H*0.04, W: r.W * 0.72, H: r.H * 0.72}
}
