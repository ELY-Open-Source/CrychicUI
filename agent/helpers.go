// Package agent holds prompt, run and answer components.
package agent

import (
	"fmt"
	"math"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/display"
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

// StepState is the state of a step or tool call.
type StepState uint8

// Step states.
const (
	StepPending StepState = iota
	StepRunning
	StepDone
	StepFailed
)

// Source is a web page an answer cites.
type Source struct {
	Title, Host, URL, Snippet string
}

func validSource(component string, s Source) {
	if s.Title == "" || s.Host == "" || s.URL == "" {
		panic(fmt.Sprintf("agent: %s source needs a title, a host and a URL", component))
	}
}

// mark draws a state: dashed, spinning, checked or crossed.
func mark(c *ui.Context, s StepState) ui.Element {
	t := c.Theme()
	box := ui.Box(c).Size(t.Space(5), t.Space(5)).Shrink(0).Center()
	return box.Children(func() {
		c.Key(s)
		switch s {
		case StepRunning:
			ic := ui.Icon(c, icon.Must("loader-circle")).Size(16, 16).TextColor(t.Accent)
			ic.Rotate(360 * float32(ic.Loop("spin", time.Second, ui.Linear)))
		case StepDone:
			kit.Arrive(ui.Icon(c, icon.Must("circle-check")).Size(16, 16).TextColor(t.Success), ui.Pose{Scale: 0.4, Fade: 1}, 0)
		case StepFailed:
			kit.Arrive(ui.Icon(c, icon.Must("circle-x")).Size(16, 16).TextColor(t.Danger), ui.Pose{Scale: 0.4, Fade: 1}, 0)
		default:
			ui.Icon(c, icon.Must("circle-dashed")).Size(16, 16).TextColor(t.TextMuted)
		}
	})
}

// stateBadge is a colored pill naming s.
func stateBadge(c *ui.Context, s StepState, pending, running, done, failed string) ui.Element {
	text, v := pending, display.BadgeNeutral
	switch s {
	case StepRunning:
		text, v = running, display.BadgeAccent
	case StepDone:
		text, v = done, display.BadgeSuccess
	case StepFailed:
		text, v = failed, display.BadgeDanger
	}
	c.Key(s)
	return display.Badge(c, text, display.BadgeOptions{Variant: v})
}

// tile is a tinted rounded square holding an icon.
func tile(c *ui.Context, glyph string, ink ui.Color, side float32) ui.Element {
	t := c.Theme()
	return ui.Box(c).Size(side, side).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).TextColor(ink).
		Children(func() { ui.Icon(c, icon.Must(glyph)).Size(16, 16) })
}

// hostTile is a site's first letter in a tinted square.
func hostTile(c *ui.Context, host string) ui.Element {
	t := c.Theme()
	r, _ := utf8.DecodeRuneInString(host)
	return ui.Box(c).Size(t.Space(5), t.Space(5)).Shrink(0).Center().Radius(t.Space(1)).Background(t.Accent.Alpha(0.14)).
		Children(func() {
			ui.Text(c, string(unicode.ToUpper(r))).FontSize(t.Rem(0.7)).Bold().TextColor(t.Accent)
		})
}

// meterInk turns accent, warning, then danger as share fills.
func meterInk(c *ui.Context, share float64) ui.Color {
	t := c.Theme()
	switch {
	case share >= 0.95:
		return t.Danger
	case share >= 0.8:
		return t.Warning
	}
	return t.Accent
}

// bar is a thin track filled to share in ink.
func bar(c *ui.Context, key string, share float64, ink ui.Color) ui.Element {
	t := c.Theme()
	track := ui.Box(c).FillWidth().Height(t.Space(1.5)).Radius(999).Background(t.SurfaceHover)
	w := track.AnimateTo(key, float32(min(max(share, 0), 1)), ui.Smooth)
	ink = track.AnimateColor(key+".ink", ink, ui.Snappy)
	return track.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
		if w > 0 {
			p.Fill(ui.Rect{X: r.X, Y: r.Y, W: max(r.H, r.W*w), H: r.H}, ink, 999)
		}
	}))
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
