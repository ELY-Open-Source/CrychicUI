package code

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Branch is a branch and its distance from its base.
type Branch struct {
	Name          string
	Ahead, Behind int
}

// BranchPickerOptions are the copy knobs of BranchPicker.
type BranchPickerOptions struct {
	Placeholder, Empty, Current string
}

// BranchPicker lists branches to search and pick; cue Select.
func BranchPicker(c *ui.Context, label string, branches []Branch, current *int, query *string, o BranchPickerOptions) (root, field ui.Element) {
	if label == "" || len(branches) == 0 || current == nil || query == nil || *current < 0 || *current >= len(branches) {
		panic(fmt.Sprintf("code: BranchPicker %q needs branches and bound values", label))
	}
	seen := map[string]bool{}
	for _, b := range branches {
		if b.Name == "" || seen[b.Name] || b.Ahead < 0 || b.Behind < 0 {
			panic(fmt.Sprintf("code: BranchPicker %q needs unique names and counts not below zero", label))
		}
		seen[b.Name] = true
	}
	def(&o.Placeholder, "Search branches")
	def(&o.Empty, "No branches match")
	def(&o.Current, "current")
	t := c.Theme()
	root = kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		ui.Row(c).FillWidth().Padding(t.Space(2)).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			field = input.SearchInput(c, query, input.SearchInputOptions{Placeholder: o.Placeholder}).Grow(1)
		})
		if field.Changed() {
			root.SetChanged()
		}
		q := strings.TrimSpace(*query)
		hits := 0
		ui.Column(c).FillWidth().Padding(t.Space(1)).FocusGroup(ui.Vertical).Role(ui.RoleList).Children(func() {
			for i, b := range branches {
				lo, hi := indexFold(b.Name, q)
				if q != "" && lo < 0 {
					continue
				}
				c.Key(b.Name)
				on := i == *current
				row := kit.Press(ui.ButtonBase(c).Label(b.Name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.25), t.Space(2)).
					Radius(t.Radius).Checked(on))
				kit.Stagger(c, root, row, hits)
				hits++
				row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(unit(row.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
				row.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						ink := t.TextMuted
						if on {
							ink = t.Accent
						}
						ui.Icon(c, icon.Must("git-branch")).Size(14, 14).Shrink(0).TextColor(ink)
						ui.RichText(c, marked(t, b.Name, lo, hi)...).Font("monospace").FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
						if b.Ahead > 0 {
							kit.Numeric(c, "↑"+strconv.Itoa(b.Ahead)).FontSize(t.Rem(0.7)).TextColor(t.Success).Shrink(0)
						}
						if b.Behind > 0 {
							kit.Numeric(c, "↓"+strconv.Itoa(b.Behind)).FontSize(t.Rem(0.7)).TextColor(t.Warning).Shrink(0)
						}
						if on {
							kit.Chip(c, o.Current, t.Accent)
						}
					})
				})
				if row.Clicked() && !on {
					*current = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("BranchPicker", "selected", "index", i)
				}
			}
			if hits == 0 {
				ui.Column(c).Key("empty").FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(4)).Transition(kit.Enter).Children(func() {
					ui.Icon(c, icon.Must("search-x")).Size(20, 20).TextColor(t.TextMuted)
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
				})
			}
		})
	})
	return root, field
}
