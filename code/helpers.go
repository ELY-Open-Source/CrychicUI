// Package code holds code viewing and tooling components.
package code

import (
	"fmt"
	"math"
	"path"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/button"
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

// line returns the spans of a line, tabs as spaces.
func (h *highlighter) line(state uint64, line string) ([]ui.Span, uint64) {
	if h.tok == nil {
		return []ui.Span{{Text: untab(line), Color: h.theme.Text}}, 0
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
		out = append(out, ui.Span{Text: untab(line[at:tk.End]), Color: ink})
		at = tk.End
	}
	if at < len(line) {
		out = append(out, ui.Span{Text: untab(line[at:]), Color: h.theme.Text})
	}
	return out, next
}

func untab(s string) string { return strings.ReplaceAll(s, "\t", "    ") }

// codeBody builds lines; mark styles rows, tail ends them.
func codeBody(c *ui.Context, lang, code string, numbers bool, mark func(i int, row ui.Element), tail func(i int)) {
	t := c.Theme()
	h := newHighlighter(c, lang)
	lines := linesOf(code)
	digits := float32(len(strconv.Itoa(len(lines))))
	state := uint64(0)
	ui.Column(c).FillWidth().Padding(t.Space(2), 0).Children(func() {
		for i, line := range lines {
			c.Key(i)
			var spans []ui.Span
			spans, state = h.line(state, line)
			row := ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(0, t.Space(3))
			if mark != nil {
				mark(i, row)
			}
			row.Children(func() {
				if numbers {
					kit.Numeric(c, strconv.Itoa(i+1)).Font("monospace").FontSize(t.Rem(0.8)).TextColor(t.TextMuted.Alpha(0.7)).
						Width(t.Space(2 * digits)).Shrink(0).TextAlign(ui.End).LineHeight(1.6)
				}
				ui.RichText(c, spans...).Font("monospace").FontSize(t.Rem(0.8)).LineHeight(1.6).Grow(1).MinWidth(0).Wrap()
				if tail != nil {
					tail(i)
				}
			})
		}
	})
}

// pair is a two-way switch; the caller plays cues.
func pair(c *ui.Context, left, right, leftGlyph, rightGlyph string, onLeft bool) (l, r ui.Element) {
	t := c.Theme()
	track := ui.Row(c).Shrink(0).Padding(t.Space(0.5)).Radius(999).Background(t.Surface).Border(1, t.Border).FocusGroup(ui.Horizontal)
	x := track.AnimateTo("thumb", unit(!onLeft), ui.Snappy)
	track.Draw(kit.Inside(999, func(p *ui.Painter, b ui.Rect) {
		pad := t.Space(0.5) + 1
		w := (b.W - 2*pad) / 2
		p.Fill(ui.Rect{X: b.X + pad + w*x, Y: b.Y + pad, W: w, H: b.H - 2*pad}, t.Background, 999)
	}))
	track.Children(func() {
		l = pairSide(c, left, leftGlyph, onLeft)
		r = pairSide(c, right, rightGlyph, !onLeft)
	})
	return l, r
}

func pairSide(c *ui.Context, name, glyph string, on bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(name).Width(t.Space(22)).Padding(t.Space(1), t.Space(2)).Radius(999).Checked(on))
	b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Text, unit(on || b.Hovered())), ui.Snappy))
	return b.Children(func() {
		ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon.Must(glyph)).Size(14, 14).Shrink(0)
			ui.Text(c, name).FontSize(t.Rem(0.8)).Shrink(0)
		})
	})
}

// action is a labeled icon button; callers play its cue.
func action(c *ui.Context, label, glyph string, primary bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius))
	if primary {
		b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())), ui.Snappy)).TextColor(t.AccentText)
	} else {
		b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy)).Border(1, t.Border)
	}
	return b.Children(func() {
		ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon.Must(glyph)).Size(14, 14).Shrink(0)
			ui.Text(c, label).FontSize(t.Rem(0.85)).Bold().Shrink(0)
		})
	})
}

