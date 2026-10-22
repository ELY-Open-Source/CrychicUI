package choice

import (
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

type rankDrag struct {
	id   uint64
	item string
}

var rankSlide = ui.ElementTransition{Animation: &ui.Snappy, Position: true}

// Ranking is rows dragged or Alt-arrowed into order.
func Ranking(c *ui.Context, label string, items *[]string) ui.Element {
	if label == "" || items == nil {
		panic("choice: Ranking needs a label and items")
	}
	distinct("Ranking", *items)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Label(label)
	id := root.ID()
	move := func(item string, to int) {
		from := slices.Index(*items, item)
		to = min(len(*items)-1, max(0, to))
		if from == to {
			return
		}
		next := slices.Delete(slices.Clone(*items), from, from+1)
		*items = slices.Insert(next, to, item)
		root.SetChanged()
		sound.Play(sound.Drop)
		kit.Log("Ranking", "moved", "from", from, "to", to)
	}
	return root.Children(func() {
		for i, item := range *items {
			c.Key(item)
			row := ui.Row(c).Gap(t.Space(2)).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Focusable().Label(item).
				Cursor(ui.CursorGrab).Drag(rankDrag{id, item}).Transition(rankSlide)
			edge := t.Border
			if d, over := ui.DragOver[rankDrag](row); over && d.id == id && d.item != item {
				edge = t.Accent
			}
			row.Background(t.Background).Border(1, row.AnimateColor("edge", edge, quick)).Shadow(0, 1, 3, 0, kit.Shade(0.06))
			if row.Dragging() {
				row.Opacity(0.4)
			}
			row.Children(func() {
				ui.Icon(c, icon.Must("grip-vertical")).Size(16, 16).TextColor(t.TextMuted)
				kit.Numeric(c, strconv.Itoa(i+1)).TextColor(t.TextMuted).MinWidth(t.Space(4))
				ui.Text(c, item).Grow(1)
			})
			if d, ok := ui.Drop[rankDrag](row); ok && d.id == id {
				move(d.item, i)
			}
			switch {
			case row.Shortcut(ui.Alt, ui.KeyUp):
				move(item, i-1)
			case row.Shortcut(ui.Alt, ui.KeyDown):
				move(item, i+1)
			}
		}
	})
}
