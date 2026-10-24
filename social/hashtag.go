package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HashtagOptions are the hashtag count labels.
type HashtagOptions struct{ Posts string }

// Hashtag opens a tag chip with count; cue Tap.
func Hashtag(c *ui.Context, tag string, count int, opened *bool, o HashtagOptions) ui.Element {
	if tag == "" || count < 0 || opened == nil {
		panic("social: Hashtag needs tag, count, and state")
	}
	def(&o.Posts, "%d posts")
	t := c.Theme()
	root := button(c, tag).Row().Gap(t.Space(2)).Radius(999).Background(t.Accent.Alpha(0.12))
	root.Children(func() {
		ui.Text(c, tag).TextColor(t.Accent)
		kit.Numeric(c, fmt.Sprintf(o.Posts, count)).TextColor(t.TextMuted)
	})
	if root.Clicked() {
		*opened = true
		root.SetChanged()
		sound.Play(sound.Tap)
		kit.Log("Hashtag", "opened", "count", count)
	}
	return root
}
