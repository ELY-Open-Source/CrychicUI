package game

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Player is one ranked player; Prev 0 means new.
type Player struct {
	Name  string
	Score int64
	Prev  int
}

// LeaderboardOptions word the board.
type LeaderboardOptions struct{ Title, New string }

// Leaderboard ranks players by score with movement arrows.
func Leaderboard(c *ui.Context, players []Player, o LeaderboardOptions) ui.Element {
	def(&o.Title, "Leaderboard")
	if len(players) == 0 {
		panic(fmt.Sprintf("game: Leaderboard %q needs players", o.Title))
	}
	names := make([]string, len(players))
	for i, p := range players {
		if p.Prev < 0 || i > 0 && p.Score > players[i-1].Score {
			panic(fmt.Sprintf("game: Leaderboard %q needs scores high to low and a last rank of at least 0", p.Name))
		}
		names[i] = p.Name
	}
	unique("Leaderboard", names)
	def(&o.New, "New")
	t := c.Theme()
	medals := []ui.Color{t.Warning, t.TextMuted, t.Warning.Mix(t.Danger, 0.45)}
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "trophy", o.Title, func() { kit.Chip(c, fmt.Sprint(len(players)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, p := range players {
				rank := i + 1
				c.Key(p.Name)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
					Transition(ui.ElementTransition{Animation: &ui.Snappy, Position: true})
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					badge := ui.Box(c).Size(t.Space(7), t.Space(7)).Shrink(0).Center().Radius(999)
					if i < 3 {
						badge.Background(medals[i].Alpha(0.18)).TextColor(medals[i])
					}
					badge.Children(func() { kit.Numeric(c, fmt.Sprint(rank)).FontSize(t.Rem(0.8)).Bold().Shrink(0) })
					display.Avatar(c, p.Name, nil).Size(t.Space(8), t.Space(8))
					ui.Text(c, p.Name).FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0).SingleLine()
					move(c, p.Prev, rank, o.New)
					kit.Numeric(c, number(c, p.Score)).FontSize(t.Rem(0.85)).Bold().Shrink(0).MaxWidth(t.Space(28)).SingleLine().
						ContentTransition(ui.ContentTransition{Numeric: true})
				})
			}
		})
	})
	return root
}

// move shows how far a player climbed or fell.
func move(c *ui.Context, prev, rank int, fresh string) {
	t := c.Theme()
	ui.Row(c).Width(t.Space(12)).Shrink(0).Gap(t.Space(0.5)).AlignItems(ui.Center).Justify(ui.End).Children(func() {
		c.Key(prev - rank)
		switch {
		case prev == 0:
			kit.Arrive(kit.Chip(c, fresh, t.Accent), ui.Pose{Scale: 0.7, Fade: 1}, 0)
		case prev > rank:
			kit.Arrive(glyph(c, "arrow-up", 14).TextColor(t.Success).Shrink(0), ui.Pose{Y: t.Space(2), Fade: 1}, 0)
			ui.Text(c, fmt.Sprint(prev-rank)).FontSize(t.Rem(0.72)).Bold().TextColor(t.Success).MinWidth(0).SingleLine()
		case prev < rank:
			kit.Arrive(glyph(c, "arrow-down", 14).TextColor(t.Danger).Shrink(0), ui.Pose{Y: -t.Space(2), Fade: 1}, 0)
			ui.Text(c, fmt.Sprint(rank-prev)).FontSize(t.Rem(0.72)).Bold().TextColor(t.Danger).MinWidth(0).SingleLine()
		default:
			ui.Box(c).Size(t.Space(2.5), t.Space(0.5)).Radius(999).Background(t.TextMuted)
		}
	})
}
