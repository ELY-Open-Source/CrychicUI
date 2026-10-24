package project

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Backlog orders issues by drag; cue Drop.
func Backlog(c *ui.Context, issues *[]string) ui.Element {
	if issues == nil || len(*issues) == 0 {
		panic("project: Backlog needs issues")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	gridID := ui.Local(root, "gridID", func() *int { return new(int) })
	moveTitle, moveTo := "", -1
	root.Children(func() {
		for i, title := range *issues {
			if title == "" {
				panic("project: Backlog issue must not be empty")
			}
			c.Key(title)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
				Border(1, t.Border).Focusable().Label(title).Cursor(ui.CursorGrab).Drag(backlogDrag{grid: *gridID, title: title})
			row.Children(func() {
				ui.Icon(c, icon.Must("grip-vertical")).Size(16, 16).TextColor(t.TextMuted)
				ui.Text(c, title).Grow(1)
			})
			if d, ok := ui.Drop[backlogDrag](row); ok && d.grid == *gridID {
				moveTitle, moveTo = d.title, i
			}
		}
	})
	if moveTitle != "" {
		from := slices.Index(*issues, moveTitle)
		if from >= 0 && from != moveTo {
			*issues = slices.Delete(*issues, from, from+1)
			*issues = slices.Insert(*issues, moveTo, moveTitle)
			root.SetChanged()
			sound.Play(sound.Drop)
			kit.Log("Backlog", "moved", "from", from, "to", moveTo)
		}
	}
	return root
}

type backlogDrag struct {
	grid  *int
	title string
}
