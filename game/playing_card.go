package game

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PlayingCardOptions word the card; Label takes rank and suit.
type PlayingCardOptions struct {
	Label, Down string
	Suits       [4]string
}

// PlayingCard flips over on a click; cue Swipe.
func PlayingCard(c *ui.Context, k Card, up *bool, o PlayingCardOptions) ui.Element {
	if up == nil {
		panic(fmt.Sprintf("game: PlayingCard %q needs a side", k.Rank))
	}
	validCard("PlayingCard", k)
	def(&o.Label, "%s of %s")
	def(&o.Down, "Face-down card")
	for i, name := range [...]string{"spades", "hearts", "diamonds", "clubs"} {
		def(&o.Suits[i], name)
	}
	t := c.Theme()
	w, h := t.Space(28), t.Space(40)
	label := o.Down
	if *up {
		label = fmt.Sprintf(o.Label, k.Rank, o.Suits[k.Suit])
	}
	b := ui.ButtonBase(c).Label(label).Size(w, h).Radius(t.Radius * 1.2).Checked(*up)
	flips := ui.Local(b, "flips", func() int { return 0 })
	flipAt := ui.Local(b, "at", func() time.Time { return time.Time{} })
	// Starting and ending full width, a replay cannot jump.
	width := b.Keyframes("flip", []ui.Keyframe{{Value: 1}, {At: 210 * time.Millisecond, Value: 0, Ease: ui.EaseIn},
		{At: 420 * time.Millisecond, Value: 1, Ease: ui.EaseOut}}, ui.KeyframeOptions{Run: uint64(*flips)})
	if *flips == 0 {
		width = 1
	}
	// The shown face turns over once the card is edge-on.
	shown := ui.Local(b, "shown", func() bool { return *up })
	if *shown != *up {
		gone := c.Now().Sub(*flipAt)
		if gone >= 210*time.Millisecond || c.Preferences().ReduceMotion {
			*shown = *up
		} else {
			c.After(210*time.Millisecond - gone)
		}
	}
	hover := b.AnimateTo("hover", unit(b.Hovered()), ui.Snappy)
	b.Scale(max(0.03, width), 1).Translate(0, -t.Space(1.5)*hover)
	b.Background(t.Background).Border(1, t.Border).Shadow(0, 1+3*hover, 3+6*hover, 0, t.Text.Alpha(0.08+0.06*hover))
	b.Children(func() {
		c.Key(*shown)
		if *shown {
			face(c, k, w)
		} else {
			back(c)
		}
	})
	if b.Clicked() {
		*up = !*up
		*flips++
		*flipAt = c.Now()
		b.SetChanged()
		sound.Play(sound.Swipe)
		kit.Log("PlayingCard", "flipped", "up", *up)
	}
	return b
}
