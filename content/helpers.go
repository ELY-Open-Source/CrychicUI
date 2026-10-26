// Package content holds documentation and reading components.
package content

import (
	"fmt"
	"math"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/editor"
)

func def(s *string, d string) {
	if *s == "" {
		*s = d
	}
}

func unit(v bool) float32 {
	if v {
		return 1
	}
	return 0
}

// glyph is a Lucide icon size DIP square.
func glyph(c *ui.Context, name string, size float32) ui.Element {
	return ui.Icon(c, icon.Must(name)).Size(size, size).Shrink(0)
}

// card is the shared arriving card surface.
func card(c *ui.Context) ui.Element {
	return kit.Arrive(surface(c), ui.Pose{Y: c.Theme().Space(2), Fade: 1}, 0)
}

// surface is a card that leaves arrival to its caller.
func surface(c *ui.Context) ui.Element {
	t := c.Theme()
	return ui.Column(c).FillWidth().Radius(t.Radius*1.5).Border(1, t.Border).Background(t.Background).Shadow(0, 1, 3, 0, t.Text.Alpha(0.06))
}

// tile is a tinted icon square.
func tile(c *ui.Context, mark string, ink ui.Color, size float32) ui.Element {
	t := c.Theme()
	box := ui.Box(c).Size(size, size).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).TextColor(ink)
	return box.Children(func() { glyph(c, mark, size*0.5) })
}

// head is a card title row with an icon tile.
func head(c *ui.Context, mark, title string) {
	t := c.Theme()
	ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
		tile(c, mark, t.Accent, t.Space(8))
		ui.Text(c, title).FontSize(t.Rem(0.95)).Bold().Grow(1).MinWidth(0).SingleLine()
	})
}

// unique panics when items repeat or are blank.
func unique(component string, items []string) {
	seen := map[string]bool{}
	for _, s := range items {
		if strings.TrimSpace(s) == "" || seen[s] {
			panic(fmt.Sprintf("content: %s item %q is blank or repeated", component, s))
		}
		seen[s] = true
	}
}

// avatar is a photo or initials circle units wide.
func avatar(c *ui.Context, name string, photo *ui.Bitmap, units float32) ui.Element {
	t := c.Theme()
	return display.Avatar(c, name, photo).Size(t.Space(units), t.Space(units)).FontSize(t.Rem(units * 0.08)).Shrink(0).AlignSelf(ui.Center)
}

// dot joins two meta texts with a drawn dot.
func dot(c *ui.Context) {
	t := c.Theme()
	ui.Box(c).Size(t.Space(0.75), t.Space(0.75)).Shrink(0).Radius(999).Background(t.TextMuted.Alpha(0.6))
}

// meta wraps short facts, each capped to the line.
func meta(c *ui.Context, parts ...string) {
	t := c.Theme()
	row := ui.Row(c).FillWidth().MinWidth(0).Wrap().Gap(t.Space(1.5)).AlignItems(ui.Center)
	seen := row.LayoutBounds().W
	kit.Settle(c, row, []float32{seen}, func() []float32 { return []float32{row.LayoutBounds().W} })
	room := t.Space(40)
	if seen > 0 {
		room = seen
	}
	row.Children(func() {
		first := true
		for _, s := range parts {
			if s == "" {
				continue
			}
			if !first {
				dot(c)
			}
			first = false
			kit.Muted(c, s).FontSize(t.Rem(0.74)).Shrink(0).MaxWidth(room).SingleLine()
		}
	})
}

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

// arc runs clockwise from a0 to a1 degrees.
func arc(cx, cy, r, a0, a1 float32) *ui.Path {
	var p ui.Path
	n := max(2, int(math.Abs(float64(a1-a0))/4))
	for i := 0; i <= n; i++ {
		a := float64(a0+(a1-a0)*float32(i)/float32(n)) * math.Pi / 180
		x, y := cx+r*float32(math.Sin(a)), cy-r*float32(math.Cos(a))
		if i == 0 {
			p.MoveTo(x, y)
		} else {
			p.LineTo(x, y)
		}
	}
	return &p
}

// ring is a circle path starting at the top, clockwise.
func ring(cx, cy, r float32) *ui.Path {
	return arc(cx, cy, r, 0, 360).Close()
}
