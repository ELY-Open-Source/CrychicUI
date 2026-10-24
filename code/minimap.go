package code

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MinimapOptions tune a Minimap.
type MinimapOptions struct {
	// Height of the strip; zero gives the default.
	Height float32
}

// mapCols is how many columns the strip's width shows.
const mapCols = 80

// mapRun is one colored token of a minimap line.
type mapRun struct {
	line, col, n int
	ink          ui.Color
}

// Minimap draws a file's tokens small under a viewport box.
func Minimap(c *ui.Context, label, lang, code string, visible int, at *int, o MinimapOptions) ui.Element {
	lines := linesOf(code)
	n := len(lines)
	if label == "" || code == "" || visible < 1 || visible > n || at == nil || *at < 0 || *at > n-visible {
		panic(fmt.Sprintf("code: Minimap %q needs code, a window and a start, got %d of %d", label, deref(at), n))
	}
	t := c.Theme()
	size("Minimap", &o.Height, t.Space(60))
	root := kit.Press(ui.ButtonBase(c).Label(label).Width(t.Space(24)).Height(o.Height).Radius(t.Radius).Clip().
		Background(t.Surface).Border(1, t.Border).Role(ui.RoleSlider).Range(0, float64(n-visible), float64(*at)))
	height := ui.Local(root, "height", func() float32 { return o.Height })
	pressY := ui.Local(root, "pressY", func() float32 { return -1 })
	dragOffset := ui.Local(root, "dragOffset", func() float32 { return 0 })
	kit.Settle(c, root, []float32{*height}, func() []float32 { return []float32{root.LayoutBounds().H} })
	if got := root.LayoutBounds().H; got > 0 {
		*height = got
	}
	runs := mapRuns(c, lang, lines)
	frac := root.AnimateTo("at", float32(*at)/float32(n), ui.Snappy)
	span := float32(visible) / float32(n)
	root.Draw(kit.Inside(t.Radius, func(p *ui.Painter, r ui.Rect) {
		pad := t.Space(1.5)
		lh := r.H / float32(n)
		cw := (r.W - 2*pad) / mapCols
		for _, run := range runs {
			w := min(float32(run.n), float32(mapCols-run.col)) * cw
			if w > 0 {
				p.Fill(ui.Rect{X: r.X + pad + float32(run.col)*cw, Y: r.Y + float32(run.line)*lh + lh*0.15, W: w, H: max(lh*0.7, 1)}, run.ink.Alpha(0.7), 0)
			}
		}
		box := ui.Rect{X: r.X, Y: r.Y + r.H*frac, W: r.W, H: max(r.H*span, t.Space(2))}
		p.Fill(box, t.Accent.Alpha(0.14), 0)
		p.Stroke(box, t.Accent.Alpha(0.5), 0, 1)
	}))
	root.HandleInput(func(ev ui.InputEvent) bool {
		if ev.Kind == ui.InputPointerDown && ev.Button == 0 {
			*pressY = ev.Y
		}
		return false
	})
	if *pressY >= 0 {
		if moved := clampAt(int(*pressY / *height * float32(n))-visible/2, n, visible); moved != *at {
			*at = moved
			root.SetChanged()
		}
		*pressY = -1
	}
	if _, dy, held := root.Dragged(); held {
		*dragOffset += dy * float32(n) / *height
		steps := int(*dragOffset)
		*dragOffset -= float32(steps)
		if moved := clampAt(*at+steps, n, visible); moved != *at {
			*at = moved
			root.SetChanged()
		}
	} else {
		*dragOffset = 0
	}
	released(root, "Minimap", float64(*at))
	move := func(to int, event string) {
		if to = clampAt(to, n, visible); to == *at {
			return
		}
		*at = to
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("Minimap", event, "at", *at)
	}
	switch {
	case root.Shortcut(0, ui.KeyDown), root.Shortcut(0, ui.KeyRight):
		move(*at+1, "down")
	case root.Shortcut(0, ui.KeyUp), root.Shortcut(0, ui.KeyLeft):
		move(*at-1, "up")
	case root.Shortcut(0, ui.KeyPageDown):
		move(*at+visible, "pagedown")
	case root.Shortcut(0, ui.KeyPageUp):
		move(*at-visible, "pageup")
	case root.Shortcut(0, ui.KeyHome):
		move(0, "home")
	case root.Shortcut(0, ui.KeyEnd):
		move(n-visible, "end")
	}
	return root
}

// mapRuns cuts highlighted lines into colored non-blank runs.
func mapRuns(c *ui.Context, lang string, lines []string) []mapRun {
	h := newHighlighter(c, lang)
	state := uint64(0)
	var out []mapRun
	for i, line := range lines {
		var spans []ui.Span
		spans, state = h.line(state, line)
		col := 0
		for _, s := range spans {
			for _, word := range strings.SplitAfter(s.Text, " ") {
				text := strings.TrimRight(word, " ")
				if k := utf8.RuneCountInString(text); k > 0 && col < mapCols {
					out = append(out, mapRun{line: i, col: col, n: k, ink: s.Color})
				}
				col += utf8.RuneCountInString(word)
			}
		}
	}
	return out
}

// clampAt keeps a viewport start inside the file.
func clampAt(at, lines, visible int) int {
	return max(0, min(lines-visible, at))
}
