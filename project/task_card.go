package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TaskCard shows title, labels, assignee, and due; it lifts.
func TaskCard(c *ui.Context, task Task) ui.Element {
	if task.ID == "" || task.Title == "" || task.Assignee == "" || task.Due == "" {
		panic("project: TaskCard has an invalid task")
	}
	t := c.Theme()
	root := kit.Lift(ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).Background(t.Surface))
	root.Children(func() {
		ui.Text(c, task.Title).Bold()
		ui.Row(c).Wrap().Gap(t.Space(1)).Children(func() {
			for _, label := range task.Labels {
				ui.Text(c, label).Padding(t.Space(0.5), t.Space(1)).Radius(999).Background(t.Accent.Alpha(0.12)).FontSize(t.Rem(0.7))
			}
		})
		ui.Row(c).FillWidth().Children(func() {
			kit.Muted(c, task.Assignee).FontSize(t.Rem(0.75)).Grow(1)
			kit.Muted(c, task.Due).FontSize(t.Rem(0.75))
		})
	})
	return root
}
