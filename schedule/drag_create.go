package schedule

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DragCreateOptions are the drag-create labels.
type DragCreateOptions struct{ Label string }

// DragCreate drags over hours to create an event; cue Pop.
func DragCreate(c *ui.Context, start, end *time.Time, created *bool, o DragCreateOptions) ui.Element {
	if start == nil || start.IsZero() || end == nil || created == nil {
		panic("schedule: DragCreate needs start and state")
	}
	def(&o.Label, "Drag to create")
	t := c.Theme()
	root := ui.Box(c).FillWidth().Height(t.Space(24)).Radius(t.Radius).Background(t.Accent.Alpha(0.1)).Cursor(ui.CursorGrab).
		Role(ui.RoleGroup).Label(o.Label)
	visualY := ui.Local(root, "visualY", func() float32 { return 0 })
	moved := ui.Local(root, "moved", func() bool { return false })
	root.AnimatePose(ui.Pose{Y: *visualY}, ui.Snappy)
	root.Children(func() {
		ui.Text(c, o.Label).Padding(t.Space(2)).TextColor(t.Accent)
	})
	if _, dy, held := root.Dragged(); held {
		*visualY += dy
		if dy != 0 {
			*moved = true
			*end = start.Add(time.Duration(float32(*visualY)/t.Space(12)) * time.Hour)
		}
	} else if *moved {
		*moved = false
		*visualY = 0
		if end.After(*start) {
			*created = true
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("DragCreate", "created")
		}
	}
	return root
}
