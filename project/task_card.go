package project

import (
	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TaskCard shows title, labels, assignee, and due; it lifts.
func TaskCard(c *ui.Context, task Task) ui.Element {
	validTask("TaskCard", task)
	t := c.Theme()
	root := kit.Lift(kit.Card(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Role(ui.RoleGroup))
	root.Background(root.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, 0.5*unit(root.Hovered())), ui.Snappy))
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0).Children(func() { taskFace(c, task) })
	return root
}

// taskFace is a task's labels, title, avatar and due line.
func taskFace(c *ui.Context, task Task) {
	t := c.Theme()
	if len(task.Labels) > 0 {
		row := ui.Row(c).FillWidth().Wrap().Gap(t.Space(1))
		room := row.LayoutBounds().W
		kit.Settle(c, row, []float32{room}, func() []float32 { return []float32{row.LayoutBounds().W} })
		row.Children(func() {
			for i, label := range task.Labels {
				c.Key(i)
				fitTag(c, label, labelInk(t, label), room)
			}
		})
	}
	title := ui.Text(c, task.Title).FontSize(t.Rem(0.85)).Bold().Wrap().FillWidth()
	if task.Done {
		title.Strikethrough().TextColor(t.TextMuted)
	}
	if task.Assignee == "" && task.Due == "" {
		return
	}
	ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
		if task.Assignee != "" {
			display.Avatar(c, task.Assignee, nil).Size(t.Space(6), t.Space(6)).FontSize(t.Rem(0.6))
			kit.Muted(c, task.Assignee).FontSize(t.Rem(0.7)).Grow(1).MinWidth(0).SingleLine()
		} else {
			ui.Box(c).Grow(1)
		}
		if task.Due != "" {
			ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
				glyph(c, "calendar", 12)
				ui.Text(c, task.Due).FontSize(t.Rem(0.7)).Shrink(0)
			})
		}
	})
}
