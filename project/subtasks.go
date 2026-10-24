package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Subtask is one checklist item.
type Subtask struct {
	Title string
	Done  bool
}

// Subtasks shows a checklist with done count bar; cue Toggle.
func Subtasks(c *ui.Context, items *[]Subtask) ui.Element {
	if items == nil || len(*items) == 0 {
		panic("project: Subtasks needs items")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	done := 0
	for _, item := range *items {
		if item.Done {
			done++
		}
	}
	root.Children(func() {
		kit.Numeric(c, fmt.Sprintf("%d/%d", done, len(*items))).Bold()
		progress(c, float64(done)/float64(len(*items)), t.Success)
		for i, item := range *items {
			c.Key(i)
			row := kit.Press(ui.ButtonBase(c).Label(item.Title).Row().FillWidth().Gap(t.Space(2)).Padding(t.Space(1)).Radius(t.Radius).
				Checked(item.Done))
			row.Children(func() {
				text := ui.Text(c, item.Title)
				if item.Done {
					text.Strikethrough()
				}
			})
			if row.Clicked() {
				(*items)[i].Done = !(*items)[i].Done
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("Subtasks", "toggled", "index", i, "done", (*items)[i].Done)
			}
		}
	})
	return root
}
