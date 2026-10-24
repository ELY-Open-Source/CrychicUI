package social

import (
	"fmt"

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
	root := ui.Row(c).FillWidth().Gap(t.Space(4)).Justify(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		stats := []struct {
			label string
			value int
		}{{o.Posts, posts}, {o.Followers, followers}, {o.Following, following}}
		for i, stat := range stats {
			c.Key(i)
			ui.Column(c).AlignItems(ui.Center).Children(func() {
				kit.Numeric(c, fmt.Sprint(stat.value)).Bold()
				kit.Muted(c, stat.label).FontSize(t.Rem(0.75))
			})
		}
	})
	return root
}
