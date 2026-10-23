// Package code holds code viewing and tooling components.
package code

import (
	"fmt"
	"math"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

// def fills an empty string with its default.
func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(b bool) float32 {
	if b {
		return 1
	}
	return 0
}

// finite reports whether v is a usable size.
func finite(v float32) bool { return !math.IsNaN(float64(v)) && !math.IsInf(float64(v), 0) }

// size panics on a bad size; zero becomes d.
func size(component string, v *float32, d float32) {
	if *v < 0 || !finite(*v) {
		panic(fmt.Sprintf("code: %s sizes must be finite and not negative", component))
	}
	if *v == 0 {
		*v = d
	}
}

// linesOf splits code into lines, keeping blanks.
func linesOf(code string) []string {
	lines := strings.Split(code, "\n")
	if n := len(lines); n > 1 && lines[n-1] == "" {
		lines = lines[:n-1]
	}
	return lines
}

// round is an icon button; the caller plays its cue.
func round(c *ui.Context, name, glyph string) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999))
	b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
	return b.Children(func() {
		ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
	})
}

// released plays Tick after a drag or key move.
func released(s ui.Element, component string, value float64) {
	moved := ui.Local(s, "moved", func() bool { return false })
	*moved = *moved || s.Changed()
	if _, _, held := s.Dragged(); *moved && !held {
		*moved = false
		sound.Play(sound.Tick)
		kit.Log(component, "changed", "value", value)
	}
}

// highlighter tokenizes lines of one language.
type highlighter struct {
	tok   editor.Tokenizer
	theme *editor.Theme
}

func newHighlighter(c *ui.Context, lang string) *highlighter {
	h := &highlighter{theme: editor.DefaultTheme(c.Theme())}
	if l := editor.LanguageNamed(lang); l != nil {
		h.tok = l.Tokenizer
	}
	return h
}

// line returns the spans of one line of code.
func (h *highlighter) line(state uint64, line string) ([]ui.Span, uint64) {
	if h.tok == nil {
		return []ui.Span{{Text: line, Color: h.theme.Text}}, 0
	}
	toks, next := h.tok.Tokenize(line, state, nil)
	out := make([]ui.Span, 0, len(toks)+1)
	at := 0
	for _, tk := range toks {
		if tk.End > len(line) || tk.End <= at {
			continue
		}
		ink := h.theme.Tokens[tk.Kind].Color
		if ink.A == 0 {
			ink = h.theme.Text
		}
		out = append(out, ui.Span{Text: line[at:tk.End], Color: ink})
		at = tk.End
	}
	if at < len(line) {
		out = append(out, ui.Span{Text: line[at:], Color: h.theme.Text})
	}
	return out, next
}

// codeLines builds highlighted mono lines of code.
func codeLines(c *ui.Context, lang, code string, mark func(i int, e ui.Element)) {
	t := c.Theme()
	h := newHighlighter(c, lang)
	state := uint64(0)
	for i, line := range linesOf(code) {
		c.Key(i)
		var spans []ui.Span
		spans, state = h.line(state, line)
		row := ui.RichText(c, spans...).Font("monospace").FontSize(t.Rem(0.85))
		if mark != nil {
			mark(i, row)
		}
	}
}

// StackFrame is one frame of a trace or stack.
type StackFrame struct {
	Func, File string
	Line       int
	App        bool
}

// frames panics on a frame missing data.
func frames(component string, list []StackFrame) {
	for _, f := range list {
		if f.Func == "" || f.File == "" || f.Line < 0 {
			panic(fmt.Sprintf("code: %s frames need a function, a file and a line", component))
		}
	}
}

// kindBadge is the letter chip of a symbol kind.
func kindBadge(c *ui.Context, kind SymbolKind) ui.Element {
	t := c.Theme()
	letters := [4]string{"f", "T", "v", "c"}
	inks := [4]ui.Color{t.Accent, t.Success, t.Warning, t.TextMuted}
	badge := ui.Box(c).Size(t.Space(4.5), t.Space(4.5)).Shrink(0).Center().Radius(t.Radius).
		Background(inks[kind].Alpha(0.15)).TextColor(inks[kind])
	return badge.Children(func() {
		ui.Text(c, letters[kind]).Font("monospace").FontSize(t.Rem(0.75)).Bold()
	})
}

// watchEditor reports buffer changes of an editor state.
func watchEditor(root ui.Element, st *editor.State, component string) {
	w := ui.Local(root, "watch", func() *editorWatch { return &editorWatch{} })
	if (*w).state != st {
		(*w).state = st
		(*w).version = st.Buffer.Version()
		return
	}
	if version := st.Buffer.Version(); version != (*w).version {
		(*w).version = version
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log(component, "edited")
	}
}

type editorWatch struct {
	state   *editor.State
	version uint64
}
