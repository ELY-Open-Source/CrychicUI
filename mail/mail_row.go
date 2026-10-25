package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MailRow shows a message and swipes to archive; cue Swipe.
func MailRow(c *ui.Context, m Mail, archived *bool) ui.Element {
	if archived == nil {
		panic("mail: MailRow needs archive state")
	}
	validMail("MailRow", m)
	t := c.Theme()
	root := ui.Grid(c).Columns(1).FillWidth().Clip().Radius(t.Radius)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	w := root.LayoutBounds().W
	kit.Settle(c, root, []float32{w}, func() []float32 { return []float32{root.LayoutBounds().W} })
	drag := ui.Local(root, "drag", func() ui.Point { return ui.Point{} })
	reach := float32(0)
	if w > 0 {
		reach = min(1, max(0, drag.X/(w/2)))
	}
	root.Children(func() {
		ui.Row(c).ColumnStart(1).RowStart(1).FillWidth().FillHeight().AlignItems(ui.Center).Padding(0, t.Space(4)).
			Background(t.Success.Alpha(0.25 + 0.75*reach)).TextColor(t.AccentText).Children(func() {
			glyph(c, "archive", 18).Scale(0.8+0.3*reach, 0.8+0.3*reach)
		})
		front := ui.Column(c).ColumnStart(1).RowStart(1).FillWidth().Padding(t.Space(2.5)).Radius(t.Radius).
			Background(t.Background).Border(1, t.Border).Shadow(0, 1, 2+6*reach, 0, kit.Shade(0.05+0.06*reach)).
			Focusable().Label(m.Subject).Cursor(ui.CursorGrab).Draggable(drag, ui.DraggableOptions{
			Axis: ui.Horizontal, Animation: ui.Snappy,
			Snap: func(p ui.Point) ui.Point {
				if w > 0 && p.X > w/2 {
					return ui.Point{X: w}
				}
				return ui.Point{}
			},
		})
		front.Children(func() { mailFace(c, m) })
		if w > 0 && kit.Resting(front) && drag.X >= w && !*archived {
			*archived = true
			*drag = ui.Point{}
			root.SetChanged()
			sound.Play(sound.Swipe)
			kit.Log("MailRow", "archived")
		}
	})
	return root
}
