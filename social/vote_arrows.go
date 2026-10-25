package social

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VoteArrowsOptions are the vote labels.
type VoteArrowsOptions struct{ Up, Down string }

// VoteArrows changes scores; cue Tick.
func VoteArrows(c *ui.Context, score, vote *int, o VoteArrowsOptions) ui.Element {
	if score == nil || vote == nil || *vote < -1 || *vote > 1 {
		panic("social: VoteArrows needs score and valid vote")
	}
	def(&o.Up, "Upvote")
	def(&o.Down, "Downvote")
	t := c.Theme()
	root := ui.Column(c).Gap(t.Space(0.5)).AlignItems(ui.Center).AlignSelf(ui.Start).Padding(t.Space(1)).Radius(999).Shrink(0).Role(ui.RoleGroup)
	root.Background(root.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.1), unit(*vote == 1)).Mix(t.Danger.Alpha(0.1), unit(*vote == -1)), ui.Snappy))
	arrow := func(delta int, label, mark string, ink ui.Color) {
		next := delta
		if *vote == delta {
			next = 0
		}
		value, ok := votedScore(*score, *vote, next)
		on := *vote == delta
		c.Key(delta)
		b := kit.Press(ui.ButtonBase(c).Label(label).Size(t.Space(8), t.Space(8)).Radius(999).Checked(on).Disabled(!ok))
		b.Background(b.AnimateColor("bg", ink.Alpha(0).Mix(ink.Alpha(0.14), unit(b.Hovered())), ui.Snappy))
		b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(ink, max(unit(on), unit(b.Hovered())*0.7)), ui.Snappy))
		b.Children(func() {
			c.Key(on)
			icon := glyph(c, mark, 20)
			if on {
				kit.Arrive(icon, ui.Pose{Y: float32(delta) * t.Space(2), Scale: 1.3}, 0)
			}
		})
		if b.Clicked() && ok {
			*score, *vote = value, next
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("VoteArrows", "voted", "vote", next, "score", value)
		}
	}
	root.Children(func() {
		arrow(1, o.Up, "arrow-big-up", t.Accent)
		count := kit.Numeric(c, number(c, *score)).FontSize(t.Rem(0.9)).Bold().Shrink(0)
		count.TextColor(count.AnimateColor("ink", t.Text.Mix(t.Accent, unit(*vote == 1)).Mix(t.Danger, unit(*vote == -1)), ui.Snappy))
		arrow(-1, o.Down, "arrow-big-down", t.Danger)
	})
	return root
}

func votedScore(score, current, next int) (int, bool) {
	value, ok := checkedAdd(score, -current)
	if !ok || next == 0 {
		return value, ok
	}
	return checkedAdd(value, next)
}

func checkedAdd(a, b int) (int, bool) {
	if b > 0 && a > math.MaxInt-b || b < 0 && a < math.MinInt-b {
		return a, false
	}
	return a + b, true
}
