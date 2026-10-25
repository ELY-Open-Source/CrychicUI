package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RankedUser is one weekly leaderboard row.
type RankedUser struct {
	User  User
	Score int
}

// RankingOptions are the leaderboard score labels.
type RankingOptions struct{ Score string }

// Ranking shows a weekly leaderboard.
func Ranking(c *ui.Context, users []RankedUser, o RankingOptions) ui.Element {
	if len(users) == 0 {
		panic("social: Ranking needs users")
	}
	def(&o.Score, "%d pts")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleList)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, row := range users {
			validUser("Ranking", row.User)
			if row.Score < 0 {
				panic("social: ranking score must not be negative")
			}
			c.Key(row.User.ID)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				kit.Numeric(c, fmt.Sprint(i+1)).TextColor(t.Accent)
				avatar(c, row.User, t.Space(7))
				ui.Text(c, row.User.Name).Grow(1)
				kit.Numeric(c, fmt.Sprintf(o.Score, row.Score)).TextColor(t.TextMuted)
			})
		}
	})
	return root
}
