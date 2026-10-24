package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PostOptions are the post labels and count formats.
type PostOptions struct {
	Like, Comment, Share, Save string
	Likes, Comments, Shares    string
}

// Post shows author, text, media, and actions.
func Post(c *ui.Context, post PostData, state *PostState, o PostOptions) ui.Element {
	validPost(post)
	if state == nil {
		panic("social: Post needs state")
	}
	def(&o.Likes, "%d likes")
	def(&o.Comments, "%d comments")
	def(&o.Shares, "%d shares")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			avatar(c, post.Author, t.Space(9))
			ui.Column(c).Children(func() {
				ui.Text(c, post.Author.Name).Bold()
				kit.Muted(c, post.Author.Handle).SingleLine()
			})
		})
		ui.Text(c, post.Text).Wrap()
		if post.Media != nil {
			ui.Image(c, post.Media).FillWidth().AspectRatio(1.6).Fit(ui.Cover).Radius(t.Radius).Clip()
		}
		c.Key("counts")
		ui.Row(c).Gap(t.Space(3)).Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Likes, post.Likes)).TextColor(t.TextMuted)
			kit.Numeric(c, fmt.Sprintf(o.Comments, post.Comments)).TextColor(t.TextMuted)
			kit.Numeric(c, fmt.Sprintf(o.Shares, post.Shares)).TextColor(t.TextMuted)
		})
		actions := PostActionsOptions{Like: o.Like, Comment: o.Comment, Share: o.Share, Save: o.Save}
		c.Key("actions")
		if PostActions(c, state, actions).Changed() {
			root.SetChanged()
		}
	})
	return root
}
