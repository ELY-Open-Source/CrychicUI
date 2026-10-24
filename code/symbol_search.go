package code

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SymbolSearchOptions word the field, the empty state and the count.
type SymbolSearchOptions struct {
	Placeholder, Empty string
	// Count words the tally, "%d" taking shown then total.
	Count string
	// Line words a line, "%d" taking the number.
	Line string
}

// SymbolSearch filters symbols as typed, marking the match; cue Select.
func SymbolSearch(c *ui.Context, label string, symbols []Symbol, query *string, chosen *int, o SymbolSearchOptions) (root, field ui.Element) {
	if label == "" || len(symbols) == 0 || query == nil || chosen == nil || *chosen < -1 || *chosen >= len(symbols) {
		panic(fmt.Sprintf("code: SymbolSearch %q needs symbols, a query and a choice or -1", label))
	}
	for _, s := range symbols {
		if s.Kind >= symbolKindCount || s.Name == "" || s.Line < 1 {
			panic(fmt.Sprintf("code: SymbolSearch %q needs a kind, a name and a line in every symbol", label))
		}
	}
	def(&o.Placeholder, "Search symbols")
	def(&o.Empty, "No symbols match")
	def(&o.Count, "%d of %d")
	def(&o.Line, ":%d")
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
			seen := map[string]int{}
			for i, s := range symbols {
				lo, hi := indexFold(s.Name, q)
				if q != "" && lo < 0 {
					continue
				}
				c.Key(ident(seen, fmt.Sprintf("%d %q %d", s.Kind, s.Name, s.Line)))
				on := i == *chosen
				b := kit.Press(ui.ButtonBase(c).Label(s.Name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1.25), t.Space(2)).
					Radius(t.Radius).Checked(on))
				kit.Stagger(c, root, b, hits)
				hits++
				b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
				b.Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						kindBadge(c, s.Kind)
						ui.RichText(c, marked(t, s.Name, lo, hi)...).Font("monospace").FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
						kit.Muted(c, fmt.Sprintf(o.Line, s.Line)).Font("monospace").FontSize(t.Rem(0.75)).Shrink(0)
					})
				})
				if b.Clicked() && !on {
					*chosen = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("SymbolSearch", "selected", "index", i)
				}
			}
			if hits == 0 {
				ui.Column(c).Key("empty").FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(4)).Transition(kit.Enter).Children(func() {
					ui.Icon(c, icon.Must("search-x")).Size(20, 20).TextColor(t.TextMuted)
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.85))
				})
			}
		})
		ui.Row(c).FillWidth().Padding(t.Space(1.5), t.Space(3)).BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Count, hits, len(symbols))).FontSize(t.Rem(0.75)).TextColor(t.TextMuted).
				ContentTransition(ui.ContentTransition{Numeric: true})
		})
	})
	return root, field
}
