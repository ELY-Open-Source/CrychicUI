// Code components for the gallery.
package code

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"strings"
)

type BranchPickerOptions struct {
	Placeholder string
	Empty       string
}

// BranchPicker is branches with search; cue Select.

func BranchPicker(c *ui.Context, label string, branches []string, current *int, query *string, o BranchPickerOptions) (root, field ui.Element) {
	if label == "" || len(branches) == 0 || current == nil || query == nil || *current < 0 || *current >= len(branches) {
		panic("code: BranchPicker needs a label, branches and bound values")
	}
	def(&o.Placeholder, "Search branches")
	def(&o.Empty, "No branches match")
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Start).Width(t.Space(72)).Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		field = input.SearchInput(c, query, input.SearchInputOptions{Placeholder: o.Placeholder})
		if field.Changed() {
			root.SetChanged()
		}
		q := strings.ToLower(*query)
		hits := 0
		list := ui.Column(c).FillWidth().Gap(0)
		list.Children(func() {
			for i, b := range branches {
				if b == "" {
					panic("code: BranchPicker needs a name for every branch")
				}
				if q != "" && !strings.Contains(strings.ToLower(b), q) {
					continue
				}
				hits++
				c.Key(b)
				row := kit.Press(ui.ButtonBase(c).Label(b).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
					Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Checked(i == *current))
				bg := t.SurfaceHover.Alpha(unit(row.Hovered()))
				if i == *current {
					bg = t.Accent.Alpha(0.12)
				}
				row.Background(row.AnimateColor("bg", bg, ui.Snappy))
				row.Children(func() {
					ui.Icon(c, icon.Must("git-branch")).Size(16, 16).TextColor(t.TextMuted)
					ui.Text(c, b).Font("monospace").FontSize(t.Rem(0.85)).SingleLine().Grow(1)
					if i == *current {
						ui.Icon(c, iconCheck).Size(16, 16).TextColor(t.Accent)
					}
				})
				if row.Clicked() && *current != i {
					*current = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("BranchPicker", "selected", "index", i)
				}
			}
		})
		if hits == 0 {
			kit.Muted(c, o.Empty).FontSize(t.Rem(0.85)).Padding(t.Space(1), t.Space(2))
		}
	})
	return root, field
}

// LangShare is one language's share of a codebase.
