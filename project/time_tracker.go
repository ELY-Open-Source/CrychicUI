package project

import (
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TimeTrackerOptions are the tracker labels.
type TimeTrackerOptions struct{ Start, Stop string }

// TimeTracker starts and stops with total; cue Toggle.
func TimeTracker(c *ui.Context, running *bool, total *time.Duration, o TimeTrackerOptions) ui.Element {
	if running == nil || total == nil || *total < 0 {
		panic("project: TimeTracker needs state")
	}
	def(&o.Start, "Start")
	def(&o.Stop, "Stop")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	started := ui.Local(root, "started", func() time.Time { return time.Time{} })
	wasRunning := ui.Local(root, "wasRunning", func() bool { return *running })
	if *running != *wasRunning {
		if *running {
			*started = c.Now()
		}
		*wasRunning = *running
	}
	if *running && (*started).IsZero() {
		*started = c.Now()
	}
	root.Children(func() {
		kit.Numeric(c, total.String()).Bold().Grow(1)
		name := o.Start
		if *running {
			name = o.Stop
		}
		b := kit.Press(ui.ButtonBase(c).Label(name).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
			Background(t.Accent).TextColor(t.AccentText))
		b.Children(func() { ui.Text(c, name).SingleLine() })
		if b.Clicked() {
			if *running {
				elapsed := c.Now().Sub(*started)
				if elapsed > time.Duration(math.MaxInt64)-*total {
					panic("project: TimeTracker total overflow")
				}
				*total += elapsed
				*running = false
			} else {
				*started = c.Now()
				*running = true
			}
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("TimeTracker", "toggled", "running", *running)
		}
	})
	return root
}
