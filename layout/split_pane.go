package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

const minPane = 40

// SplitPane sets two panes side by side; *size is width.
func SplitPane(c *ui.Context, size *float32, first, second func()) ui.Element {
	return split(c, "SplitPane", false, size, first, second)
}

// split lays out two panes around a draggable handle.
func split(c *ui.Context, name string, vertical bool, size *float32, first, second func()) ui.Element {
	if size == nil || !finite(*size) || *size < 0 {
		panic("layout: " + name + " needs a finite, non-negative size")
	}
	t := c.Theme()
	grip := t.Space(2)
	var root ui.Element
	var room float32
	pane, cursor := func(e ui.Element, v float32) ui.Element { return e.Width(v).ClipX() }, ui.CursorResizeEW
	if vertical {
		root = ui.Column(c)
		room = root.LayoutBounds().H
		pane, cursor = func(e ui.Element, v float32) ui.Element { return e.Height(v).ClipY() }, ui.CursorResizeNS
	} else {
		root = ui.Row(c)
		room = root.LayoutBounds().W
	}
	return root.Fill().AlignItems(ui.Stretch).Children(func() {
		pane(ui.Column(c), *size).Children(first)
		handle := ui.Box(c).Cursor(cursor).Center().Label("Divider")
		if vertical {
			handle.Height(grip)
		} else {
			handle.Width(grip)
		}
		dx, dy, moved := handle.Dragged()
		active := handle.Hovered() || handle.Dragging()
		color := t.Border
		if active {
			color = t.Accent
		}
		line := ui.Box(c)
		if vertical {
			line.FillWidth().Height(1)
		} else {
			line.FillHeight().Width(1)
		}
		handle.Children(func() { line.Background(line.AnimateColor("line", color, ui.Snappy)) })
		if moved {
			d := dx
			if vertical {
				d = dy
			}
			next := max(minPane, *size+d)
			if room > 0 {
				next = min(next, max(minPane, room-grip-minPane))
			}
			if next != *size {
				*size = next
				root.SetChanged()
			}
		}
		if _, _, ended := handle.DragEnded(); ended {
			sound.Play(sound.Tick)
			kit.Log(name, "resized", "size", *size)
		}
		ui.Column(c).Grow(1).Basis(0).Children(second)
	})
}
