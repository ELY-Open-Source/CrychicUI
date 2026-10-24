// Package dashboard holds analytics components.
package dashboard

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/i18n"
	"github.com/egoist/mygo/ui"
)

// Point is one value in a series.
type Point struct {
	Label string
	Value float64
}

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

func finite(v float64) bool { return !math.IsNaN(v) && !math.IsInf(v, 0) }

func progress(c *ui.Context, value float64, color ui.Color) ui.Element {
	if value < 0 || value > 1 || !finite(value) {
		panic("dashboard: progress must be in [0,1]")
	}
	t := c.Theme()
	bar := ui.Box(c).FillWidth().Height(t.Space(1.5)).Radius(999).Clip().Background(t.SurfaceHover)
	w := bar.AnimateTo("progress", float32(value), ui.Smooth)
	tint := bar.AnimateColor("ink", color, ui.Snappy)
	bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
		if w > 0 {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(r.H, r.W*w), H: r.H}, tint, 999)
		}
	}))
	return bar
}

// series draws an animated line over a soft area.
func series(c *ui.Context, who string, points []Point, ink ui.Color) ui.Element {
	if len(points) < 2 {
		panic(fmt.Sprintf("dashboard: %s needs two points", who))
	}
	lo, hi := points[0].Value, points[0].Value
	for _, point := range points {
		if !finite(point.Value) {
			panic(fmt.Sprintf("dashboard: %s value at %q must be finite", who, point.Label))
		}
		lo, hi = min(lo, point.Value), max(hi, point.Value)
	}
	span := hi - lo
	if !finite(span) {
		panic(fmt.Sprintf("dashboard: %s range overflows", who))
	}
	t := c.Theme()
	box := ui.Box(c).FillWidth().Height(t.Space(16))
	draw := box.Keyframes("crychic.draw", []ui.Keyframe{{Value: 0}, {At: 800 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	ys := make([]float32, len(points))
	for i, point := range points {
		y := float32(0.5)
		if span > 0 {
			y = float32((point.Value - lo) / span)
		}
		ys[i] = box.AnimateTo(fmt.Sprint("y", i), y, ui.Smooth)
	}
	tint := box.AnimateColor("ink", ink, ui.Snappy)
	box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		pad := t.Space(1)
		at := func(i int) (float32, float32) {
			return r.X + pad + (r.W-2*pad)*float32(i)/float32(len(ys)-1), r.Y + pad + (r.H-2*pad)*(1-ys[i])
		}
		var line, area ui.Path
		x0, _ := at(0)
		xn, yn := at(len(ys) - 1)
		area.MoveTo(x0, r.Y+r.H)
		for i := range ys {
			x, y := at(i)
			if i == 0 {
				line.MoveTo(x, y)
			} else {
				line.LineTo(x, y)
			}
			area.LineTo(x, y)
		}
		area.LineTo(xn, r.Y+r.H).Close()
		p.FillPath(&area, tint.Alpha(0.12*draw))
		p.StrokePath(line.Trim(0, draw), t.Space(0.5), tint)
		if draw >= 1 {
			d := t.Space(2)
			p.Fill(ui.Rect{X: xn - d/2, Y: yn - d/2, W: d, H: d}, tint, 999)
		}
	}))
	return box
}

// statusInk colors a health state.
func statusInk(t *ui.Theme, s Status) ui.Color {
	return [...]ui.Color{t.Success, t.Warning, t.Danger}[s]
}

// statusGlyph is a health state's icon.
func statusGlyph(s Status) string {
	return [...]string{"circle-check", "triangle-alert", "circle-x"}[s]
}

// deltaChip tints a signed change by its leading sign.
func deltaChip(c *ui.Context, text string) ui.Element {
	t := c.Theme()
	ink, glyph := t.TextMuted, "minus"
	switch {
	case strings.HasPrefix(text, "+"):
		ink, glyph = t.Success, "trending-up"
	case falling(text):
		ink, glyph = t.Danger, "trending-down"
	}
	pill := ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999)
	pill.Background(pill.AnimateColor("bg", ink.Alpha(0.12), ui.Snappy)).TextColor(pill.AnimateColor("ink", ink, ui.Snappy))
	return pill.Children(func() {
		ui.Icon(c, icon.Must(glyph)).Size(12, 12).ContentTransition(ui.ContentTransition{})
		kit.Numeric(c, text).FontSize(t.Rem(0.75)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{Numeric: true})
	})
}

// falling reports a change written with a minus sign.
func falling(text string) bool { return strings.HasPrefix(text, "-") || strings.HasPrefix(text, "−") }

// palette is the i-th category color from the theme.
func palette(t *ui.Theme, i int) ui.Color {
	inks := [...]ui.Color{t.Accent, t.Success, t.Warning, t.Danger, t.Accent.Mix(t.Success, 0.5), t.Warning.Mix(t.Danger, 0.5)}
	return inks[(i%len(inks)+len(inks))%len(inks)]
}

// number formats an int or float with the locale's grouping.
func number(c *ui.Context, v any) string {
	return i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), v, i18n.NumberOptions{Integer: true})
}

// tile is a rounded square holding a tinted icon.
func tile(c *ui.Context, glyph string, ink ui.Color) ui.Element {
	t := c.Theme()
	return ui.Box(c).Size(t.Space(9), t.Space(9)).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).
		TextColor(ink).Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16) })
}

// unique panics unless names are set and distinct.
func unique(component string, names []string) {
	seen := make(map[string]bool, len(names))
	for _, n := range names {
		if n == "" || seen[n] {
			panic(fmt.Sprintf("dashboard: %s name %q must be unique and not empty", component, n))
		}
		seen[n] = true
	}
}

// card is the shared widget surface that arrives once.
func card(c *ui.Context) ui.Element {
	t := c.Theme()
	return kit.Arrive(kit.Card(c).FillWidth().Role(ui.RoleGroup), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
}
