// Code components for the gallery.
package code

import (
	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"strconv"
	"strings"
)

type SymbolSearchOptions struct {
	Placeholder string
	Empty       string
}

// SymbolSearch filters symbols as typed; cue Select.

func SymbolSearch(c *ui.Context, label string, symbols []Symbol, query *string, chosen *int, o SymbolSearchOptions) (root, field ui.Element) {
	if label == "" || len(symbols) == 0 || query == nil || chosen == nil || *chosen < -1 || *chosen >= len(symbols) {
		panic("code: SymbolSearch needs a label, symbols and bound values")
	}
	def(&o.Placeholder, "Search symbols")
	def(&o.Empty, "No symbols match")
	t := c.Theme()
	root = ui.Column(c).AlignSelf(ui.Start).Width(t.Space(80)).Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		field = input.SearchInput(c, query, input.SearchInputOptions{Placeholder: o.Placeholder})
		if field.Changed() {
			root.SetChanged()
		}
		q := strings.ToLower(*query)
		hits := 0
		list := ui.Column(c).FillWidth().Gap(0)
		list.Children(func() {
			for i, s := range symbols {
				if s.Kind >= symbolKindCount || s.Name == "" || s.Line < 1 {
					panic("code: SymbolSearch needs a kind, a name and a line in every symbol")
				}
				if q != "" && !strings.Contains(strings.ToLower(s.Name), q) {
					continue
				}
				hits++
				c.Key(i)
				row := kit.Press(ui.ButtonBase(c).Label(s.Name).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
					Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Checked(i == *chosen))
				bg := t.SurfaceHover.Alpha(unit(row.Hovered()))
				if i == *chosen {
					bg = t.Accent.Alpha(0.12)
				}
				row.Background(row.AnimateColor("bg", bg, ui.Snappy))
				row.Children(func() {
					kindBadge(c, s.Kind)
					ui.Text(c, s.Name).SingleLine().Grow(1)
					kit.Muted(c, strconv.Itoa(s.Line)).FontSize(t.Rem(0.8)).Font("monospace")
				})
				if row.Clicked() && *chosen != i {
					*chosen = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("SymbolSearch", "selected", "index", i)
				}
			}
		})
		if hits == 0 {
			kit.Muted(c, o.Empty).FontSize(t.Rem(0.85)).Padding(t.Space(1), t.Space(2))
		}
	})
	return root, field
}

// BranchPickerOptions are the copy knobs of BranchPicker.
