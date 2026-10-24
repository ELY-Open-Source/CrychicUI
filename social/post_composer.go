package social

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PostComposerOptions tune the composer labels and limit.
type PostComposerOptions struct {
	Placeholder, Media, Send, Count string
	Limit                           int
}

// PostComposer edits text with count and media; cue Send.
func PostComposer(c *ui.Context, draft *string, media ui.ImageSource, mediaAdded *bool, o PostComposerOptions) (root, send ui.Element) {
	if draft == nil || mediaAdded == nil || o.Limit < 0 {
		panic("social: PostComposer needs draft, media state, and limit")
	}
	def(&o.Placeholder, "What's happening?")
	def(&o.Media, "Add media")
	def(&o.Send, "Post")
	def(&o.Count, "%d/%d")
	if o.Limit == 0 {
		o.Limit = 280
	}
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		in := ui.TextArea(c, draft).Label(o.Placeholder).Placeholder(o.Placeholder).FillWidth().Lines(3, 6)
		if in.Changed() {
			root.SetChanged()
		}
		count := len([]rune(*draft))
		c.Key("count")
		kit.Numeric(c, fmt.Sprintf(o.Count, count, o.Limit)).TextColor(t.TextMuted)
		if media != nil {
			ui.Image(c, media).FillWidth().AspectRatio(1.8).Fit(ui.Cover).Radius(t.Radius).Clip()
		}
		c.Key("controls")
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			c.Key("media")
			add := button(c, o.Media)
			add.Children(func() { ui.Text(c, o.Media) })
			if add.Clicked() {
				*mediaAdded = true
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("PostComposer", "media")
			}
			valid := strings.TrimSpace(*draft) != "" && count <= o.Limit
			c.Key("send")
			send = button(c, o.Send).Background(t.Accent).TextColor(t.AccentText).Disabled(!valid)
			send.Children(func() { ui.Text(c, o.Send) })
		})
	})
	if send.Clicked() {
		send.SetChanged()
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("PostComposer", "sent", "runes", len([]rune(*draft)))
	}
	return root, send
}
