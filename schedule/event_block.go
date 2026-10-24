package schedule

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// EventBlock drags a duration block; cue Drop.
func EventBlock(c *ui.Context, event *Event) ui.Element {
	if event == nil {
		panic("schedule: EventBlock needs an event")
	}
	validEvent("EventBlock", *event)
	t := c.Theme()
	ink := inks(t)[event.Tint]
	height := float32(event.End.Sub(event.Start)) / float32(time.Hour) * t.Space(12)
	root := ui.Row(c).FillWidth().Height(height).Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2), t.Space(1.5), t.Space(3)).
		Radius(t.Radius).Clip().Cursor(ui.CursorGrab).Role(ui.RoleGroup).Label(event.Title)
	dragY := ui.Local(root, "dragY", func() float32 { return 0 })
	visualY := ui.Local(root, "visualY", func() float32 { return 0 })
	moved := ui.Local(root, "moved", func() bool { return false })
	held := ui.Local(root, "held", func() bool { return false })
	root.AnimatePose(ui.Pose{Y: *visualY, Scale: 1 + 0.02*unit(*held)}, ui.Snappy)
	root.Background(root.AnimateColor("bg", ink.Alpha(0.14).Mix(ink.Alpha(0.24), unit(*held || root.Hovered())), ui.Snappy))
	if *held {
		root.Shadow(0, 6, 16, 0, kit.Shade(0.14))
	}
	root.Draw(func(p *ui.Painter, r ui.Rect) { p.Fill(ui.Rect{X: r.X, Y: r.Y, W: t.Space(1), H: r.H}, ink, 0) })
	root.Children(func() {
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).TextColor(ink.Mix(t.Text, 0.35)).Children(func() {
			ui.Text(c, event.Title).FontSize(t.Rem(0.8)).Bold().SingleLine().MinWidth(0)
			kit.Numeric(c, eventText(*event, event.Start.Location())).FontSize(t.Rem(0.7))
			if event.Place != "" {
				ui.Text(c, event.Place).FontSize(t.Rem(0.7)).SingleLine().MinWidth(0).Opacity(0.8)
			}
		})
		glyph(c, "grip-vertical", 14).TextColor(ink).AlignSelf(ui.Center)
	})
	_, dy, drag := root.Dragged()
	*held = drag
	if drag {
		*dragY += dy
		*visualY += dy
		steps := int(*dragY / t.Space(12))
		if steps != 0 {
			delta := time.Duration(steps) * time.Hour
			event.Start, event.End = event.Start.Add(delta), event.End.Add(delta)
			*dragY -= float32(steps) * t.Space(12)
			*moved = true
			root.SetChanged()
		}
	} else {
		*dragY = 0
		*visualY = 0
		if *moved {
			*moved = false
			sound.Play(sound.Drop)
			kit.Log("EventBlock", "moved", "start", event.Start.Format("15:04"))
		}
	}
	return root
}
