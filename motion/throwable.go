package motion

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ThrowableOptions tune a Throwable.
type ThrowableOptions struct {
	// Slots is how many places the tile can land; 3.
	Slots int
}

// Throwable is a tile thrown along slots, landing in one.
func Throwable(c *ui.Context, label string, slot *int, o ThrowableOptions) ui.Element {
	if o.Slots == 0 {
		o.Slots = 3
	}
	if label == "" || slot == nil || o.Slots < 2 || *slot < 0 || *slot >= o.Slots {
		panic("motion: Throwable needs a label, two or more Slots and a slot among them")
	}
	t := c.Theme()
	size := t.Space(16)
	step := size + t.Space(2)
	dir := float32(1)
	if c.Direction() == ui.RightToLeft {
		dir = -1
	}
	points := make([]ui.Point, o.Slots)
	for i := range points {
		points[i] = ui.Point{X: dir * step * float32(i)}
	}
	span := step * float32(o.Slots-1)
	root := ui.Box(c).Width(span + size).Height(size)
	drag := ui.Local(root, "drag", func() ui.Point { return points[*slot] })
	synced := ui.Local(root, "synced", func() int { return *slot })
	pitch := ui.Local(root, "pitch", func() float32 { return dir * step })
	if *slot != *synced || dir*step != *pitch {
		*synced, *pitch, *drag = *slot, dir*step, points[*slot]
	}
	return root.Children(func() {
		for i := range o.Slots {
			c.Key(i)
			ui.Box(c).Absolute().Top(0).Start(step*float32(i)).Size(size, size).Radius(t.Radius).Border(1, t.Border)
		}
		tile := ui.ButtonBase(c).Label(label).Absolute().Top(0).Start(0).Size(size, size).Radius(t.Radius).Center().
			Background(t.Accent).TextColor(t.AccentText).Draggable(drag, ui.DraggableOptions{Axis: ui.Horizontal,
			Bounds: &ui.Rect{X: min(0, dir*span), W: span}, Snap: ui.Nearest(points...), Animation: ui.Snappy})
		at := int(math.Round(float64(abs(drag.X) / step)))
		switch {
		case tile.Shortcut(0, ui.KeyRight):
			*drag = points[min(max(at+int(dir), 0), o.Slots-1)]
		case tile.Shortcut(0, ui.KeyLeft):
			*drag = points[min(max(at-int(dir), 0), o.Slots-1)]
		}
		kit.Press(tile).Children(func() { ui.Text(c, label).Bold() })
		if at != *slot && resting(tile) {
			*slot, *synced = at, at
			root.SetChanged()
			sound.Play(sound.Drop)
			kit.Log("Throwable", "landed", "slot", at)
		}
	})
}
