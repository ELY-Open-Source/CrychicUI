// Package home holds smart home components.
package home

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
	t := c.Theme()
	root := ui.Column(c).FillWidth().Radius(t.Radius*1.5).Border(1, t.Border).Background(t.Background).
		Shadow(0, 1, 3, 0, t.Text.Alpha(0.06))
	return kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
}

// tile is a tinted icon square; lit fills it.
func tile(c *ui.Context, mark string, ink ui.Color, size float32, lit bool) ui.Element {
	t := c.Theme()
	box := ui.Box(c).Size(size, size).Shrink(0).Center().Radius(t.Radius)
	box.Background(box.AnimateColor("bg", t.SurfaceHover.Mix(ink, 0.14*unit(lit)), ui.Snappy))
	box.TextColor(box.AnimateColor("ink", t.TextMuted.Mix(ink, unit(lit)), ui.Snappy))
	return box.Children(func() {
		c.Key(mark)
		kit.Arrive(glyph(c, mark, size*0.5), ui.Pose{Scale: 0.6, Fade: 1}, 0)
	})
}

// meter is a bar filled to share in ink.
func meter(c *ui.Context, share float32, ink ui.Color) ui.Element {
	t := c.Theme()
	bar := ui.Box(c).FillWidth().Height(t.Space(1.5)).Radius(999).Background(t.SurfaceHover)
	grow := bar.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	w := grow * bar.AnimateTo("share", max(0, min(1, share)), ui.Smooth)
	col := bar.AnimateColor("ink", ink, ui.Snappy)
	return bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
		p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * w, H: r.H}, col, 999)
	}))
}

// unique panics when items repeat or are blank.
func unique(component string, items []string) {
	seen := map[string]bool{}
	for _, s := range items {
		if strings.TrimSpace(s) == "" || seen[s] {
			panic(fmt.Sprintf("home: %s item %q is blank or repeated", component, s))
		}
		seen[s] = true
	}
}

// number groups a number for the current locale.
func number(c *ui.Context, n any) string {
	return i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), n, i18n.NumberOptions{Integer: true})
}

// toggle is a labeled switch with a sliding knob.
func toggle(c *ui.Context, label string, on bool) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(9), t.Space(5)).Shrink(0).Radius(999).Checked(on))
	b.Background(b.AnimateColor("track", t.SurfaceHover.Mix(t.Border, 0.6*unit(b.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy))
	x := b.AnimateTo("knob", unit(on), ui.Snappy)
	pad := t.Space(0.5)
	return b.Draw(func(p *ui.Painter, r ui.Rect) {
		d := r.H - 2*pad
		p.Fill(ui.Rect{X: r.X + pad + x*(r.W-r.H), Y: r.Y + pad, W: d, H: d}, t.Background, d/2)
	})
}

// round is a round icon button in bg.
func round(c *ui.Context, label, mark string, size float32, bg, ink ui.Color) ui.Element {
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(label).Size(size, size).Shrink(0).Center().Radius(999).TextColor(ink))
	b.Background(b.AnimateColor("bg", bg.Mix(t.Text, 0.08*unit(b.Hovered())), ui.Snappy))
	return b.Children(func() {
		c.Key(mark)
		kit.Arrive(glyph(c, mark, size*0.45), ui.Pose{Scale: 0.6, Fade: 1}, 0)
	})
}

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

// countdown is whole seconds left; it wakes each second.
func countdown(c *ui.Context, until time.Time) int64 {
	now := c.Now()
	// Unix seconds: Sub saturates past 292 years.
	secs := until.Unix() - now.Unix()
	if until.Nanosecond() > now.Nanosecond() {
		secs++
	}
	if secs <= 0 {
		return 0
	}
	frac := time.Duration((until.Nanosecond() - now.Nanosecond() + 1e9) % 1e9)
	if frac == 0 {
		frac = time.Second
	}
	c.After(frac + time.Millisecond)
	return secs
}
