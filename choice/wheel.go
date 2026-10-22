package choice

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Wheel is a snapping picker wheel that ticks per row.
func Wheel(c *ui.Context, label string, items []string, selected *int) ui.Element {
	distinct("Wheel", items)
	if label == "" || selected == nil || *selected < 0 || *selected >= len(items) {
		panic("choice: Wheel needs a label and a valid selection")
	}
	t := c.Theme()
	row, n := t.Space(8), len(items)
	at := func(y float32) int { return min(n-1, max(0, int(math.Round(float64(-y/row))))) }
	root := ui.Box(c).Height(row*5).MinWidth(t.Space(30)).Radius(t.Radius).Background(t.Background).Border(1, t.Border).
		Clip().Focusable().Label(label).Role(ui.RoleSlider)
	off := ui.Local(root, "off", func() ui.Point { return ui.Point{Y: -float32(*selected) * row} })
	last := ui.Local(root, "last", func() int { return *selected })
	tick := ui.Local(root, "tick", func() int { return *selected })
	user := ui.Local(root, "user", func() bool { return false })
	spin := ui.Local(root, "spin", func() float32 { return 0 })
	root.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind != ui.InputScroll {
			return false
		}
		*spin += ev.DY
		return true
	})
	step := int(*spin / row)
	*spin -= float32(step) * row
	switch {
	case root.Shortcut(0, ui.KeyDown):
		step = 1
	case root.Shortcut(0, ui.KeyUp):
		step = -1
	}
	if v := min(n-1, max(0, *selected+step)); v != *selected {
		*selected, *user = v, true
		root.SetChanged()
		kit.Log("Wheel", "keyed", "index", v)
	}
	if was := ui.Local(root, "row", func() float32 { return row }); *selected != *last || *was != row {
		*off, *last, *was = ui.Point{Y: -float32(*selected) * row}, *selected, row
	}
	root.Range(0, float64(n-1), float64(*selected)).Step(1)
	return root.Children(func() {
		ui.Box(c).Absolute().Top(row * 2).Left(t.Space(1)).Right(t.Space(1)).Height(row).Radius(t.Radius).Background(t.Surface)
		col := ui.Column(c).Absolute().Top(row * 2).Left(0).Right(0).Cursor(ui.CursorGrab)
		col.Draggable(off, ui.DraggableOptions{
			Axis:      ui.Vertical,
			Bounds:    &ui.Rect{Y: -float32(n-1) * row, H: float32(n-1) * row},
			Snap:      func(p ui.Point) ui.Point { return ui.Point{Y: -float32(at(p.Y)) * row} },
			Animation: ui.Snappy,
		})
		shown := col.DragOffset().Y
		*user = *user || col.Dragging()
		if i := at(shown); i != *tick {
			*tick = i
			if *user {
				sound.Play(sound.Tick)
			}
		}
		if !col.Dragging() && shown == off.Y {
			*user = false
		}
		if i := at(off.Y); !col.Dragging() && i != *selected {
			*selected, *last = i, i
			root.SetChanged()
			kit.Log("Wheel", "dragged", "index", i)
		}
		col.Children(func() {
			for i, item := range items {
				d := float32(math.Abs(float64(float32(i)*row+shown))) / row
				c.Key(i)
				ui.Text(c, item).SingleLine().Height(row).TextAlign(ui.Center).Padding(t.Space(2), t.Space(3)).
					Opacity(max(0.25, 1-0.3*d)).Scale(max(0.8, 1-0.06*d), max(0.8, 1-0.06*d))
			}
		})
	})
}
