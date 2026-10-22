package motion

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CardStackOptions tune a CardStack.
type CardStackOptions struct {
	// Empty shows once every card is gone; "No more cards".
	Empty string
}

// CardStack stacks cards; drag the top away for the next.
func CardStack(c *ui.Context, n int, top *int, card func(i int), o CardStackOptions) ui.Element {
	if n < 0 || top == nil || *top < 0 || *top > n || card == nil {
		panic("motion: CardStack needs n ≥ 0, a top from 0 to n and a card")
	}
	if o.Empty == "" {
		o.Empty = "No more cards"
	}
	t := c.Theme()
	away := t.Space(120)
	root := ui.Grid(c).Columns(1)
	drag := ui.Local(root, "drag", func() ui.Point { return ui.Point{} })
	shown := ui.Local(root, "top", func() int { return *top })
	if *top != *shown {
		*shown, *drag = *top, ui.Point{}
	}
	fling := func(p ui.Point) ui.Point {
		switch {
		case p.X > away/4:
			return ui.Point{X: away}
		case p.X < -away/4:
			return ui.Point{X: -away}
		}
		return ui.Point{}
	}
	return root.Children(func() {
		if *top == n {
			kit.Muted(c, o.Empty).Padding(t.Space(4))
			return
		}
		for i := min(*top+3, n) - 1; i >= *top; i-- {
			c.Key(i)
			e := kit.Card(c).Padding(t.Space(4)).ColumnStart(1).RowStart(1)
			d := e.AnimateTo("depth", float32(i-*top), ui.Snappy)
			e.Scale(1-0.05*d, 1-0.05*d).Translate(0, d*t.Space(2))
			if i == *top {
				e.Draggable(drag, ui.DraggableOptions{Axis: ui.Horizontal, Snap: fling, Animation: ui.Snappy})
				x := e.DragOffset().X
				e.Rotate(12 * x / away).Opacity(1 - min(abs(x)/away, 1))
				if drag.X != 0 && resting(e) {
					*top++
					*shown, *drag = *top, ui.Point{}
					root.SetChanged()
					sound.Play(sound.Swipe)
					kit.Log("CardStack", "dismissed", "index", i)
				}
			}
			e.Children(func() { card(i) })
		}
	})
}