// edited reports a field edit of a tool; cue Tick.
func edited(root ui.Element, component string) {
	root.SetChanged()
	sound.Play(sound.Tick)
	kit.Log(component, "edited")
}

// fieldLabel is the visible caption over an input.
func fieldLabel(c *ui.Context, text string) ui.Element {
	return kit.Muted(c, text).FontSize(c.Theme().Rem(0.75)).Bold()
}

// captioned shows a caption over a field built before it.
func captioned(c *ui.Context, caption string, field func()) {
	t := c.Theme()
	ui.Grid(c).Columns(1).FillWidth().Gap(t.Space(1)).Children(func() {
		ui.Column(c).RowStart(2).ColumnStart(1).MinWidth(0).Children(field)
		fieldLabel(c, caption).RowStart(1).ColumnStart(1)
	})
}

// indexFold finds q in s ignoring case, as byte bounds.
func indexFold(s, q string) (int, int) {
	n := utf8.RuneCountInString(q)
	for i := range s {
		end := i
		for k := 0; k < n && end < len(s); k++ {
			_, w := utf8.DecodeRuneInString(s[end:])
			end += w
		}
		if strings.EqualFold(s[i:end], q) {
			return i, end
		}
	}
	return -1, -1
}

// marked spans s with lo..hi bold in the accent.
func marked(t *ui.Theme, s string, lo, hi int) []ui.Span {
	if lo < 0 || lo == hi {
		return []ui.Span{{Text: s}}
	}
	return []ui.Span{{Text: s[:lo]}, {Text: s[lo:hi], Color: t.Accent, Weight: 700}, {Text: s[hi:]}}
}

// ident keys an item by its text and repeat count.
func ident(seen map[string]int, s string) string {
	seen[s]++
	return s + "#" + strconv.Itoa(seen[s])
}

// rowExit fades a row out; stagger brings rows in.
var rowExit = ui.ElementTransition{Duration: 160 * time.Millisecond, Exit: &ui.Motion{}}

// codeHeader is a file bar: glyph, title, language and actions.
func codeHeader(c *ui.Context, title, lang string, actions func()) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1), t.Space(1), t.Space(3)).
		Background(t.SurfaceHover).BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
		ui.Icon(c, icon.Must(glyphFor(title))).Size(16, 16).Shrink(0).TextColor(t.TextMuted)
		ui.Text(c, title).Font("monospace").FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
		if lang != "" && lang != title {
			langChip(c, lang)
		}
		if actions != nil {
			actions()
		}
	})
}

// langChip is a language name after its colored dot.
func langChip(c *ui.Context, lang string) {
	t := c.Theme()
	ink, _ := langInk(t, lang)
	ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Shrink(0).Children(func() {
		ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(ink)
		kit.Muted(c, lang).FontSize(t.Rem(0.75)).Shrink(0)
	})
}

// langInk is a language's theme color; known reports a match.
func langInk(t *ui.Theme, lang string) (ink ui.Color, known bool) {
	switch strings.ToLower(lang) {
	case "go":
		return t.Accent, true
	case "js", "javascript", "ts", "typescript", "json":
		return t.Warning, true
	case "py", "python", "sh", "bash", "shell":
		return t.Success, true
	case "rs", "rust", "html", "css":
		return t.Danger, true
	}
	return t.TextMuted, false
}

// glyphFor picks a Lucide file icon by extension.
func glyphFor(name string) string {
	switch strings.ToLower(path.Ext(name)) {
	case ".go", ".js", ".ts", ".py", ".rs", ".json", ".html", ".css", ".c", ".java":
		return "file-code"
	case ".sh", ".bash", ".zsh":
		return "file-terminal"
	case ".md", ".txt":
		return "file-text"
	case "":
		return "code"
	}
	return "file"
}

