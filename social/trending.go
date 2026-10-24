package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Topic is one ranked trending topic.
type Topic struct {
	Tag   string
	Posts int
}

// TrendingOptions are the trending count labels.
type TrendingOptions struct{ Posts string }

// Trending ranks trending topics.
func Trending(c *ui.Context, topics []Topic, o TrendingOptions) ui.Element {
	if len(topics) == 0 {
		panic("social: Trending needs topics")
	}
	def(&o.Posts, "%d posts")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleList)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, topic := range topics {
			if topic.Tag == "" || topic.Posts < 0 {
				panic("social: invalid topic")
			}
			c.Key(topic.Tag)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				kit.Numeric(c, fmt.Sprint(i+1)).TextColor(t.Accent)
				ui.Text(c, topic.Tag).Grow(1)
				kit.Numeric(c, fmt.Sprintf(o.Posts, topic.Posts)).TextColor(t.TextMuted)
			})
		}
	})
	return root
}
