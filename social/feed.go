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
	ids := make([]string, len(posts))
	for i, post := range posts {
		validPost("Feed", post)
		ids[i] = post.ID
	}
	unique("Feed", ids)
	def(&o.New, "New posts")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2.5)).Role(ui.RoleGroup)
	root.Children(func() {
		c.Key("new")
		ui.Column(c).FillWidth().AlignItems(ui.Center).Children(func() {
			if !*fresh {
				return
			}
			pill := kit.Press(ui.ButtonBase(c).Label(o.New).Padding(t.Space(1.5), t.Space(3.5)).Radius(999).
				Shadow(0, 2, 8, 0, t.Text.Alpha(0.16)).TextColor(t.AccentText))
			pill.Background(pill.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(pill.Hovered())), ui.Snappy))
			pill.Transition(ui.ElementTransition{Exit: &ui.Motion{Y: -t.Space(3)}})
			kit.Arrive(pill, ui.Pose{Y: -t.Space(3), Fade: 1}, 0).Children(func() {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					glyph(c, "arrow-up", 14)
					ui.Text(c, o.New).FontSize(t.Rem(0.8)).Bold().Shrink(0)
				})
			})
			if pill.Clicked() {
				*fresh = false
				root.SetChanged()
				sound.Play(sound.Notify)
				kit.Log("Feed", "revealed", "count", len(posts))
			}
		})
		c.Key("posts")
		list := ui.Column(c).FillWidth().Gap(t.Space(2.5))
		list.Children(func() {
			for i, post := range posts {
				c.Key(post.ID)
				row := ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).Background(t.Background)
				kit.Stagger(c, list, row, i).Children(func() { postBody(c, post) })
			}
		})
	})
	return root
}
