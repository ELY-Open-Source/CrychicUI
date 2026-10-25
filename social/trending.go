package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Topic is one ranked trending topic; Category is optional.
type Topic struct {
	Tag, Category string
	Posts         int
}

// TrendingOptions word the card; Posts formats the grouped count.
type TrendingOptions struct{ Title, Posts string }

// Trending ranks trending topics.
func Trending(c *ui.Context, topics []Topic, o TrendingOptions) ui.Element {
	if len(topics) == 0 {
		panic("social: Trending needs topics")
	}
	tags := make([]string, len(topics))
	for i, topic := range topics {
		if topic.Posts < 0 {
			panic(fmt.Sprintf("social: Trending topic %q has negative posts", topic.Tag))
		}
		tags[i] = topic.Tag
	}
	unique("Trending", tags)
	def(&o.Title, "Trending")
	def(&o.Posts, "%s posts")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "trending-up", o.Title, nil)
		list := ui.Column(c).FillWidth().Padding(t.Space(1.5))
		list.Children(func() {
			for i, topic := range topics {
				c.Key(topic.Tag)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					rank := kit.Numeric(c, fmt.Sprint(i+1)).FontSize(t.Rem(1.05)).Bold().Width(t.Space(5)).Shrink(0)
					rank.TextColor(t.TextMuted)
					if i < 3 {
						rank.TextColor(t.Accent)
					}
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
						if topic.Category != "" {
							kit.Muted(c, topic.Category).FontSize(t.Rem(0.7)).FillWidth().MinWidth(0).SingleLine()
						}
						ui.Text(c, topic.Tag).FontSize(t.Rem(0.88)).Bold().FillWidth().MinWidth(0).SingleLine()
						kit.Numeric(c, fmt.Sprintf(o.Posts, number(c, topic.Posts))).FontSize(t.Rem(0.72)).TextColor(t.TextMuted)
					})
					if i == 0 {
						glyph(c, "flame", 16).TextColor(t.Danger)
					}
				})
			}
		})
	})
	return root
}
