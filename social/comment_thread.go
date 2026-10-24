package social

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Comment is one nested discussion comment.
type Comment struct {
	ID, Text string
	Author   User
	Replies  []Comment
}

// CommentThreadOptions are the reply labels.
type CommentThreadOptions struct{ Placeholder, Reply string }

// CommentThread shows nested comments and sends replies; cue Send.
func CommentThread(c *ui.Context, comments []Comment, draft *string, o CommentThreadOptions) (root, send ui.Element) {
	if len(comments) == 0 || draft == nil {
		panic("social: CommentThread needs comments and draft")
	}
	def(&o.Placeholder, "Write a reply")
	def(&o.Reply, "Reply")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
			for _, comment := range comments {
				commentView(c, comment, 0)
			}
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Children(func() {
			in := ui.TextArea(c, draft).Label(o.Placeholder).Placeholder(o.Placeholder).FillWidth().Lines(2, 5)
			if in.Changed() {
				root.SetChanged()
			}
			send = button(c, o.Reply).Background(t.Accent).TextColor(t.AccentText).Disabled(strings.TrimSpace(*draft) == "")
			send.Children(func() { ui.Text(c, o.Reply) })
		})
	})
	if send.Clicked() {
		send.SetChanged()
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("CommentThread", "replied", "runes", len([]rune(*draft)))
	}
	return root, send
}

func commentView(c *ui.Context, comment Comment, depth int) {
	validUser(comment.Author)
	if comment.ID == "" || comment.Text == "" {
		panic("social: invalid comment")
	}
	t := c.Theme()
	c.Key(comment.ID)
	card := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2)).MarginStart(t.Space(float32(depth) * 3)).
		Radius(t.Radius).Background(t.Surface)
	card.Children(func() {
		ui.Text(c, comment.Author.Name).Bold()
		ui.Text(c, comment.Text).Wrap()
		for _, reply := range comment.Replies {
			commentView(c, reply, depth+1)
		}
	})
}
