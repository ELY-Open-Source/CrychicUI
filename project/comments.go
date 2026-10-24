package project

import (
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Comment is one issue comment.
type Comment struct {
	Author, Text string
}

// CommentsOptions are the comment labels.
type CommentsOptions struct{ Placeholder, Send, Author string }

// Comments shows a thread with composer; cue Send.
func Comments(c *ui.Context, comments *[]Comment, draft *string, o CommentsOptions) (root, send ui.Element) {
	if comments == nil || draft == nil {
		panic("project: Comments needs comments and draft")
	}
	def(&o.Placeholder, "Write a comment")
	def(&o.Send, "Send")
	def(&o.Author, "You")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	pending := ui.Local(root, "pending", func() string { return "" })
	root.Children(func() {
		for i, comment := range *comments {
			if comment.Author == "" || comment.Text == "" {
				panic("project: Comments has an invalid comment")
			}
			c.Key(i)
			ui.Column(c).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, comment.Author).Bold()
				ui.Text(c, comment.Text).Wrap()
			})
		}
		c.Key("composer")
		ui.Row(c).Gap(t.Space(2)).Children(func() {
			in := ui.TextInput(c, draft).Label(o.Placeholder).Placeholder(o.Placeholder).Grow(1)
			in.OnTextKey(func(k ui.TextKey) bool {
				if k.Key != ui.KeyEnter || strings.TrimSpace(k.Text()) == "" {
					return false
				}
				*pending = strings.TrimSpace(k.Text())
				k.SetText("")
				return true
			})
			if in.Changed() {
				root.SetChanged()
			}
			send = kit.Press(ui.ButtonBase(c).Label(o.Send).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius).
				Background(t.Accent).TextColor(t.AccentText).Disabled(strings.TrimSpace(*draft) == ""))
			send.Children(func() { ui.Text(c, o.Send).SingleLine() })
		})
	})
	if *pending != "" || send.Clicked() {
		text := *pending
		if text == "" {
			text = strings.TrimSpace(*draft)
		}
		if text != "" {
			*comments = append(*comments, Comment{Author: o.Author, Text: text})
			*draft = ""
			send.SetChanged()
			root.SetChanged()
			sound.Play(sound.Send)
			kit.Log("Comments", "sent", "runes", utf8.RuneCountInString(text))
		}
		*pending = ""
	}
	return root, send
}
