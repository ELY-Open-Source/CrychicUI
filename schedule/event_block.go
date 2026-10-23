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
	height := float32(event.End.Sub(event.Start)) / float32(time.Hour) * t.Space(12)
	root := ui.Box(c).FillWidth().Height(height).Radius(t.Radius).Clip().Background(t.Accent.Alpha(0.2)).Cursor(ui.CursorGrab)
	dragY := ui.Local(root, "dragY", func() float32 { return 0 })
	visualY := ui.Local(root, "visualY", func() float32 { return 0 })
	moved := ui.Local(root, "moved", func() bool { return false })
	root.AnimatePose(ui.Pose{Y: *visualY}, ui.Snappy)
	root.Children(func() {
		ui.Text(c, event.Title).FontSize(t.Rem(0.75)).TextColor(t.Accent).Padding(t.Space(1))
	})
	if _, dy, held := root.Dragged(); held {
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
			kit.Log("EventBlock", "moved")
		}
	}
	return root
}
