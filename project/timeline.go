package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TimelineTask is one Gantt interval as day fractions.
type TimelineTask struct {
	Title      string
	Start, End float64
}

// Timeline shows Gantt bars by task; drag resizes; cue Tick.
func Timeline(c *ui.Context, tasks *[]TimelineTask) ui.Element {
	if tasks == nil || len(*tasks) == 0 {
		panic("project: Timeline needs tasks")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, task := range *tasks {
			if task.Title == "" || task.Start < 0 || task.End <= task.Start || task.End > 1 || !finite(task.Start) || !finite(task.End) {
				panic("project: Timeline has an invalid task")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, task.Title).Width(t.Space(24))
				track := ui.Box(c).Grow(1).Height(t.Space(6)).Radius(t.Radius).Clip().Background(t.Surface).Cursor(ui.CursorGrab)
				moved := ui.Local(track, "moved", func() bool { return false })
				kit.Settle(c, track, []float32{track.LayoutBounds().W}, func() []float32 { return []float32{track.LayoutBounds().W} })
				track.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
					x := r.X + r.W*float32(task.Start)
					w := r.W * float32(task.End-task.Start)
					p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, t.Accent.Alpha(0.3), t.Radius)
				}))
				if dx, _, held := track.Dragged(); held {
					width := track.LayoutBounds().W
					if dx != 0 && width > 0 {
						next := min(1, max((*tasks)[i].Start+0.05, (*tasks)[i].End+float64(dx)/float64(width)))
						if next != (*tasks)[i].End {
							(*tasks)[i].End = next
							root.SetChanged()
							*moved = true
						}
					}
				} else if *moved {
					*moved = false
					sound.Play(sound.Tick)
					kit.Log("Timeline", "resized", "index", i)
				}
			})
		}
	})
	return root
}
