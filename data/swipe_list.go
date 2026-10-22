package data

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SwipeList rows swipe right to archive, left to delete.
func SwipeList(c *ui.Context, label string, items, archived *[]string) ui.Element {
	if label == "" || items == nil || archived == nil {
		panic(fmt.Sprintf("data: SwipeList %q needs items and an archive", label))
	}
	unique("SwipeList", label, *items)
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleList).Label(label)
	drop := func(it string, keep bool) {
		*items = slices.DeleteFunc(slices.Clone(*items), func(s string) bool { return s == it })
		root.SetChanged()
		if keep {
			*archived = append(*archived, it)
			sound.Play(sound.Swipe)
		} else {
			sound.Play(sound.Delete)
		}
		kit.Log("SwipeList", "removed", "archived", keep)
	}
	gone := ui.ElementTransition{Animation: &ui.Snappy, Exit: &ui.Motion{Collapse: true}}
	// Rows fill the list; exiting rows must not be measured.
	w := root.LayoutBounds().W
	kit.Settle(c, root, []float32{w}, func() []float32 { return []float32{root.LayoutBounds().W} })
	root.Children(func() {
		for _, it := range slices.Clone(*items) {
			c.Key(it)
			row := ui.Grid(c).Columns(1).FillWidth().Clip().Radius(t.Radius).Transition(gone)
			drag := ui.Local(row, "drag", func() ui.Point { return ui.Point{} })
			row.Children(func() {
				ui.Row(c).ColumnStart(1).RowStart(1).FillWidth().FillHeight().Children(func() {
					ui.Row(c).Grow(1).FillHeight().AlignItems(ui.Center).Padding(0, t.Space(4)).Background(t.Success).Children(func() {
						ui.Icon(c, icon.Must("archive")).Size(16, 16).TextColor(t.AccentText)
					})
					ui.Row(c).Grow(1).FillHeight().AlignItems(ui.Center).Justify(ui.End).Padding(0, t.Space(4)).Background(t.Danger).Children(func() {
						ui.Icon(c, icon.Must("trash-2")).Size(16, 16).TextColor(t.AccentText)
					})
				})
				front := ui.Row(c).ColumnStart(1).RowStart(1).FillWidth().Focusable().Label(it).Padding(t.Space(3)).Background(t.Background).
					Border(1, t.Border).Radius(t.Radius).Draggable(drag, ui.DraggableOptions{Axis: ui.Horizontal, Animation: ui.Snappy,
					Snap: func(p ui.Point) ui.Point {
						switch {
						case w > 0 && p.X > w/2:
							return ui.Point{X: w}
						case w > 0 && p.X < -w/2:
							return ui.Point{X: -w}
						}
						return ui.Point{}
					}})
				front.Children(func() { ui.Text(c, it).SingleLine() })
				switch {
				case front.Shortcut(ui.Shift, ui.KeyDelete):
					drop(it, true)
				case front.Shortcut(0, ui.KeyDelete), front.Shortcut(0, ui.KeyBackspace):
					drop(it, false)
				case w > 0 && kit.Resting(front) && drag.X >= w:
					drop(it, true)
				case w > 0 && kit.Resting(front) && drag.X <= -w:
					drop(it, false)
				}
			})
		}
	})
	return root
}
