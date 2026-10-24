package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MailRow shows a message and swipes to archive; cue Swipe.
func MailRow(c *ui.Context, m Mail, archived *bool) ui.Element {
	if archived == nil {
		panic("mail: MailRow needs archive state")
	}
	validMail(m)
	t := c.Theme()
	root := ui.Grid(c).Columns(1).FillWidth().Clip().Radius(t.Radius)
	w := root.LayoutBounds().W
	kit.Settle(c, root, []float32{w}, func() []float32 { return []float32{root.LayoutBounds().W} })
	drag := ui.Local(root, "drag", func() ui.Point { return ui.Point{} })
	root.Children(func() {
		ui.Row(c).ColumnStart(1).RowStart(1).FillWidth().FillHeight().AlignItems(ui.Center).
			Padding(0, t.Space(4)).Background(t.Accent).Children(func() {
			ui.Icon(c, icon.Must("archive")).Size(16, 16).TextColor(t.AccentText)
		})
		front := ui.Column(c).ColumnStart(1).RowStart(1).FillWidth().Gap(t.Space(1)).Padding(t.Space(2)).
			Background(t.Background).Border(1, t.Border).Focusable().Label(m.Subject).Draggable(drag, ui.DraggableOptions{
			Axis: ui.Horizontal, Animation: ui.Snappy,
			Snap: func(p ui.Point) ui.Point {
				if w > 0 && p.X > w/2 {
					return ui.Point{X: w}
				}
				return ui.Point{}
			},
		})
		front.Children(func() { mailText(c, m) })
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
