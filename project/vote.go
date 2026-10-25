package project

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// VoteOptions word the card; Vote formats each button's label.
type VoteOptions struct{ Title, Vote, Total, One, Top string }

// Vote adds dot votes on cards; cue Pop.
func Vote(c *ui.Context, cards []string, votes *[]int, o VoteOptions) ui.Element {
	def(&o.Title, "Dot vote")
	def(&o.Vote, "Vote for %s")
	def(&o.Total, "%d votes")
	def(&o.One, "1 vote")
	def(&o.Top, "Top")
	if len(cards) == 0 || votes == nil || len(*votes) != len(cards) {
		panic(fmt.Sprintf("project: Vote %q needs one count per card", o.Title))
	}
	unique("Vote", cards)
	total, top := 0, 0
	for i, n := range *votes {
		if n < 0 || total > math.MaxInt-n {
			panic(fmt.Sprintf("project: Vote %q has an invalid count %d for %q", o.Title, n, cards[i]))
		}
		total += n
		top = max(top, n)
	}
	count := fmt.Sprintf(o.Total, total)
	if total == 1 {
		count = o.One
	}
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "vote", o.Title, func() { kit.Chip(c, count, t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(2.5))
		list.Children(func() {
			for i, text := range cards {
				c.Key(text)
				n := (*votes)[i]
				lead := n > 0 && n == top
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius)
				row.Border(1, row.AnimateColor("line", t.Border.Mix(t.Accent, unit(lead)), ui.Snappy))
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.Accent.Alpha(0.06), unit(lead)), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, text).FontSize(t.Rem(0.82)).Grow(1).MinWidth(0).Wrap()
							if lead {
								kit.Arrive(kit.Chip(c, o.Top, t.Accent), ui.Pose{Scale: 0.7, Fade: 1}, 0)
							}
						})
						ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
							for d := range min(n, 5) {
								c.Key(d)
								kit.Arrive(ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Radius(999).Background(t.Accent), ui.Pose{Scale: 0.3, Fade: 1}, 0)
							}
							if n > 5 {
								kit.Muted(c, fmt.Sprintf("+%d", n-5)).FontSize(t.Rem(0.7)).Shrink(0)
							}
						})
					})
					kit.Numeric(c, fmt.Sprint(n)).FontSize(t.Rem(0.95)).Bold().Shrink(0).TextColor(t.Text.Mix(t.Accent, unit(lead)))
					full := n == math.MaxInt || total == math.MaxInt
					b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Vote, text)).Size(t.Space(8), t.Space(8)).Shrink(0).Center().
						Radius(999).Disabled(full))
					b.Background(b.AnimateColor("bg", t.Accent.Alpha(0.12).Mix(t.Accent.Alpha(0.22), unit(b.Hovered())), ui.Snappy))
					b.TextColor(t.Accent).Children(func() { glyph(c, "plus", 16) })
					if b.Clicked() && !full {
						(*votes)[i]++
						root.SetChanged()
						sound.Play(sound.Pop)
						kit.Log("Vote", "voted", "index", i, "count", (*votes)[i])
					}
				})
			}
		})
	})
	return root
}
