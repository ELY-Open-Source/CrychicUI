package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ColumnBrowser shows Finder-style folder columns that slide in; cue Select.
func ColumnBrowser(c *ui.Context, label string, columns [][]string, selected *[]int) ui.Element {
	if label == "" || len(columns) == 0 || selected == nil || len(*selected) != len(columns) {
		panic(fmt.Sprintf("files: ColumnBrowser %q needs columns and one selection each", label))
	}
	for i, column := range columns {
		seen := map[string]bool{}
		for _, name := range column {
			if name == "" || seen[name] {
				panic(fmt.Sprintf("files: ColumnBrowser %q column %d needs unique names", label, i))
			}
			seen[name] = true
		}
		if s := (*selected)[i]; len(column) == 0 || s < -1 || s >= len(column) {
			panic(fmt.Sprintf("files: ColumnBrowser %q column %d needs items and a selection or -1", label, i))
		}
	}
	t := c.Theme()
	root := kit.Card(c).FillWidth().Clip().Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		ui.ScrollHorizontal(c).FillWidth().Height(t.Space(56)).Children(func() {
			for i, column := range columns {
				if i > 0 && (*selected)[i-1] < 0 {
					break
				}
				// Keyed by its folder path: a new one slides in.
				path := make([]string, i)
				for k := range i {
					path[k] = columns[k][(*selected)[k]]
				}
				c.Key(fmt.Sprintf("%q", path))
				col := ui.Scroll(c).Width(t.Space(42)).FillHeight().Shrink(0).Padding(t.Space(1)).BorderWidth(0, 1, 0, 0).BorderColor(t.Border).
					FocusGroup(ui.Vertical).Role(ui.RoleList)
				kit.Arrive(col, ui.Pose{X: t.Space(4), Fade: 1}, 0).Children(func() {
					for j, name := range column {
						c.Key(name)
						on := (*selected)[i] == j
						last := i == len(columns)-1
						row := kit.Press(ui.ButtonBase(c).Label(name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1), t.Space(1.5)).
							Radius(t.Radius).Checked(on))
						row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(unit(row.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy))
						row.TextColor(row.AnimateColor("ink", t.Text.Mix(t.AccentText, unit(on)), ui.Snappy))
						row.Children(func() {
							ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
								glyph := "folder"
								if last {
									glyph = glyphFor(name)
								}
								ui.Icon(c, icon.Must(glyph)).Size(14, 14).Shrink(0)
								ui.Text(c, name).FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
								if !last {
									ui.Icon(c, icon.Must("chevron-right")).Size(12, 12).Shrink(0).Opacity(0.6)
								}
							})
						})
						if row.Clicked() && !on {
							(*selected)[i] = j
							for k := i + 1; k < len(*selected); k++ {
								(*selected)[k] = -1
							}
							root.SetChanged()
							sound.Play(sound.Select)
							kit.Log("ColumnBrowser", "selected", "column", i, "index", j)
						}
					}
				})
			}
		})
	})
	return root
}

// glyphFor picks a file icon, or a folder without extension.
func glyphFor(name string) string {
	if g := glyph(name); g != "file" {
		return g
	}
	return "folder"
}
