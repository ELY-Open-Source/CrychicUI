// Package learn holds learning and study components.
package learn

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

// dot is a small drawn separator between meta texts.
func dot(c *ui.Context) {
	t := c.Theme()
	ui.Box(c).Size(t.Space(0.75), t.Space(0.75)).Shrink(0).Radius(999).Background(t.TextMuted)
}

// unique panics when items repeat or are blank.
func unique(component string, items []string) {
	seen := map[string]bool{}
	for _, s := range items {
		if strings.TrimSpace(s) == "" || seen[s] {
			panic(fmt.Sprintf("learn: %s item %q is blank or repeated", component, s))
		}
		seen[s] = true
	}
}

// number groups an integer for the current locale.
func number(c *ui.Context, n any) string {
	return i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), n, i18n.NumberOptions{Integer: true})
}

// percent rounds done of total; never 100 early, 0 late.
func percent(done, total int) int {
	p := int(math.Round(float64(done) / float64(total) * 100))
	if done < total {
		p = min(p, 99)
	}
	if done > 0 {
		p = max(p, 1)
	}
	return p
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

// star is a filled five-point star inside r.
func star(r ui.Rect) *ui.Path {
	var p ui.Path
	cx, cy := r.X+r.W/2, r.Y+r.H*0.54
	outer, inner := min(r.W, r.H)/2, min(r.W, r.H)/5
	for i := range 10 {
		rad := outer
		if i%2 == 1 {
			rad = inner
		}
		th := float64(i) * math.Pi / 5
		x, y := cx+rad*float32(math.Sin(th)), cy-rad*float32(math.Cos(th))
		if i == 0 {
			p.MoveTo(x, y)
		} else {
			p.LineTo(x, y)
		}
	}
	return p.Close()
}

// rating draws five stars filled to score of 5.
func rating(c *ui.Context, score float64, size float32) {
	t := c.Theme()
	ui.Row(c).Gap(size * 0.15).Shrink(0).Children(func() {
		for i := range 5 {
			fill := float32(max(0, min(1, score-float64(i))))
			ui.Box(c).Size(size, size).Shrink(0).Draw(func(p *ui.Painter, r ui.Rect) {
				p.FillPath(star(r), t.TextMuted.Alpha(0.25))
				if fill > 0 {
					p.Clip(ui.Rect{X: r.X, Y: r.Y, W: r.W * fill, H: r.H}, 0, func() { p.FillPath(star(r), t.Warning) })
				}
			})
		}
	})
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
