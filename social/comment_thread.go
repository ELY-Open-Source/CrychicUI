package social

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Comment is one nested discussion comment.
type Comment struct {
	ID, Text, Time string
	Author         User
	Replies        []Comment
}

// CommentThreadOptions are the reply labels.
type CommentThreadOptions struct{ Placeholder, Reply string }

// CommentThread shows nested comments and sends replies; cue Send.
func CommentThread(c *ui.Context, comments []Comment, draft *string, o CommentThreadOptions) (root, send ui.Element) {
	if len(comments) == 0 || draft == nil {
		panic("social: CommentThread needs comments and draft")
	}
	unique("CommentThread", commentIDs(comments, nil))
	def(&o.Placeholder, "Write a reply")
	def(&o.Reply, "Reply")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	root.Children(func() {
		list := ui.Column(c).FillWidth().Gap(t.Space(3))
		list.Children(func() {
			for i, comment := range comments {
				commentView(c, list, comment, i)
			}
		})
		ui.Column(c).FillWidth().Children(func() {
			field := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.End).Padding(t.Space(2), t.Space(2), t.Space(2), t.Space(3)).
				Radius(t.Radius).Background(t.Background)
			field.Border(1, field.AnimateColor("ring", t.Border.Mix(t.Accent, unit(field.FocusWithin())), ui.Snappy))
			field.Children(func() {
				in := ui.TextAreaBase(c, draft).Label(o.Placeholder).Placeholder(o.Placeholder).Grow(1).MinWidth(0).Lines(1, 5).FontSize(t.Rem(0.85))
				if in.Changed() {
					root.SetChanged()
				}
				send = kit.Press(ui.ButtonBase(c).Label(o.Reply).Size(t.Space(8), t.Space(8)).Shrink(0).Radius(999).TextColor(t.AccentText).
					Disabled(strings.TrimSpace(*draft) == ""))
				send.Background(send.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(send.Hovered())), ui.Snappy))
				send.Children(func() { glyph(c, "send", 16) })
			})
		})
	})
	if send.Clicked() {
		send.SetChanged()
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("CommentThread", "replied", "runes", utf8.RuneCountInString(*draft))
	}
	return root, send
}

func commentView(c *ui.Context, list ui.Element, comment Comment, i int) {
	validUser("CommentThread", comment.Author)
	if comment.ID == "" || comment.Text == "" {
		panic(fmt.Sprintf("social: CommentThread comment %q needs id and text", comment.ID))
	}
	t := c.Theme()
	c.Key(comment.ID)
	row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Start)
	kit.Stagger(c, list, row, i).Children(func() {
		avatar(c, comment.Author, t.Space(8)).AlignSelf(ui.Start)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
			ui.Column(c).Gap(t.Space(0.5)).Padding(t.Space(2), t.Space(3)).Radius(t.Radius * 1.5).Background(t.Surface).Children(func() {
				ui.Text(c, comment.Author.Name).FontSize(t.Rem(0.8)).Bold().Shrink(0)
				ui.Text(c, comment.Text).FontSize(t.Rem(0.85)).Wrap().FillWidth()
			})
			if comment.Time != "" {
				kit.Muted(c, comment.Time).FontSize(t.Rem(0.7)).Padding(0, t.Space(3))
			}
			if len(comment.Replies) == 0 {
				return
			}
			replies := ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(1), 0, 0, t.Space(3)).
				BorderWidth(0, 0, 0, 1).BorderColor(t.Border)
			replies.Children(func() {
				for j, reply := range comment.Replies {
					commentView(c, replies, reply, j)
				}
			})
		})
	})
}

// commentIDs collects every comment and reply id.
func commentIDs(comments []Comment, ids []string) []string {
	for _, comment := range comments {
		ids = commentIDs(comment.Replies, append(ids, comment.ID))
	}
	return ids
}
