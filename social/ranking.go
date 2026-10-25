package social

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RankedUser is one weekly leaderboard row.
type RankedUser struct {
	User  User
	Score int
}

// RankingOptions word the card; Score formats the grouped score.
type RankingOptions struct{ Title, Score string }

// Ranking shows a weekly leaderboard.
func Ranking(c *ui.Context, users []RankedUser, o RankingOptions) ui.Element {
	if len(users) == 0 {
		panic("social: Ranking needs users")
	}
	ids := make([]string, len(users))
	for i, row := range users {
		validUser("Ranking", row.User)
		if row.Score < 0 || i > 0 && row.Score > users[i-1].Score {
			panic(fmt.Sprintf("social: Ranking %q needs a score of at most the row above, not negative", row.User.Name))
		}
		ids[i] = row.User.ID
	}
	unique("Ranking", ids)
	def(&o.Title, "This week")
	def(&o.Score, "%s pts")
	t := c.Theme()
	medals := []ui.Color{t.Warning, t.TextMuted, t.Warning.Mix(t.Danger, 0.45)}
	top := float32(max(users[0].Score, 1))
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "trophy", o.Title, nil)
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, row := range users {
				c.Key(row.User.ID)
				line := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius)
				if i == 0 {
					line.Background(t.Warning.Alpha(0.08))
				}
				kit.Stagger(c, list, line, i).Children(func() {
					rank := ui.Box(c).Size(t.Space(6), t.Space(6)).Shrink(0).Center().Radius(999)
					if i < len(medals) {
						rank.Background(medals[i].Alpha(0.18)).TextColor(medals[i])
					} else {
						rank.TextColor(t.TextMuted)
					}
					rank.Children(func() { ui.Text(c, fmt.Sprint(i+1)).FontSize(t.Rem(0.75)).Bold().Shrink(0) })
					avatar(c, row.User, t.Space(9))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
						ui.Text(c, row.User.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
						bar := ui.Box(c).FillWidth().Height(t.Space(1)).Radius(999).Background(t.SurfaceHover)
						grow := bar.Keyframes("grow", []ui.Keyframe{{Value: 0}, {At: 600 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
						fill := grow * bar.AnimateTo("share", float32(row.Score)/top, ui.Smooth)
						ink := t.Accent
						if i == 0 {
							ink = t.Warning
						}
						bar.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
							p.Fill(ui.Rect{X: r.X, Y: r.Y, W: r.W * fill, H: r.H}, ink, 999)
						}))
					})
					kit.Numeric(c, fmt.Sprintf(o.Score, number(c, row.Score))).FontSize(t.Rem(0.8)).Bold().Shrink(0)
				})
			}
		})
	})
	return root
}
