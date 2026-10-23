package code

import (
	"fmt"

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

// Outline lists the symbols of a file by kind.
func Outline(c *ui.Context, label string, symbols []Symbol, selected *int) ui.Element {
	if label == "" || len(symbols) == 0 || selected == nil || *selected < -1 || *selected >= len(symbols) {
		panic(fmt.Sprintf("code: Outline %q needs a label, symbols and a choice or -1", label))
	}
	for _, s := range symbols {
		if s.Kind >= symbolKindCount || s.Name == "" || s.Line < 1 {
			panic(fmt.Sprintf("code: Outline %q needs a kind, a name and a line", label))
		}
	}
	t := c.Theme()
	letters := [4]string{"f", "T", "v", "c"}
	inks := [4]ui.Color{t.Accent, t.Success, t.Warning, t.TextMuted}
	root := ui.Column(c).FillWidth().Gap(t.Space(0.5)).FocusGroup(ui.Vertical).Role(ui.RoleList).Label(label)
	return root.Children(func() {
		for i, s := range symbols {
			c.Key(i)
			on := i == *selected
			b := kit.Press(ui.ButtonBase(c).Label(s.Name).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(t.Space(1), t.Space(2)).Radius(t.Radius).Checked(on))
			bg := t.SurfaceHover.Alpha(unit(b.Hovered()))
			if on {
				bg = t.Accent.Alpha(0.1)
			}
			b.Background(b.AnimateColor(fmt.Sprint("bg.", i), bg, ui.Snappy)).Children(func() {
				badge := ui.Box(c).Size(t.Space(4.5), t.Space(4.5)).Shrink(0).Center().Radius(t.Radius).Background(inks[s.Kind].Alpha(0.15)).TextColor(inks[s.Kind])
				badge.Children(func() { ui.Text(c, letters[s.Kind]).Font("monospace").FontSize(t.Rem(0.75)).Bold() })
				ui.Text(c, s.Name).Font("monospace").FontSize(t.Rem(0.85)).Grow(1).SingleLine()
				kit.Muted(c, fmt.Sprintf(":%d", s.Line)).Font("monospace").FontSize(t.Rem(0.8))
			})
			if b.Clicked() && !on {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Outline", "selected", "index", i)
			}
		}
	})
}
