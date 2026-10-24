package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Story is one user's story ring.
type Story struct {
	ID   string
	User User
	Seen bool
}

// StoryRing opens unseen story avatars; cue Open.
func StoryRing(c *ui.Context, stories []Story, opened *string) ui.Element {
	if len(stories) == 0 || opened == nil {
		panic("social: StoryRing needs stories and state")
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, story := range stories {
			if story.ID == "" {
				panic("social: story needs id")
			}
			validUser(story.User)
			c.Key(story.ID)
			b := button(c, story.User.Name).Column().Gap(t.Space(1)).AlignItems(ui.Center).Radius(999)
			b.Children(func() {
				avatar(c, story.User, t.Space(12)).Border(1, t.Border.Mix(t.Accent, unit(!story.Seen)))
				ui.Text(c, story.User.Name).SingleLine()
			})
			if b.Clicked() {
				*opened = story.ID
				root.SetChanged()
				sound.Play(sound.Open)
				kit.Log("StoryRing", "opened")
			}
		}
	})
	return root
}
