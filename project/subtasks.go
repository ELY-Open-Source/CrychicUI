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

// SubtasksOptions word the list; Count takes done and total.
type SubtasksOptions struct{ Title, Count string }

// Subtasks shows a checklist with done count bar; cue Toggle.
func Subtasks(c *ui.Context, items *[]Subtask, o SubtasksOptions) ui.Element {
	if items == nil || len(*items) == 0 {
		panic("project: Subtasks needs items")
	}
	names := make([]string, len(*items))
	done := 0
	for i, item := range *items {
		names[i] = item.Title
		if item.Done {
			done++
		}
	}
	unique("Subtasks", names)
	def(&o.Title, "Subtasks")
	def(&o.Count, "%d of %d")
	t := c.Theme()
	all := done == len(*items)
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "list-checks", o.Title, func() {
			ink := t.TextMuted
			if all {
				ink = t.Success
			}
			kit.Chip(c, fmt.Sprintf(o.Count, done, len(*items)), ink)
		})
		ui.Column(c).FillWidth().Padding(t.Space(3), t.Space(3), t.Space(1), t.Space(3)).Children(func() {
			progress(c, float64(done)/float64(len(*items)), t.Accent.Mix(t.Success, unit(all)))
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(1))
		list.Children(func() {
			for i, item := range *items {
				c.Key(item.Title)
				row := kit.Press(ui.ButtonBase(c).Label(item.Title).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.5), t.Space(2)).
					Radius(t.Radius).Checked(item.Done))
				row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(0.7*unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
						box := ui.Box(c).Size(t.Space(4.5), t.Space(4.5)).Shrink(0).Center().Radius(t.Space(1)).TextColor(t.AccentText)
						box.Border(1, box.AnimateColor("line", t.Border.Mix(t.Success, unit(item.Done)), ui.Snappy))
						box.Background(box.AnimateColor("fill", t.Success.Alpha(unit(item.Done)), ui.Snappy)).Children(func() {
							if item.Done {
								c.Key("check")
								kit.Arrive(glyph(c, "check", 12), ui.Pose{Scale: 0.4, Fade: 1}, 0)
							}
						})
						text := ui.Text(c, item.Title).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).Wrap().
							TextColor(row.AnimateColor("ink", t.Text.Mix(t.TextMuted, unit(item.Done)), ui.Snappy))
						if item.Done {
							text.Strikethrough()
						}
					})
				})
				if row.Clicked() {
					(*items)[i].Done = !item.Done
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("Subtasks", "toggled", "index", i, "done", !item.Done)
				}
			}
		})
	})
	return root
}
