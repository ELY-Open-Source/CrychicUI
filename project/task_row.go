package project

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// TaskRow checks a task; the title strikes; cue Toggle.
func TaskRow(c *ui.Context, task Task, done *bool) ui.Element {
	if task.ID == "" || task.Title == "" || task.Assignee == "" || task.Due == "" || done == nil {
		panic("project: TaskRow has an invalid task")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).Radius(t.Radius)
	root.Children(func() {
		check := kit.Press(ui.ButtonBase(c).Label(task.Title).Size(t.Space(5), t.Space(5)).Center().Radius(999).
			Border(1, t.Border).Checked(*done))
		check.Children(func() {
			if *done {
				ui.Icon(c, icon.Must("check")).Size(16, 16)
			}
		})
		if check.Clicked() {
			*done = !*done
			root.SetChanged()
			sound.Play(sound.Toggle)
			kit.Log("TaskRow", "toggled", "done", *done)
		}
		text := ui.Text(c, task.Title).Grow(1)
		if *done {
			text.Strikethrough()
		}
		kit.Muted(c, task.Assignee).FontSize(t.Rem(0.75))
		kit.Muted(c, task.Due).FontSize(t.Rem(0.75))
	})
	return root
}
