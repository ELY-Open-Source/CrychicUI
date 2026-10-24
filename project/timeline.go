package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TimelineTask is one Gantt interval as span fractions.
type TimelineTask struct {
	Title      string
	Start, End float64
}

// TimelineOptions head the span with evenly spaced Ticks.
type TimelineOptions struct{ Ticks []string }

// Timeline shows Gantt bars by task; drag resizes; cue Tick.
func Timeline(c *ui.Context, tasks *[]TimelineTask, o TimelineOptions) ui.Element {
	if tasks == nil || len(*tasks) == 0 {
		panic("project: Timeline needs tasks")
	}
	names := make([]string, len(*tasks))
	for i, task := range *tasks {
		if task.Start < 0 || task.End <= task.Start || task.End > 1 || !finite(task.Start) || !finite(task.End) {
			panic(fmt.Sprintf("project: Timeline task %q has an invalid span", task.Title))
		}
		names[i] = task.Title
	}
	unique("Timeline", names)
	t := c.Theme()
	palette := []ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Role(ui.RoleGroup)
	root.Children(func() {
		if len(o.Ticks) > 0 {
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				ui.Box(c).Width(t.Space(24)).Shrink(0)
				ui.Row(c).Grow(1).MinWidth(0).Children(func() {
					for i, tick := range o.Ticks {
						c.Key(i)
						ui.Box(c).Grow(1).Basis(0).MinWidth(0).Children(func() {
							kit.Muted(c, tick).FontSize(t.Rem(0.65)).Bold().SingleLine().MinWidth(0)
						})
					}
				})
			})
		}
		for i, task := range *tasks {
			ink := palette[i%len(palette)]
			c.Key(task.Title)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
			kit.Stagger(c, root, row, i).Children(func() {
				ui.Text(c, task.Title).FontSize(t.Rem(0.8)).Width(t.Space(24)).Shrink(0).SingleLine()
				track := ui.Box(c).Grow(1).MinWidth(0).Height(t.Space(7)).Radius(t.Radius).Clip().Cursor(ui.CursorResizeEW)
				moved := ui.Local(track, "moved", func() bool { return false })
				dx, _, held := track.Dragged()
				track.Background(track.AnimateColor("bg", t.SurfaceHover.Alpha(0.5).Mix(t.SurfaceHover, unit(track.Hovered())), ui.Snappy))
				fill := track.AnimateColor("ink", ink.Alpha(0.7).Mix(ink, unit(held || *moved)), ui.Snappy)
				lo := track.AnimateTo("lo", float32(task.Start), ui.Snappy)
				hi := track.AnimateTo("hi", float32(task.End), ui.Snappy)
				kit.Settle(c, track, []float32{track.LayoutBounds().W}, func() []float32 { return []float32{track.LayoutBounds().W} })
				ticks := len(o.Ticks)
				track.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
					for k := 1; k < ticks; k++ {
						x := r.X + r.W*float32(k)/float32(ticks)
						p.Line(x, r.Y, x, r.Y+r.H, 1, t.Border.Alpha(0.7))
					}
					inset := t.Space(1)
					bar := ui.Rect{X: r.X + r.W*lo, Y: r.Y + inset, W: r.W * (hi - lo), H: r.H - 2*inset}
					p.Fill(bar, fill, t.Space(1))
					grip := t.Space(0.5)
					p.Fill(ui.Rect{X: bar.X + bar.W - 3*grip, Y: bar.Y + bar.H/4, W: grip, H: bar.H / 2}, t.AccentText.Alpha(0.8), grip/2)
				}))
				if held {
					width := track.LayoutBounds().W
					if dx != 0 && width > 0 {
						cur := (*tasks)[i]
						next := min(1, max(cur.Start+min(0.05, cur.End-cur.Start), cur.End+float64(dx)/float64(width)))
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
