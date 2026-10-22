package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PiPBoxOptions tune a PiPBox.
type PiPBoxOptions struct {
	// Height is the area's height; zero gives the default.
	Height float32
}

// PiPBox is thrown to a corner, 0 to 3.
func PiPBox(c *ui.Context, label string, corner *int, o PiPBoxOptions, content func()) ui.Element {
	if label == "" || corner == nil || content == nil || !valid(*corner, 4) || !finite(o.Height) || o.Height < 0 {
		panic(fmt.Sprintf("overlay: PiPBox %q needs a corner 0..3, a finite height and content", label))
	}
	t := c.Theme()
	h := o.Height
	if h == 0 {
		h = t.Space(80)
	}
	pad, bw, bh := t.Space(3), t.Space(40), t.Space(24)
	root := ui.Box(c).FillWidth().Height(h).Radius(t.Radius).Background(t.Surface).Border(1, t.Border)
	drag := ui.Local(root, "drag", func() ui.Point { return ui.Point{} })
	synced := ui.Local(root, "synced", func() [3]float32 { return [3]float32{-1} })
	held := ui.Local(root, "held", func() bool { return false })
	area := root.LayoutBounds()
	kit.Settle(c, root, []float32{area.W, area.H}, func() []float32 {
		a := root.LayoutBounds()
		return []float32{a.W, a.H}
	})
	room := ui.Rect{X: pad, Y: pad, W: max(0, area.W-bw-2*pad), H: max(0, area.H-bh-2*pad)}
	spot := func(i int) ui.Point {
		p := ui.Point{X: room.X, Y: room.Y}
		if i%2 == 1 {
			p.X += room.W
		}
		if i >= 2 {
			p.Y += room.H
		}
		return p
	}
	nearest := func(p ui.Point) ui.Point {
		best := spot(0)
		for i := 1; i < 4; i++ {
			if s := spot(i); dist(s, p) < dist(best, p) {
				best = s
			}
		}
		return best
	}
	key := [3]float32{float32(*corner), room.W, room.H}
	if !*held && key != *synced {
		*drag, *synced = spot(*corner), key
	}
	land := func(i int) {
		*corner, *drag = i, spot(i)
		*synced = [3]float32{float32(i), room.W, room.H}
		root.SetChanged()
		sound.Play(sound.Drop)
		kit.Log("PiPBox", "corner", "index", i)
	}
	root.Children(func() {
		box := ui.ButtonBase(c).Label(label).Absolute().Left(0).Top(0).Size(bw, bh).Radius(t.Radius).Background(t.Background).
			Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).Clip()
		box.Draggable(drag, ui.DraggableOptions{Snap: nearest, Bounds: &room, Animation: ui.Bouncy}).Children(content)
		if _, _, down := box.Dragged(); down {
			*held = true
		}
		to := *corner
		switch {
		case box.Shortcut(0, ui.KeyLeft):
			to = to &^ 1
		case box.Shortcut(0, ui.KeyRight):
			to = to | 1
		case box.Shortcut(0, ui.KeyUp):
			to = to % 2
		case box.Shortcut(0, ui.KeyDown):
			to = to%2 + 2
		}
		if to != *corner {
			land(to)
		}
		if *held && kit.Resting(box) {
			*held = false
			for i := range 4 {
				if *drag == spot(i) && i != *corner {
					land(i)
				}
			}
		}
	})
	return root
}
