package game

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HandOptions word the cards; Label takes rank and suit.
type HandOptions struct {
	Label string
	Suits [4]string
}

// Hand fans cards; hover lifts; clicks pick one; cue Swipe.
func Hand(c *ui.Context, cards []Card, picked *int, o HandOptions) ui.Element {
	if picked == nil || len(cards) == 0 || len(cards) > 13 || *picked < -1 || *picked >= len(cards) {
		panic("game: Hand needs 1 to 13 cards and a pick in range")
	}
	names := make([]string, len(cards))
	for i, k := range cards {
		validCard("Hand", k)
		names[i] = fmt.Sprint(k.Rank, "/", k.Suit)
	}
	unique("Hand", names)
	def(&o.Label, "%s of %s")
	for i, name := range [...]string{"spades", "hearts", "diamonds", "clubs"} {
		def(&o.Suits[i], name)
	}
	t := c.Theme()
	w, h := t.Space(18), t.Space(26)
	n := len(cards)
	step := min(t.Space(10), w*0.6)
	mid := float64(n-1) / 2
	// Size the box to every rotated, dropped card.
	lo, hi, bottom := float32(math.Inf(1)), float32(math.Inf(-1)), float32(0)
	for i := range n {
		off := float64(i) - mid
		sin, cos := math.Sincos(off * 6 * math.Pi / 180)
		sin = math.Abs(sin)
		cx := float32(i)*step + w/2
		hw := float32(float64(w)/2*cos + float64(h)/2*sin)
		hh := float32(float64(w)/2*sin + float64(h)/2*cos)
		lo, hi = min(lo, cx-hw), max(hi, cx+hw)
		bottom = max(bottom, t.Space(8)+h/2+float32(off*off)*t.Space(0.6)+hh)
	}
	left := t.Space(2) - lo
	root := ui.Box(c).Width(hi - lo + t.Space(4)).Height(bottom + t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, k := range cards {
			off := float64(i) - mid
			angle := float32(off * 6)
			drop := float32(off*off) * t.Space(0.6)
			on := *picked == i
			c.Key(names[i])
			b := ui.ButtonBase(c).Label(fmt.Sprintf(o.Label, k.Rank, o.Suits[k.Suit])).Absolute().
				Left(left+float32(i)*step).Top(t.Space(8)).Size(w, h).Radius(t.Radius).Checked(on)
			lift := b.AnimateTo("lift", max(unit(b.Hovered())*0.6, unit(on)), ui.Snappy)
			// One transform: Arrive would overwrite the lift.
			deal := b.Keyframes("deal", []ui.Keyframe{{Value: 1}, {At: 300 * time.Millisecond, Value: 0, Ease: ui.EaseOut}},
				ui.KeyframeOptions{Delay: time.Duration(i) * 50 * time.Millisecond})
			b.Translate(0, drop-lift*t.Space(7)+deal*t.Space(10)).Rotate(angle*(1-lift*0.7) - float32(off)*12*deal).Opacity(1 - deal)
			b.Background(t.Background).Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy)).
				Shadow(0, 1+3*lift, 3+6*lift, 0, t.Text.Alpha(0.08+0.06*lift))
			b.Children(func() { face(c, k, w) })
			if b.Clicked() {
				if on {
					*picked = -1
				} else {
					*picked = i
				}
				root.SetChanged()
				sound.Play(sound.Swipe)
				kit.Log("Hand", "picked", "index", *picked)
			}
		}
	})
	return root
}