// termChrome is a terminal title bar with three lights.
func termChrome(c *ui.Context, title string, fg ui.Color) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).
		BorderWidth(0, 0, 1, 0).BorderColor(fg.Alpha(0.12)).Children(func() {
		for _, ink := range []ui.Color{t.Danger, t.Warning, t.Success} {
			ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background(ink)
		}
		ui.Text(c, title).FontSize(t.Rem(0.75)).TextColor(fg.Alpha(0.6)).Grow(1).MinWidth(0).TextAlign(ui.Center).SingleLine()
		ui.Box(c).Width(t.Space(9.5)).Shrink(0)
	})
}

// valueInk colors a debugger value by its look.
func valueInk(t *ui.Theme, v string) ui.Color {
	switch {
	case strings.HasPrefix(v, `"`), strings.HasPrefix(v, "'"):
		return t.Success
	case v == "true", v == "false", v == "nil", v == "null":
		return t.Warning
	case v != "" && (v[0] == '-' || v[0] >= '0' && v[0] <= '9'):
		return t.Accent
	}
	return t.Text
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

var (
	iconChevron     = icon.Must("chevron-right")
	iconPR          = icon.Must("git-pull-request")
	iconCheck       = icon.Must("check")
	iconX           = icon.Must("x")
	iconClock       = icon.Must("clock-3")
	iconSpin        = icon.Must("rotate-cw")
	iconCheckCircle = icon.Must("circle-check")
	iconXCircle     = icon.Must("circle-x")
	iconPackage     = icon.Must("package")
	iconCommit      = icon.Must("git-commit-horizontal")
)

// kindBadge is the icon tile of a symbol kind.
func kindBadge(c *ui.Context, kind SymbolKind) ui.Element {
	t := c.Theme()
	glyphs := [4]string{"square-function", "box", "variable", "pi"}
	inks := [4]ui.Color{t.Accent, t.Success, t.Warning, t.TextMuted}
	badge := ui.Box(c).Size(t.Space(5.5), t.Space(5.5)).Shrink(0).Center().Radius(t.Radius).
		Background(inks[kind].Alpha(0.14)).TextColor(inks[kind])
	return badge.Children(func() { ui.Icon(c, icon.Must(glyphs[kind])).Size(14, 14) })
}

// avatarChip is a small initials circle.
func avatarChip(c *ui.Context, name string) ui.Element {
	t := c.Theme()
	initials := ""
	for _, r := range name {
		initials += string(r)
		if len(initials) == 2 {
			break
		}
	}
	e := ui.Box(c).Size(t.Space(8), t.Space(8)).Radius(999).Background(t.SurfaceHover).Center().Label(name)
	return e.Children(func() {
		ui.Text(c, initials).FontSize(t.Rem(0.7)).Bold().TextColor(t.TextMuted)
	})
}

// copyButton copies text; it plays Pop.
func copyButton(c *ui.Context, name, text string) ui.Element {
	return button.CopyButton(c, name, text)
}

func age(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "now"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	case d < 30*24*time.Hour:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	case d < 365*24*time.Hour:
		return fmt.Sprintf("%dmo", max(1, int(d.Hours()/24/30.44)))
	default:
		return fmt.Sprintf("%dy", int(d.Hours()/(24*365)))
	}
}

// took formats a run time as 45s or 1m 12s.
func took(d time.Duration) string {
	switch {
	case d < 10*time.Second:
		return strconv.FormatFloat(d.Seconds(), 'f', 1, 64) + "s"
	case d < time.Minute:
		return strconv.Itoa(int(d.Seconds())) + "s"
	}
	return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
}

// compact formats a count like 1.2k.

func compact(n int) string {
	switch {
	case n < 1000:
		return fmt.Sprintf("%d", n)
	case n < 1000000:
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	default:
		return fmt.Sprintf("%.1fM", float64(n)/1000000)
	}
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

func inversePair(t *ui.Theme) (ui.Color, ui.Color) {
	bg, fg := t.Inverse, t.InverseText
	if bg.A == 0 {
		bg = t.Text
	}
	if fg.A == 0 {
		fg = t.Background
	}
	return bg, fg
}

// t40 is the default panel height.

func t40(c *ui.Context) float32 { return c.Theme().Space(40) }
