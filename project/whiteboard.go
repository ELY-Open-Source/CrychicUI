package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// StickyNote is one whiteboard note at a fractional position.
type StickyNote struct {
	ID, Text string
	X, Y     float64
}

// Whiteboard drags sticky notes on a canvas; cue Drop.
func Whiteboard(c *ui.Context, notes *[]StickyNote) ui.Element {
	if notes == nil || len(*notes) == 0 {
		panic("project: Whiteboard needs notes")
	}
	ids := make([]string, len(*notes))
	for i, note := range *notes {
		ids[i] = note.ID
	}
	unique("Whiteboard", ids)
	t := c.Theme()
	root := ui.Box(c).FillWidth().Height(t.Space(80)).Radius(t.Radius).Clip().Background(t.Surface).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	step := t.Space(5)
	root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		dot := t.Space(0.25)
		for y := r.Y + step; y < r.Y+r.H; y += step {
			for x := r.X + step; x < r.X+r.W; x += step {
				p.Fill(ui.Rect{X: x - dot, Y: y - dot, W: 2 * dot, H: 2 * dot}, t.Border, dot)
			}
		}
	}))
	root.Children(func() {
		for i, note := range *notes {
			if note.ID == "" || note.Text == "" || note.X < 0 || note.X > 1 || note.Y < 0 || note.Y > 1 || !finite(note.X) || !finite(note.Y) {
				panic(fmt.Sprintf("project: Whiteboard has an invalid note %q", note.ID))
			}
			c.Key(note.ID)
			ink := labelInk(t, note.ID)
			card := ui.ButtonBase(c).Absolute().Width(t.Space(36)).Padding(t.Space(2.5)).Radius(t.Space(1)).
				Background(t.Background.Mix(ink, 0.16)).Border(1, ink.Alpha(0.35)).Cursor(ui.CursorGrab).Focusable().Label(note.Text)
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
			dx, dy, held := card.Dragged()
			lift := card.AnimateTo("lift", unit(held), ui.Snappy)
			card.AnimatePose(ui.Pose{Scale: 1 + 0.03*unit(held)}, ui.Snappy)
			card.Shadow(0, 1+3*lift, 3+9*lift, 0, kit.Shade(0.06+0.08*lift))
			card.Children(func() { ui.Text(c, note.Text).FontSize(t.Rem(0.8)).Wrap().FillWidth() })
			moved := ui.Local(card, "moved", func() bool { return false })
			if held {
				if used[0] > 0 && used[1] > 0 && (dx != 0 || dy != 0) {
					x := max(0, min(xLimit, (*notes)[i].X+float64(dx)/float64(used[0])))
					y := max(0, min(yLimit, (*notes)[i].Y+float64(dy)/float64(used[1])))
					if x != (*notes)[i].X || y != (*notes)[i].Y {
						(*notes)[i].X, (*notes)[i].Y = x, y
						root.SetChanged()
						*moved = true
					}
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
