// Package health holds health and fitness components.
package health

import (
	"fmt"
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

// tile is a tinted icon square.
func tile(c *ui.Context, mark string, ink ui.Color, size float32) ui.Element {
	t := c.Theme()
	box := ui.Box(c).Size(size, size).Shrink(0).Center().Radius(t.Radius).Background(ink.Alpha(0.12)).TextColor(ink)
	return box.Children(func() { glyph(c, mark, size*0.5) })
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
			panic(fmt.Sprintf("health: %s item %q is blank or repeated", component, s))
		}
		seen[s] = true
	}
}

// number groups an integer for the current locale.
func number(c *ui.Context, n any) string {
	return i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), n, i18n.NumberOptions{Integer: true})
}

// ring is a circle path starting at the top, clockwise.
func ring(cx, cy, r float32) *ui.Path {
	const k = 0.5522847498
	var p ui.Path
	p.MoveTo(cx, cy-r)
	p.CubeTo(cx+r*k, cy-r, cx+r, cy-r*k, cx+r, cy)
	p.CubeTo(cx+r, cy+r*k, cx+r*k, cy+r, cx, cy+r)
	p.CubeTo(cx-r*k, cy+r, cx-r, cy+r*k, cx-r, cy)
	p.CubeTo(cx-r, cy-r*k, cx-r*k, cy-r, cx, cy-r)
	return p.Close()
}

// dial is a ring track with an arc to share.
func dial(c *ui.Context, size, share float32, ink ui.Color, inside func()) ui.Element {
	t := c.Theme()
	d := ui.Box(c).Size(size, size).Shrink(0).Center()
	grow := d.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 700 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	arc := grow * d.AnimateTo("share", max(0, min(1, share)), ui.Smooth)
	col := d.AnimateColor("ink", ink, ui.Snappy)
	d.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		w := size * 0.1
		cx, cy, rr := r.X+r.W/2, r.Y+r.H/2, r.W/2-w/2
		p.StrokePath(ring(cx, cy, rr), w, t.SurfaceHover)
		if arc > 0 {
			p.StrokePath(ring(cx, cy, rr).Trim(0, arc), w, col)
		}
	}))
	return d.Children(inside)
}

// trend draws values from lo to hi as a line.
func trend(c *ui.Context, values []float64, lo, hi float64, ink ui.Color, height float32, guide float64) ui.Element {
	t := c.Theme()
	box := ui.Box(c).FillWidth().Height(height)
	k := box.Keyframes("draw", []ui.Keyframe{{Value: 0}, {At: 800 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
	col := box.AnimateColor("ink", ink, ui.Snappy)
	return box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
		pad := t.Space(1)
		y := func(v float64) float32 { return r.Y + pad + (r.H-2*pad)*float32((hi-v)/(hi-lo)) }
		x := func(i int) float32 { return r.X + pad + (r.W-2*pad)*float32(i)/float32(len(values)-1) }
		if guide >= lo && guide <= hi {
			gy := y(guide)
			for gx := r.X; gx < r.X+r.W; gx += t.Space(2) {
				p.Line(gx, gy, min(gx+t.Space(1), r.X+r.W), gy, 1, t.TextMuted.Alpha(0.6))
			}
		}
		var line, area ui.Path
		line.MoveTo(x(0), y(values[0]))
		area.MoveTo(x(0), r.Y+r.H).LineTo(x(0), y(values[0]))
		for i := 1; i < len(values); i++ {
			line.LineTo(x(i), y(values[i]))
			area.LineTo(x(i), y(values[i]))
		}
		area.LineTo(x(len(values)-1), r.Y+r.H)
		p.Clip(ui.Rect{X: r.X, Y: r.Y, W: r.W * k, H: r.H}, 0, func() { p.FillPath(area.Close(), col.Alpha(0.12)) })
		p.StrokePath(line.Trim(0, k), t.Space(0.5), col)
		p.FillPath(ring(x(len(values)-1), y(values[len(values)-1]), t.Space(1)*k), col)
	}))
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

// clock formats seconds as m:ss, or h:mm:ss past an hour.
func clock(secs int64) string {
	if secs >= 3600 {
		return fmt.Sprintf("%d:%02d:%02d", secs/3600, secs/60%60, secs%60)
	}
	return fmt.Sprintf("%d:%02d", secs/60, secs%60)
}
