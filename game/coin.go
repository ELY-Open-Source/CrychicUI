package game

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// CoinOptions word the coin; Count takes side and tally.
type CoinOptions struct {
	Heads, Tails, Flip, Count string
	Rand                      func(n int) int
}

// Coin flips a coin edge over edge; cue Pop.
func Coin(c *ui.Context, heads *bool, o CoinOptions) ui.Element {
	def(&o.Flip, "Flip")
	if heads == nil {
		panic(fmt.Sprintf("game: Coin %q needs a side", o.Flip))
	}
	def(&o.Heads, "Heads")
	def(&o.Tails, "Tails")
	def(&o.Count, "%s %d")
	if o.Rand == nil {
		o.Rand = pick
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5), t.Space(4)).AlignItems(ui.Center).Role(ui.RoleGroup)
	flips := ui.Local(root, "flips", func() int { return 0 })
	tally := ui.Local(root, "tally", func() [2]int { return [2]int{} })
	root.Children(func() {
		c.Key(*flips)
		size := t.Space(24)
		coin := ui.Box(c).Size(size, size).Shrink(0).Center()
		spin := float32(1)
		if *flips > 0 {
			spin = coin.Keyframes("spin", []ui.Keyframe{{Value: 0}, {At: 900 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
		}
		// Six half turns land on the result face.
		turn := float64(spin) * 6 * math.Pi
		up := *heads
		if int(turn/math.Pi)%2 == 1 && spin < 1 {
			up = !up
		}
		width := float32(math.Max(0.04, math.Abs(math.Cos(turn))))
		lift := float32(math.Sin(float64(spin) * math.Pi))
		coin.Scale(width, 1).Translate(0, -t.Space(8)*lift)
		coin.Draw(func(p *ui.Painter, r ui.Rect) {
			cx, cy, rr := r.X+r.W/2, r.Y+r.H/2, r.W/2
			p.FillPath(ring(cx, cy, rr), t.Warning)
			p.FillPath(ring(cx, cy, rr*0.84), t.Warning.Mix(t.Background, 0.25))
		})
		coin.Children(func() {
			mark := "crown"
			if !up {
				mark = "star"
			}
			glyph(c, mark, size*0.42).TextColor(t.Warning.Mix(t.Text, 0.45))
		})
		side := o.Tails
		if *heads {
			side = o.Heads
		}
		c.Key("side")
		res := ui.Text(c, side).FontSize(t.Rem(1.1)).Bold().Shrink(0)
		res.Opacity(res.AnimateTo("shown", unit(spin >= 1), ui.Smooth))
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			kit.Chip(c, fmt.Sprintf(o.Count, o.Heads, tally[0]), t.Warning)
			kit.Chip(c, fmt.Sprintf(o.Count, o.Tails, tally[1]), t.TextMuted)
		})
		b := kit.Press(ui.ButtonBase(c).Label(o.Flip).Padding(t.Space(2), t.Space(5)).Radius(999).Shrink(0))
		b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(b.Hovered())), ui.Snappy)).TextColor(t.AccentText)
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, "coins", 16)
				ui.Text(c, o.Flip).FontSize(t.Rem(0.85)).Bold().Shrink(0)
			})
		})
		if b.Clicked() {
			n := o.Rand(2)
			if n < 0 || n > 1 {
				panic(fmt.Sprintf("game: Coin %q Rand gave %d for 2 sides", o.Flip, n))
			}
			*heads = n == 0
			tally[n]++
			*flips++
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("Coin", "flipped", "heads", *heads)
		}
	})
	return root
}
