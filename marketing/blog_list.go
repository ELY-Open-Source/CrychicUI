package marketing

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BlogListOptions word the list; Empty shows with no posts.
type BlogListOptions struct {
	Title, Empty string
	PostWords
}

// BlogList lists posts by title; cue Tap.
func BlogList(c *ui.Context, posts []Post, opened *string, o BlogListOptions) ui.Element {
	def(&o.Title, "Latest posts")
	def(&o.Empty, "No posts yet")
	o.fill()
	titles := make([]string, len(posts))
	for i, p := range posts {
		validPost("BlogList", p)
		titles[i] = p.Title
	}
	unique("BlogList", titles)
	if opened == nil || *opened != "" && !slices.Contains(titles, *opened) {
		panic(fmt.Sprintf("marketing: BlogList %q needs an opened post among its posts", o.Title))
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(1)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(0, t.Space(1), t.Space(2), t.Space(1)).Children(func() {
			glyph(c, "newspaper", 16).TextColor(t.TextMuted)
			ui.Text(c, o.Title).FontSize(t.Rem(0.95)).Bold().Grow(1).MinWidth(0).SingleLine()
			kit.Chip(c, strconv.Itoa(len(posts)), t.TextMuted)
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(1))
		list.Children(func() {
			if len(posts) == 0 {
				ui.Column(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(2)).Padding(t.Space(6), 0).Children(func() {
					glyph(c, "newspaper", 28).TextColor(t.TextMuted.Alpha(0.6))
					kit.Muted(c, o.Empty).FontSize(t.Rem(0.84)).Shrink(0)
				})
				return
			}
			for i, p := range posts {
				on := p.Title == *opened
				c.Key(strconv.Quote(p.Title))
				b := kit.Press(ui.ButtonBase(c).Label(p.Title).Checked(on).FillWidth().Row().Gap(t.Space(3)).AlignItems(ui.Center).
					Padding(t.Space(2)).Radius(t.Radius))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent, 0.08*unit(on)), ui.Snappy))
				nudge := b.AnimateTo("nudge", unit(b.Hovered()), ui.Snappy)
				kit.Stagger(c, list, b, i).Children(func() {
					thumb := ui.Box(c).Size(t.Space(16), t.Space(12)).Shrink(0).Center().Radius(t.Radius).Clip().Background(t.SurfaceHover)
					thumb.Children(func() {
						if p.Cover != nil {
							ui.Image(c, p.Cover).Fill().Fit(ui.Cover)
						} else {
							glyph(c, "newspaper", 18).TextColor(t.TextMuted)
						}
					})
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Row(c).FillWidth().MinWidth(0).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							if p.Tag != "" {
								ui.Text(c, p.Tag).FontSize(t.Rem(0.68)).Bold().TextColor(t.Accent).Shrink(0).MaxWidth(t.Space(24)).SingleLine()
								dot(c)
							}
							kit.Muted(c, p.Author).FontSize(t.Rem(0.68)).Grow(1).MinWidth(0).SingleLine()
						})
						ui.Text(c, p.Title).FontSize(t.Rem(0.86)).Bold().TextColor(t.Text.Mix(t.Accent, unit(on))).FillWidth().MinWidth(0).SingleLine()
						meta(c, p, o.PostWords)
					})
					glyph(c, "chevron-right", 16).TextColor(t.TextMuted).Translate(t.Space(0.75)*nudge, 0)
				})
				if b.Clicked() {
					sound.Play(sound.Tap)
					kit.Log("BlogList", "opened", "index", i)
					if !on {
						*opened = p.Title
						root.SetChanged()
					}
				}
			}
		})
	})
	return root
}
