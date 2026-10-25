package social

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ProfileStatsOptions are the profile count labels.
type ProfileStatsOptions struct{ Posts, Followers, Following string }

// ProfileStats rolls post, follower, and following counts.
func ProfileStats(c *ui.Context, posts, followers, following int, o ProfileStatsOptions) ui.Element {
	if posts < 0 || followers < 0 || following < 0 {
		panic("social: ProfileStats counts must not be negative")
	}
	def(&o.Posts, "Posts")
	def(&o.Followers, "Followers")
	def(&o.Following, "Following")
	t := c.Theme()
	root := card(c).Row().Role(ui.RoleGroup)
	root.Children(func() {
		stats := []struct {
			label string
			value int
		}{{o.Posts, posts}, {o.Followers, followers}, {o.Following, following}}
		for i, stat := range stats {
			c.Key(i)
			cell := ui.Column(c).Grow(1).Basis(0).MinWidth(0).AlignItems(ui.Center).Gap(t.Space(0.5)).Padding(t.Space(3), t.Space(2))
			if i > 0 {
				cell.BorderWidth(0, 0, 0, 1).BorderColor(t.Border)
			}
			kit.Arrive(cell, ui.Pose{Y: t.Space(1), Fade: 1}, time.Duration(i)*60*time.Millisecond).Children(func() {
				kit.Numeric(c, number(c, stat.value)).FontSize(t.Rem(1.2)).Bold().FillWidth().Wrap().TextAlign(ui.Center)
				kit.Muted(c, stat.label).FontSize(t.Rem(0.75)).MaxWidth(t.Space(30)).SingleLine()
			})
		}
	})
	return root
}
