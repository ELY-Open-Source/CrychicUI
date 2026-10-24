package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FeedOptions are the feed labels.
type FeedOptions struct{ New string }

// Feed lists posts and reveals new posts; cue Notify.
func Feed(c *ui.Context, posts []PostData, fresh *bool, o FeedOptions) ui.Element {
	if len(posts) == 0 || fresh == nil {
		panic("social: Feed needs posts and state")
	}
	def(&o.New, "New posts")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		if *fresh {
			pill := button(c, o.New).Radius(999).Background(t.Accent).TextColor(t.AccentText)
			pill.Children(func() { ui.Text(c, o.New) })
			if pill.Clicked() {
				*fresh = false
				root.SetChanged()
				sound.Play(sound.Notify)
				kit.Log("Feed", "revealed", "count", len(posts))
			}
		}
		for _, post := range posts {
			validPost(post)
			c.Key(post.ID)
			card := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border)
			kit.Arrive(card, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
			card.Children(func() {
				ui.Text(c, post.Author.Name).Bold()
				ui.Text(c, post.Text).Wrap()
				if post.Media != nil {
					ui.Image(c, post.Media).FillWidth().AspectRatio(1.8).Fit(ui.Cover).Radius(t.Radius).Clip()
				}
			})
		}
	})
	return root
}
