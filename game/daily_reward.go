package game

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Reward is one day's prize.
type Reward struct{ Icon, Amount string }

// DailyRewardOptions word the calendar; Day takes a number.
type DailyRewardOptions struct{ Title, Day, Claim, Done, Streak string }

// DailyReward shows days of rewards; claiming today's plays Success.
func DailyReward(c *ui.Context, rewards []Reward, claimed *int, o DailyRewardOptions) ui.Element {
	def(&o.Title, "Daily rewards")
	if claimed == nil || len(rewards) == 0 || *claimed < 0 || *claimed > len(rewards) {
		panic(fmt.Sprintf("game: DailyReward %q needs rewards and 0 <= claimed <= days", o.Title))
	}
	for i, r := range rewards {
		if r.Icon == "" || r.Amount == "" {
			panic(fmt.Sprintf("game: DailyReward %q day %d needs an icon and amount", o.Title, i+1))
		}
	}
	def(&o.Day, "Day %d")
	def(&o.Claim, "Claim")
	def(&o.Done, "All claimed")
	def(&o.Streak, "%d-day streak")
	t := c.Theme()
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "gift", o.Title, func() {
			if *claimed > 0 {
				kit.Chip(c, fmt.Sprintf(o.Streak, *claimed), t.Warning)
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			grid := ui.Grid(c).Columns(min(7, len(rewards))).FillWidth().Gap(t.Space(1.5))
			grid.Children(func() {
				for i, r := range rewards {
					got, today := i < *claimed, i == *claimed
					c.Key(i)
					cell := ui.Column(c).MinWidth(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(1)).Radius(t.Radius)
					cell.Background(cell.AnimateColor("bg", t.SurfaceHover.Alpha(0.6).Mix(t.Warning.Alpha(0.14), unit(today)).
						Mix(t.Success.Alpha(0.1), unit(got)), ui.Snappy))
					cell.Border(1, cell.AnimateColor("line", t.Border.Alpha(0).Mix(t.Warning, unit(today)), ui.Snappy))
					kit.Stagger(c, grid, cell, i).Children(func() {
						kit.Muted(c, fmt.Sprintf(o.Day, i+1)).FontSize(t.Rem(0.6)).Bold().FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
						ui.Box(c).Size(22, 22).Center().Children(func() {
							c.Key(got)
							if got {
								kit.Arrive(glyph(c, "circle-check", 22).TextColor(t.Success), ui.Pose{Scale: 0.4, Fade: 1}, 0)
							} else {
								glyph(c, r.Icon, 22).TextColor(t.Warning.Mix(t.TextMuted, unit(!today)))
							}
						})
						ui.Text(c, r.Amount).FontSize(t.Rem(0.66)).Bold().FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
					})
				}
			})
			all := *claimed == len(rewards)
			label := o.Claim
			if all {
				label = o.Done
			}
			b := kit.Press(ui.ButtonBase(c).Label(label).FillWidth().Padding(t.Space(2)).Radius(t.Radius).Disabled(all))
			b.Background(b.AnimateColor("bg", t.Warning.Mix(t.Warning.Mix(t.Text, 0.15), unit(b.Hovered())), ui.Snappy)).TextColor(t.Background)
			b.Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					glyph(c, "gift", 16)
					ui.Text(c, label).FontSize(t.Rem(0.85)).Bold().Shrink(0)
				})
			})
			if b.Clicked() && !all {
				*claimed++
				root.SetChanged()
				sound.Play(sound.Success)
				kit.Log("DailyReward", "claimed", "day", *claimed)
			}
		})
	})
	return root
}
