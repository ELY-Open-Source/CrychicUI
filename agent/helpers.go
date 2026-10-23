// Package agent holds prompt, run and answer components.
package agent

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
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
		panic(fmt.Sprintf("agent: %s sizes must be finite and not negative", component))
	}
	if *v == 0 {
		*v = d
	}
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

// tag is a small muted pill of text.
func tag(c *ui.Context, text string) ui.Element {
	t := c.Theme()
	return kit.Muted(c, text).FontSize(t.Rem(0.75)).Padding(t.Space(0.5), t.Space(1.5)).Radius(999).Background(t.Surface)
}

// ends gives back and forth glyphs and keys.
func ends(c *ui.Context) (back, forth string) {
	back, forth = "chevron-left", "chevron-right"
	if c.Direction() == ui.RightToLeft {
		back, forth = forth, back
	}
	return back, forth
}
