// Package game holds game interface components.
package game

import (
	"fmt"
	"math"
	"math/rand/v2"
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
			panic(fmt.Sprintf("game: %s item %q is blank or repeated", component, s))
		}
		seen[s] = true
	}
}

// number groups an integer for the current locale.
func number(c *ui.Context, n any) string {
	return i18n.FormatNumber(c.Region().FormatLocale(c.Locale()), n, i18n.NumberOptions{Integer: true})
}

// pick is the default random choice in [0, n).
func pick(n int) int { return rand.IntN(n) }

// ellipse is a closed ellipse path.
func ellipse(cx, cy, rx, ry float32) *ui.Path {
	const k = 0.5522847498
	var p ui.Path
	p.MoveTo(cx, cy-ry)
	p.CubeTo(cx+rx*k, cy-ry, cx+rx, cy-ry*k, cx+rx, cy)
	p.CubeTo(cx+rx, cy+ry*k, cx+rx*k, cy+ry, cx, cy+ry)
	p.CubeTo(cx-rx*k, cy+ry, cx-rx, cy+ry*k, cx-rx, cy)
	p.CubeTo(cx-rx, cy-ry*k, cx-rx*k, cy-ry, cx, cy-ry)
	return p.Close()
}

// ring is a circle path starting at the top, clockwise.
func ring(cx, cy, r float32) *ui.Path { return ellipse(cx, cy, r, r) }

// heart is a filled heart path inside r.
func heart(r ui.Rect) *ui.Path {
	var p ui.Path
	x, y, w, h := r.X, r.Y, r.W, r.H
	p.MoveTo(x+w/2, y+h*0.92)
	p.CubeTo(x+w*0.1, y+h*0.62, x-w*0.02, y+h*0.3, x+w*0.2, y+h*0.12)
	p.CubeTo(x+w*0.36, y, x+w*0.48, y+h*0.08, x+w/2, y+h*0.22)
	p.CubeTo(x+w*0.52, y+h*0.08, x+w*0.64, y, x+w*0.8, y+h*0.12)
	p.CubeTo(x+w*1.02, y+h*0.3, x+w*0.9, y+h*0.62, x+w/2, y+h*0.92)
	return p.Close()
}

// wedge is a pie slice, radians clockwise from the top.
func wedge(cx, cy, r float32, a, b float64) *ui.Path {
	var p ui.Path
	p.MoveTo(cx, cy)
	steps := max(2, int(math.Ceil((b-a)/(math.Pi/24))))
	for i := 0; i <= steps; i++ {
		th := a + (b-a)*float64(i)/float64(steps)
		p.LineTo(cx+r*float32(math.Sin(th)), cy-r*float32(math.Cos(th)))
	}
	return p.Close()
}

// Rarity grades an item.
type Rarity int

const (
	Common Rarity = iota
	Uncommon
	Rare
	Epic
	Legendary
)

func rarityInk(t *ui.Theme, r Rarity) ui.Color {
	return [...]ui.Color{t.TextMuted, t.Success, t.Accent, t.Accent.Mix(t.Danger, 0.5), t.Warning}[r]
}

func validRarity(component, label string, r Rarity) {
	if r < Common || r > Legendary {
		panic(fmt.Sprintf("game: %s %q rarity %d is not Common to Legendary", component, label, r))
	}
}

// Suit is a playing card suit.
type Suit int

const (
	SuitSpades Suit = iota
	SuitHearts
	SuitDiamonds
	SuitClubs
)

// Card is one playing card; Rank is printed, as "Q".
type Card struct {
	Rank string
	Suit Suit
}

func validCard(component string, k Card) {
	if strings.TrimSpace(k.Rank) == "" || k.Suit < SuitSpades || k.Suit > SuitClubs {
		panic(fmt.Sprintf("game: %s card %q needs a rank and a suit", component, k.Rank))
	}
}

func suitMark(s Suit) string { return [...]string{"spade", "heart", "diamond", "club"}[s] }

func suitInk(t *ui.Theme, s Suit) ui.Color {
	if s == SuitHearts || s == SuitDiamonds {
		return t.Danger
	}
	return t.Text
}

// face is a corner index above a large pip.
func face(c *ui.Context, k Card, w float32) {
	t := c.Theme()
	ink := suitInk(t, k.Suit)
	ui.Column(c).FillWidth().FillHeight().Padding(w * 0.06).AlignItems(ui.Start).TextColor(ink).Children(func() {
		ui.Column(c).AlignItems(ui.Center).Gap(0).Shrink(0).Children(func() {
			ui.Text(c, k.Rank).FontSize(w * 0.2).Bold().Shrink(0)
			glyph(c, suitMark(k.Suit), w*0.14)
		})
		ui.Box(c).Grow(1).FillWidth().Center().Children(func() { glyph(c, suitMark(k.Suit), w*0.38) })
	})
}

// back is a card's patterned back.
func back(c *ui.Context) {
	t := c.Theme()
	ui.Box(c).FillWidth().FillHeight().Padding(t.Space(1)).Children(func() {
		b := ui.Box(c).FillWidth().FillHeight().Radius(t.Radius * 0.6).Background(t.Accent).Center()
		b.Draw(kit.Inside(t.Radius*0.6, func(p *ui.Painter, r ui.Rect) {
			step := t.Space(2)
			for y := r.Y + step/2; y < r.Y+r.H; y += step {
				for x := r.X + step/2; x < r.X+r.W; x += step {
					p.FillPath(ring(x, y, step*0.18), t.AccentText.Alpha(0.35))
				}
			}
		}))
		b.Children(func() { glyph(c, "sparkles", 20).TextColor(t.AccentText) })
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

// stars draws n of three stars, filled ones in ink.
func stars(c *ui.Context, n int, size float32, ink ui.Color) {
	t := c.Theme()
	ui.Row(c).Gap(size * 0.15).Children(func() {
		for i := range 3 {
			b := ui.Box(c).Size(size, size).Shrink(0)
			on := i < n
			b.Draw(func(p *ui.Painter, r ui.Rect) {
				if on {
					p.FillPath(star(r), ink)
				} else {
					p.FillPath(star(r), t.TextMuted.Alpha(0.25))
				}
			})
		}
	})
}
