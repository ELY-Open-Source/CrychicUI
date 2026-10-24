package project

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Comment is one issue comment; When is optional.
type Comment struct {
	Author, Text, When string
}

// CommentsOptions are the comment labels; Author signs new comments.
type CommentsOptions struct{ Title, Placeholder, Send, Author, Empty string }

// Comments shows a thread with composer; cue Send.
func Comments(c *ui.Context, comments *[]Comment, draft *string, o CommentsOptions) (root, send ui.Element) {
	if comments == nil || draft == nil {
		panic("project: Comments needs comments and draft")
	}
	for i, comment := range *comments {
		if comment.Author == "" || comment.Text == "" {
			panic(fmt.Sprintf("project: Comments has an invalid comment %d", i))
		}
	}
	def(&o.Title, "Comments")
	def(&o.Placeholder, "Write a comment")
	def(&o.Send, "Send")
	def(&o.Author, "You")
	def(&o.Empty, "No comments yet")
	t := c.Theme()
	root = card(c)
	pending := ui.Local(root, "pending", func() string { return "" })
	root.Children(func() {
		kit.CardHead(c, "messages-square", o.Title, func() { kit.Chip(c, fmt.Sprint(len(*comments)), t.TextMuted) })
		thread := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3))
		thread.Children(func() {
			if len(*comments) == 0 {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
					glyph(c, "message-square-dashed", 16)
					ui.Text(c, o.Empty).FontSize(t.Rem(0.8))
				})
			}
			for i, comment := range *comments {
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Start)
				kit.Stagger(c, thread, row, i).Children(func() {
					display.Avatar(c, comment.Author, nil).Size(t.Space(8), t.Space(8)).FontSize(t.Rem(0.65))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Text(c, comment.Author).FontSize(t.Rem(0.8)).Bold().Shrink(0)
							if comment.When != "" {
								kit.Muted(c, comment.When).FontSize(t.Rem(0.7)).Grow(1).MinWidth(0).SingleLine()
							}
						})
						ui.Text(c, comment.Text).FontSize(t.Rem(0.8)).Wrap().FillWidth().Padding(t.Space(2), t.Space(2.5)).
							Radius(t.Radius).Background(t.SurfaceHover.Alpha(0.6))
					})
				})
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(3)).
			BorderWidth(1, 0, 0, 0).BorderColor(t.Border).Children(func() {
			display.Avatar(c, o.Author, nil).Size(t.Space(7), t.Space(7)).FontSize(t.Rem(0.6))
			field := ui.Row(c).Grow(1).MinWidth(0).Padding(t.Space(1), t.Space(2.5)).Radius(999).Border(1, t.Border)
			field.Children(func() {
				in := ui.TextInput(c, draft).Label(o.Placeholder).Placeholder(o.Placeholder).Grow(1).MinWidth(0)
				in.OnTextKey(func(k ui.TextKey) bool {
					if k.Key != ui.KeyEnter || strings.TrimSpace(k.Text()) == "" {
						return false
					}
					*pending = strings.TrimSpace(k.Text())
					k.SetText("")
					return true
				})
				field.BorderColor(field.AnimateColor("line", t.Border.Mix(t.Accent, unit(in.Focused())), ui.Snappy))
				if in.Changed() {
					root.SetChanged()
				}
			})
			empty := strings.TrimSpace(*draft) == ""
			send = kit.Press(ui.ButtonBase(c).Label(o.Send).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).
				TextColor(t.AccentText).Disabled(empty))
			send.Background(send.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(send.Hovered())), ui.Snappy))
			send.Children(func() { glyph(c, "send", 14) })
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
