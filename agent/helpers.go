// Package agent holds prompt, run and answer components.
package agent

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
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

// clock writes d as m:ss, or h:mm:ss past an hour.
func clock(d time.Duration) string {
	s := int64(d / time.Second)
	if s >= 3600 {
		return fmt.Sprintf("%d:%02d:%02d", s/3600, s/60%60, s%60)
	}
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

// ends gives back and forth glyphs and keys.
func ends(c *ui.Context) (back, forth string) {
	back, forth = "chevron-left", "chevron-right"
	if c.Direction() == ui.RightToLeft {
		back, forth = forth, back
	}
	return back, forth
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

// prefixLines returns s through the end of its nth line.
func prefixLines(s string, n int) string {
	if n <= 0 {
		return ""
	}
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			n--
			if n == 0 {
				return s[:i+1]
			}
		}
	}
	return s
}
