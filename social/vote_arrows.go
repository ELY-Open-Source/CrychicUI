package social

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
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
	root := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		for _, item := range []struct {
			delta int
			label string
			icon  string
		}{{1, o.Up, "chevron-up"}, {-1, o.Down, "chevron-down"}} {
			next := item.delta
			if *vote == item.delta {
				next = 0
			}
			value, ok := votedScore(*score, *vote, next)
			c.Key(item.delta)
			b := button(c, item.label).Checked(*vote == item.delta).Disabled(!ok)
			b.Children(func() { ui.Icon(c, icon.Must(item.icon)).Size(16, 16) })
			if b.Clicked() && ok {
				*score, *vote = value, next
				root.SetChanged()
				sound.Play(sound.Tick)
				kit.Log("VoteArrows", "voted", "vote", next, "score", value)
			}
		}
		kit.Numeric(c, fmt.Sprint(*score))
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
