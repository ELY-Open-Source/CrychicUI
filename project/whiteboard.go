package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StickyNote is one whiteboard note.
type StickyNote struct {
	ID, Text string
	X, Y     float64
}

// Whiteboard drags sticky notes on a canvas; cue Drop.
func Whiteboard(c *ui.Context, notes *[]StickyNote) ui.Element {
	if notes == nil || len(*notes) == 0 {
		panic("project: Whiteboard needs notes")
	}
	t := c.Theme()
	root := ui.Box(c).FillWidth().Height(t.Space(80)).Radius(t.Radius).Clip().Background(t.Surface)
	root.Children(func() {
		for i, note := range *notes {
			if note.ID == "" || note.Text == "" || note.X < 0 || note.X > 1 || note.Y < 0 || note.Y > 1 || !finite(note.X) || !finite(note.Y) {
				panic("project: Whiteboard has an invalid note")
			}
			c.Key(note.ID)
			card := ui.ButtonBase(c).Absolute().Width(t.Space(36)).Padding(t.Space(2)).Radius(t.Space(0.5)).
				Background(t.Accent.Alpha(0.2)).Cursor(ui.CursorGrab).Focusable().Label(note.Text)
			room := func() []float32 {
				a, s := root.LayoutBounds(), card.LayoutBounds()
				return []float32{a.W, a.H, max(0, a.W-s.W), max(0, a.H-s.H)}
			}
			used := room()
			kit.Settle(c, card, used, room)
			xLimit, yLimit := 1.0, 1.0
			if used[0] > 0 {
				xLimit = float64(used[2]) / float64(used[0])
			}
			if used[1] > 0 {
				yLimit = float64(used[3]) / float64(used[1])
			}
			(*notes)[i].X = max(0, min(xLimit, (*notes)[i].X))
			(*notes)[i].Y = max(0, min(yLimit, (*notes)[i].Y))
			card.Left(float32((*notes)[i].X) * used[0]).Top(float32((*notes)[i].Y) * used[1])
			card.Children(func() { ui.Text(c, note.Text).Wrap() })
			moved := ui.Local(card, "moved", func() bool { return false })
			if dx, dy, held := card.Dragged(); held {
				if used[0] > 0 && used[1] > 0 && (dx != 0 || dy != 0) {
					(*notes)[i].X = max(0, min(xLimit, (*notes)[i].X+float64(dx)/float64(used[0])))
					(*notes)[i].Y = max(0, min(yLimit, (*notes)[i].Y+float64(dy)/float64(used[1])))
					root.SetChanged()
					*moved = true
				}
			} else if *moved {
				*moved = false
				sound.Play(sound.Drop)
				kit.Log("Whiteboard", "moved", "id", note.ID)
			}
		}
	})
	kit.Settle(c, root, []float32{root.LayoutBounds().W, root.LayoutBounds().H}, func() []float32 {
		return []float32{root.LayoutBounds().W, root.LayoutBounds().H}
	})
	return root
}
