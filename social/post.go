package social

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PostOptions word the actions; Likes, Comments, Shares format grouped counts.
type PostOptions struct {
	Like, Comment, Share, Save string
	Likes, Comments, Shares    string
}

// Post shows author, text, media, and actions.
func Post(c *ui.Context, post PostData, state *PostState, o PostOptions) ui.Element {
	validPost("Post", post)
	if state == nil {
		panic(fmt.Sprintf("social: Post %q needs state", post.ID))
	}
	def(&o.Likes, "%s")
	def(&o.Comments, "%s")
	def(&o.Shares, "%s")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3), t.Space(3), t.Space(1.5), t.Space(3)).Children(func() {
			postBody(c, post)
			c.Key("counts")
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
				counts := []struct {
					mark, format string
					n            uint64
				}{{"heart", o.Likes, uint64(post.Likes) + uint64(unit(state.Liked))}, {"message-circle", o.Comments, uint64(post.Comments)}, {"repeat-2", o.Shares, uint64(post.Shares)}}
				for i, count := range counts {
					c.Key(i)
					ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
						mark := glyph(c, count.mark, 14)
						if i == 0 {
							mark.TextColor(mark.AnimateColor("liked", t.TextMuted.Mix(t.Danger, unit(state.Liked)), ui.Snappy))
						}
						kit.Numeric(c, fmt.Sprintf(count.format, number(c, count.n))).FontSize(t.Rem(0.78)).Shrink(0)
					})
				}
			})
		})
		ui.Box(c).FillWidth().Height(1).Background(t.Border)
		actions := PostActionsOptions{Like: o.Like, Comment: o.Comment, Share: o.Share, Save: o.Save}
		c.Key("actions")
		ui.Column(c).FillWidth().Padding(t.Space(1), t.Space(1.5)).Children(func() {
			if PostActions(c, state, actions).Changed() {
				root.SetChanged()
			}
		})
	})
	return root
}

// postBody is a post's byline, text and media.
func postBody(c *ui.Context, post PostData) {
	t := c.Theme()
	byline(c, post.Author, post.Time, t.Space(10))
	ui.Text(c, post.Text).FontSize(t.Rem(0.9)).Wrap().FillWidth()
	if post.Media != nil {
		c.Key("media")
		media := ui.Box(c).FillWidth().AspectRatio(16.0 / 9).Radius(t.Radius).Clip().Background(t.SurfaceHover)
		kit.Arrive(media, ui.Pose{Scale: 0.98, Fade: 1}, 80*time.Millisecond).Children(func() {
			ui.Image(c, post.Media).Fill().Fit(ui.Cover)
		})
	}
}
