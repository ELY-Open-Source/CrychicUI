package code

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SymbolKind is what a symbol is.
type SymbolKind uint8

// Symbol kinds.
const (
	SymFunc SymbolKind = iota
	SymType
	SymVar
	SymConst
	symbolKindCount
)

// Symbol is one symbol of a file.
type Symbol struct {
	Kind SymbolKind
	Name string
	Line int
}

// OutlineOptions name the kind sections and word lines.
type OutlineOptions struct {
	Funcs, Types, Vars, Consts string
	// Line words a line, "%d" taking the number.
	Line string
}

// Outline lists a file's symbols in sections by kind.
func Outline(c *ui.Context, label string, symbols []Symbol, selected *int, o OutlineOptions) ui.Element {
	if label == "" || len(symbols) == 0 || selected == nil || *selected < -1 || *selected >= len(symbols) {
		panic(fmt.Sprintf("code: Outline %q needs a label, symbols and a choice or -1", label))
	}
	var count [symbolKindCount]int
	for _, s := range symbols {
		if s.Kind >= symbolKindCount || s.Name == "" || s.Line < 1 {
			panic(fmt.Sprintf("code: Outline %q needs a kind, a name and a line", label))
		}
		count[s.Kind]++
	}
	def(&o.Funcs, "Functions")
	def(&o.Types, "Types")
	def(&o.Vars, "Variables")
	def(&o.Consts, "Constants")
	def(&o.Line, ":%d")
	names := [symbolKindCount]string{o.Funcs, o.Types, o.Vars, o.Consts}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).FocusGroup(ui.Vertical).Role(ui.RoleList).Label(label)
	return root.Children(func() {
		k := 0
		for _, kind := range []SymbolKind{SymType, SymFunc, SymVar, SymConst} {
			if count[kind] == 0 {
				continue
			}
			c.Key(kind)
			ui.Column(c).FillWidth().Children(func() { k = outlineSection(c, root, kind, names[kind], count[kind], symbols, selected, o.Line, k) })
		}
	})
}

// outlineSection builds a kind's rows; it returns the stagger.
func outlineSection(c *ui.Context, root ui.Element, kind SymbolKind, name string, n int, symbols []Symbol, selected *int, line string, k int) int {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2), t.Space(0.5), t.Space(2)).Children(func() {
		kit.Muted(c, name).FontSize(t.Rem(0.7)).Bold().Grow(1).MinWidth(0).SingleLine()
		kit.Numeric(c, strconv.Itoa(n)).FontSize(t.Rem(0.7)).TextColor(t.TextMuted)
	})
	seen := map[string]int{}
	for i, s := range symbols {
		if s.Kind != kind {
			continue
		}
		c.Key(ident(seen, fmt.Sprintf("%q %d", s.Name, s.Line)))
		on := i == *selected
		b := kit.Press(ui.ButtonBase(c).Label(s.Name).FillWidth().AlignItems(ui.Stretch).Padding(t.Space(1), t.Space(2)).
			Radius(t.Radius).Checked(on))
		stagger(c, root, b, k)
		k++
		b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())).Mix(t.Accent.Alpha(0.12), unit(on)), ui.Snappy))
		b.Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				kindBadge(c, s.Kind)
				txt := ui.Text(c, s.Name).Font("monospace").FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
				if on {
					txt.Bold()
				}
				kit.Muted(c, fmt.Sprintf(line, s.Line)).Font("monospace").FontSize(t.Rem(0.75)).Shrink(0)
			})
		})
		if b.Clicked() && !on {
			*selected = i
			root.SetChanged()
			sound.Play(sound.Select)
			kit.Log("Outline", "selected", "index", i)
		}
	}
	return k
}
