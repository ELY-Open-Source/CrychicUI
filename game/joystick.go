package game

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// JoystickOptions word the stick.
type JoystickOptions struct{ Label string }

// Joystick drags a thumb in a circle; release springs home.
func Joystick(c *ui.Context, x, y *float32, o JoystickOptions) ui.Element {
	def(&o.Label, "Joystick")
	// Written as in-range tests so NaN fails them.
	if x == nil || y == nil || !(*x >= -1 && *x <= 1 && *y >= -1 && *y <= 1) {
		panic(fmt.Sprintf("game: Joystick %q needs x and y in -1..1", o.Label))
	}
	t := c.Theme()
	size, thumb := t.Space(36), t.Space(14)
	reach := (size - thumb) / 2
	root := ui.Box(c).Size(size, size).Shrink(0).Role(ui.RoleGroup)
	off := ui.Local(root, "off", func() [2]float32 { return [2]float32{*x * reach, *y * reach} })
	root.Children(func() {
		k := ui.ButtonBase(c).Label(o.Label).Absolute().Size(thumb, thumb).Radius(999).Cursor(ui.CursorGrab)
		dx, dy, held := k.Dragged()
		if held && (dx != 0 || dy != 0) {
			nx, ny := off[0]+dx, off[1]+dy
			if d := float32(math.Hypot(float64(nx), float64(ny))); d > reach {
				nx, ny = nx*reach/d, ny*reach/d
			}
			*off = [2]float32{nx, ny}
			*x, *y = nx/reach, ny/reach
			root.SetChanged()
		}
		moved := off[0] != 0 || off[1] != 0
		if !held && moved {
			*off = [2]float32{}
			*x, *y = 0, 0
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("Joystick", "released")
		}
		spring := ui.Bouncy
		if held {
			spring = ui.Snappy
		}
		px := k.AnimateTo("x", off[0], spring)
		py := k.AnimateTo("y", off[1], spring)
		k.Left(reach + px).Top(reach + py)
		lift := k.AnimateTo("lift", unit(held), ui.Snappy)
		k.Background(k.AnimateColor("bg", t.Accent.Mix(t.AccentHover, max(lift, unit(k.Hovered()))), ui.Snappy)).
			Shadow(0, 1+3*lift, 3+6*lift, 0, t.Text.Alpha(0.12+0.08*lift))
		root.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
			cx, cy := r.X+r.W/2, r.Y+r.H/2
			p.FillPath(ring(cx, cy, r.W/2), t.SurfaceHover)
			p.StrokePath(ring(cx, cy, r.W/2-t.Space(0.5)), t.Space(0.5), t.Border)
			p.StrokePath(ring(cx, cy, reach*0.5), t.Space(0.25), t.Border.Alpha(0.6))
			if px != 0 || py != 0 {
				var arm ui.Path
				arm.MoveTo(cx, cy).LineTo(cx+px, cy+py)
				p.StrokePath(&arm, t.Space(1.5), t.Accent.Alpha(0.3))
			}
		}))
	})
	return root
}
