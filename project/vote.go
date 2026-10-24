package project

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VoteOptions are the vote labels.
type VoteOptions struct{ Vote string }

// Vote adds dot votes on cards; cue Pop.
func Vote(c *ui.Context, cards []string, votes *[]int, o VoteOptions) ui.Element {
	if len(cards) == 0 || votes == nil || len(*votes) != len(cards) {
		panic("project: Vote needs cards and counts")
	}
	def(&o.Vote, "Vote")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	root.Children(func() {
		for i, card := range cards {
			if card == "" || (*votes)[i] < 0 {
				panic("project: Vote card must not be empty")
			}
			c.Key(i)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				ui.Text(c, card).Grow(1)
				kit.Numeric(c, fmt.Sprint((*votes)[i])).TextColor(t.Accent)
				b := kit.Press(ui.ButtonBase(c).Label(o.Vote).Padding(t.Space(1), t.Space(2)).Radius(999).
					Disabled((*votes)[i] == math.MaxInt))
				b.Children(func() { ui.Text(c, o.Vote).SingleLine() })
				if b.Clicked() && (*votes)[i] < math.MaxInt {
					(*votes)[i]++
					root.SetChanged()
					sound.Play(sound.Pop)
					kit.Log("Vote", "voted", "index", i, "count", (*votes)[i])
				}
			})
		}
	})
	return root
}
