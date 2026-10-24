package project

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ChecklistItem is one editable checklist item.
type ChecklistItem struct {
	Title string
	Done  bool
}

// ChecklistOptions are the checklist labels.
type ChecklistOptions struct{ Add, Toggle string }

// Checklist edits a checklist; cue Toggle, Pop.
func Checklist(c *ui.Context, items *[]ChecklistItem, o ChecklistOptions) ui.Element {
	if items == nil {
		panic("project: Checklist needs items")
	}
	def(&o.Add, "Add item")
	def(&o.Toggle, "Toggle")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, item := range *items {
			if item.Title == "" {
				panic("project: Checklist item must not be empty")
			}
			c.Key(i)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center)
			last := ui.Local(row, "last", func() string { return item.Title })
			row.Children(func() {
				toggle := kit.Press(ui.ButtonBase(c).Label(o.Toggle).Size(t.Space(6), t.Space(6)).Center().Radius(999).Checked(item.Done))
				toggle.Children(func() {
					if item.Done {
						ui.Icon(c, icon.Must("check")).Size(16, 16).TextColor(t.Success)
					}
				})
				if toggle.Clicked() {
					(*items)[i].Done = !(*items)[i].Done
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("Checklist", "toggled", "index", i, "done", (*items)[i].Done)
				}
				edit := ui.EditableText(c, &(*items)[i].Title).Grow(1)
				if (*items)[i].Done {
					edit.Strikethrough()
				}
				if edit.Changed() {
					if strings.TrimSpace((*items)[i].Title) == "" {
						(*items)[i].Title = *last
						kit.Log("Checklist", "empty rejected")
					} else {
						*last = (*items)[i].Title
						root.SetChanged()
						sound.Play(sound.Pop)
						kit.Log("Checklist", "edited", "index", i)
					}
				} else {
					*last = (*items)[i].Title
				}
			})
		}
		c.Key("add")
		add := kit.Press(ui.ButtonBase(c).Label(o.Add).Padding(t.Space(1), t.Space(2)).Radius(t.Radius))
		add.Children(func() { ui.Text(c, o.Add).FontSize(t.Rem(0.75)).TextColor(t.TextMuted) })
		if add.Clicked() {
			*items = append(*items, ChecklistItem{Title: fmt.Sprintf("%s %d", o.Add, len(*items)+1)})
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("Checklist", "added", "count", len(*items))
		}
	})
	return root
}
