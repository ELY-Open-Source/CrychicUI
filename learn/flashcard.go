package learn

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FlashcardOptions word the card's two sides and its hint.
type FlashcardOptions struct{ Question, Answer, Hint string }

// Flashcard flips between question and answer; cue Swipe.
func Flashcard(c *ui.Context, front, back string, flipped *bool, o FlashcardOptions) ui.Element {
	if flipped == nil || strings.TrimSpace(front) == "" || strings.TrimSpace(back) == "" {
		panic(fmt.Sprintf("learn: Flashcard %q needs a front, a back and a side", front))
	}
	def(&o.Question, "Question")
	def(&o.Answer, "Answer")
	def(&o.Hint, "Click to flip")
	t := c.Theme()
	text := front
	if *flipped {
		text = back
	}
	b := ui.ButtonBase(c).Label(text).FillWidth().MinHeight(t.Space(48)).Radius(t.Radius * 1.5).Checked(*flipped)
	flips := ui.Local(b, "flips", func() int { return 0 })
	flipAt := ui.Local(b, "at", func() time.Time { return time.Time{} })
	width := b.Keyframes("flip", []ui.Keyframe{{Value: 1}, {At: 180 * time.Millisecond, Value: 0, Ease: ui.EaseIn},
		{At: 360 * time.Millisecond, Value: 1, Ease: ui.EaseOut}}, ui.KeyframeOptions{Run: uint64(*flips)})
	if *flips == 0 {
		width = 1
	}
	// The shown side turns over once the card is edge-on.
	shown := ui.Local(b, "shown", func() bool { return *flipped })
	if *shown != *flipped {
		gone := c.Now().Sub(*flipAt)
		if gone >= 180*time.Millisecond || c.Preferences().ReduceMotion {
			*shown = *flipped
		} else {
			c.After(180*time.Millisecond - gone)
		}
	}
	hover := b.AnimateTo("hover", unit(b.Hovered()), ui.Snappy)
	ink := t.Accent
	if *shown {
		ink = t.Success
	}
	b.Scale(max(0.03, width), 1).Translate(0, -t.Space(1)*hover)
	b.Background(b.AnimateColor("bg", t.Background.Mix(ink.Alpha(0.06), unit(*shown)), ui.Snappy)).
		Border(1, b.AnimateColor("line", t.Border.Mix(ink, 0.5*unit(*shown)+0.3*hover), ui.Snappy)).
		Shadow(0, 1+3*hover, 3+6*hover, 0, t.Text.Alpha(0.06+0.05*hover))
	b.Children(func() {
		c.Key(*shown)
		side, body, mark := o.Question, front, "circle-help"
		if *shown {
			side, body, mark = o.Answer, back, "lightbulb"
		}
		ui.Column(c).Fill().Padding(t.Space(4)).Gap(t.Space(3)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(ink).Children(func() {
				glyph(c, mark, 14)
				ui.Text(c, side).FontSize(t.Rem(0.72)).Bold().MinWidth(0).SingleLine()
			})
			ui.Box(c).FillWidth().Grow(1).Center().Children(func() {
				ui.Text(c, body).FontSize(t.Rem(1.15)).Bold().TextAlign(ui.Center).FillWidth().Wrap()
			})
			ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Justify(ui.Center).TextColor(t.TextMuted).Children(func() {
				glyph(c, "repeat", 12)
				ui.Text(c, o.Hint).FontSize(t.Rem(0.68)).MinWidth(0).SingleLine()
			})
		})
	})
	if b.Clicked() {
		*flipped = !*flipped
		*flips++
		*flipAt = c.Now()
		b.SetChanged()
		sound.Play(sound.Swipe)
		kit.Log("Flashcard", "flipped", "answer", *flipped)
	}
	return b
}
