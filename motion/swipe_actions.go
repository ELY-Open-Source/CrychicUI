package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SwipeActions reveals actions when its row is dragged aside.
func SwipeActions(c *ui.Context, actions []string, chosen *int, content func()) ui.Element {
	if len(actions) == 0 || chosen == nil || content == nil {
		panic("motion: SwipeActions needs actions, a chosen pointer and content")
	}
	for _, a := range actions {
		if a == "" {
			panic("motion: SwipeActions actions need labels")
		}
	}
	t := c.Theme()
	each := t.Space(20)
	reveal := each * float32(len(actions))
	open := -reveal
	if c.Direction() == ui.RightToLeft {
		open = reveal
	}
	root := ui.Grid(c).Columns(1).ClipX().Radius(t.Radius).Border(1, t.Border)
	drag := ui.Local(root, "drag", func() ui.Point { return ui.Point{} })
	shut := ui.Local(root, "shut", func() bool { return true })
	was := ui.Local(root, "open", func() float32 { return open })
	if *was != open {
		*was = open
		if !*shut {
			*drag = ui.Point{X: open}
		}
	}
	return root.Children(func() {
		ui.Row(c).ColumnStart(1).RowStart(1).Justify(ui.End).Children(func() {
			for i, a := range actions {
				c.Key(i)
				b := ui.ButtonBase(c).Label(a).Width(each).Center()
				bg := t.Surface
				if b.Hovered() {
					bg = t.SurfaceHover
				}
				kit.Press(b.Background(bg)).Children(func() { ui.Text(c, a) })
				if kit.Clicked(b, sound.Tap) {
					*chosen, *drag, *shut = i, ui.Point{}, true
					root.SetChanged()
					kit.Log("SwipeActions", "chose", "index", i)
				}
			}
		})
		front := ui.Box(c).ColumnStart(1).RowStart(1).Background(t.Background).Draggable(drag, ui.DraggableOptions{Axis: ui.Horizontal,
			Bounds: &ui.Rect{X: min(open, 0), W: reveal}, Snap: ui.Nearest(ui.Point{}, ui.Point{X: open}), Animation: ui.Snappy})
		front.Children(content)
		if at := drag.X == 0; at != *shut && kit.Resting(front) {
			*shut = at
			sound.Play(sound.Swipe)
			kit.Log("SwipeActions", "swiped", "open", !at)
		}
	})
}
