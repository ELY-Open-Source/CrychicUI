package chat

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	"github.com/egoist/mygo/ui/markdown"
	"github.com/egoist/mygo/ui/markdown/parser"
)

// MarkdownMessage is a bubble rendered with ui/markdown.
func MarkdownMessage(c *ui.Context, m Message) ui.Element {
	messages("MarkdownMessage", []Message{m})
	t := c.Theme()
	root := shell(c, false)
	if m.Mine {
		root.AlignSelf(ui.End)
	}
	md := ui.Local(root, "md", func() *markdown.Options { return &markdown.Options{} })
	anchor := ui.Local(root, "anchor", func() string { return "" })
	left := ui.Local(root, "left", func() int { return 0 })
	done := ui.Local(root, "done", func() bool { return false })
	(*md).OnLink = func(url string) {
		switch {
		case strings.HasPrefix(url, "#"):
			*anchor = strings.ToLower(url[1:])
			*left = 2
			c.After(time.Millisecond)
			sound.Play(sound.Open)
			kit.Log("MarkdownMessage", "anchor")
		case ui.SafeToOpen(url):
			c.OpenURL(url)
			sound.Play(sound.Open)
			kit.Log("MarkdownMessage", "link")
		default:
			r := c.Router()
			bad := strings.ContainsAny(url, "\\\x00") || strings.HasPrefix(url, "//")
			if i := strings.IndexAny(url, ":/?#"); r != nil && !bad && (i < 0 || url[i] != ':') {
				r.Push(url)
				sound.Play(sound.Open)
				kit.Log("MarkdownMessage", "route")
			}
		}
	}
	th := markdown.DefaultTheme(t)
	(*md).CodeBlock = func(c *ui.Context, b *markdown.CodeBlockView) {
		box := ui.Column(c).AlignSelf(ui.Stretch).Clip().Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		box.Role(ui.RoleGroup).Label(b.Language).Direction(ui.LeftToRight).Children(func() {
			ui.Row(c).AlignItems(ui.Center).Padding(t.Space(1), t.Space(2)).Children(func() {
				kit.Muted(c, b.Language).FontSize(t.Rem(0.8)).Grow(1)
				if cp := b.CopyButton(c); cp.Clicked() {
					sound.Play(sound.Pop)
					kit.Log("MarkdownMessage", "copied", "runes", utf8.RuneCountInString(b.Code()))
				}
			})
			b.Body(c)
		})
	}
	(*md).Block = func(c *ui.Context, d *parser.Document, i int) bool {
		b := d.Block(i)
		if *anchor == "" || b.Kind != parser.Heading || b.Level < 1 || b.Level > 6 {
			return false
		}
		hs := th.Headings[b.Level-1]
		text := parser.PlainText(d.Inlines(i))
		h := ui.Column(c).Gap(t.Space(1))
		h.Children(func() {
			e := ui.Text(c, text).FontSize(hs.Size * t.FontSize).FontWeight(hs.Weight)
			if hs.Color.A != 0 {
				e.TextColor(hs.Color)
			}
			if hs.Rule {
				ui.Box(c).FillWidth().Height(1).Background(th.Rule)
			}
		})
		if markdown.Slug(text) == *anchor && !*done {
			*done = true
			h.ScrollIntoView()
			c.After(time.Millisecond)
		}
		return true
	}
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	*done = false
	out := root.Children(func() {
		markdown.View(c, m.Text, *md)
	})
	if *anchor != "" {
		*left--
		if *left <= 0 {
			*anchor = ""
			c.After(time.Millisecond)
		}
	}
	return out
}
