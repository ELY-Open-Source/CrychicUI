package project

import (
	"fmt"
	"slices"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ChecklistItem is one editable checklist item; ID is unique.
type ChecklistItem struct {
	ID, Title string
	Done      bool
}

// ChecklistOptions word the card; Toggle formats each box's label.
type ChecklistOptions struct{ Title, Count, Add, NewItem, Toggle, Empty string }

// Checklist edits a checklist; cue Toggle, Pop.
func Checklist(c *ui.Context, items *[]ChecklistItem, o ChecklistOptions) ui.Element {
	if items == nil {
		panic("project: Checklist needs items")
	}
	def(&o.Title, "Checklist")
	def(&o.Count, "%d of %d")
	def(&o.Add, "Add item")
	def(&o.NewItem, "New item")
	def(&o.Toggle, "Toggle %s")
	def(&o.Empty, "No items yet")
	done := 0
	ids := make([]string, len(*items))
	for i, item := range *items {
		if item.Title == "" {
			panic(fmt.Sprintf("project: Checklist %q item %q is empty", o.Title, item.ID))
		}
		ids[i] = item.ID
		if item.Done {
			done++
		}
	}
	unique("Checklist", ids)
	t := c.Theme()
	share, ink := 0.0, t.Accent
	if len(*items) > 0 {
		share = float64(done) / float64(len(*items))
	}
	if len(*items) > 0 && done == len(*items) {
		ink = t.Success
	}
	root := card(c)
	made := ui.Local(root, "made", func() int { return 0 })
	root.Children(func() {
		kit.CardHead(c, "list-checks", o.Title, func() {
			kit.Chip(c, fmt.Sprintf(o.Count, done, len(*items)), ink)
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(2.5), t.Space(2))
		list.Children(func() {
			ui.Box(c).FillWidth().Padding(0, t.Space(1), t.Space(1.5), t.Space(1)).Children(func() { progress(c, share, ink) })
			if len(*items) == 0 {
				kit.Muted(c, o.Empty).FontSize(t.Rem(0.8)).Padding(t.Space(1.5), t.Space(1))
			}
			c.Key("rows")
			rows := ui.Column(c).FillWidth().Gap(t.Space(0.5))
			rows.Children(func() {
				for i, item := range *items {
					c.Key(item.ID)
					row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1)).Radius(t.Radius)
					last := ui.Local(row, "last", func() string { return item.Title })
					kit.Stagger(c, rows, row, i).Children(func() {
						box := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Toggle, item.Title)).Size(t.Space(5), t.Space(5)).Shrink(0).
							Center().Radius(t.Space(1.25)).Checked(item.Done))
						box.Background(box.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(box.Hovered())).Mix(t.Success, unit(item.Done)), ui.Snappy))
						box.Border(1, box.AnimateColor("line", t.Border.Mix(t.Success, unit(item.Done)), ui.Snappy))
						box.Children(func() {
							if item.Done {
								kit.Arrive(glyph(c, "check", 14).TextColor(t.AccentText), ui.Pose{Scale: 0.4, Fade: 1}, 0)
							}
						})
						if box.Clicked() {
							(*items)[i].Done = !(*items)[i].Done
							root.SetChanged()
							sound.Play(sound.Toggle)
							kit.Log("Checklist", "toggled", "index", i, "done", (*items)[i].Done)
						}
						edit := ui.EditableText(c, &(*items)[i].Title).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0)
						edit.TextColor(edit.AnimateColor("ink", t.Text.Mix(t.TextMuted, unit(item.Done)), ui.Snappy))
						if item.Done {
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
			})
			c.Key("add")
			add := kit.Press(ui.ButtonBase(c).Label(o.Add).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.5), t.Space(1)).Radius(t.Radius))
			add.Background(add.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(add.Hovered())), ui.Snappy))
			add.TextColor(add.AnimateColor("ink", t.TextMuted.Mix(t.Accent, unit(add.Hovered())), ui.Snappy))
			add.Children(func() {
				ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					glyph(c, "plus", 16)
					ui.Text(c, o.Add).FontSize(t.Rem(0.8)).Shrink(0)
				})
			})
			if add.Clicked() {
				id := ""
				for id == "" || slices.Contains(ids, id) {
					*made++
					id = fmt.Sprint("new-", *made)
				}
				*items = append(*items, ChecklistItem{ID: id, Title: o.NewItem})
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("Checklist", "added", "count", len(*items))
			}
		})
	})
	return root
}
