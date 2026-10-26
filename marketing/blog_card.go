package marketing

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// BlogCardOptions word the card.
type BlogCardOptions struct{ PostWords }

// BlogCard is a post card with a cover; cue Tap.
func BlogCard(c *ui.Context, p Post, o BlogCardOptions) ui.Element {
	validPost("BlogCard", p)
	o.fill()
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(p.Title).Column().AlignItems(ui.Stretch).FillWidth().Radius(t.Radius * 1.5).Clip())
	hover := root.AnimateTo("hover", unit(root.Hovered()), ui.Snappy)
	root.Background(t.Background).Border(1, t.Border.Mix(t.Accent, hover)).Shadow(0, 1, 3, 0, t.Text.Alpha(0.06))
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0).Children(func() {
		if p.Cover != nil {
			ui.Box(c).FillWidth().AspectRatio(16.0 / 9).Clip().Background(t.SurfaceHover).Children(func() {
				ui.Image(c, p.Cover).Fill().Fit(ui.Cover).Scale(1+0.05*hover, 1+0.05*hover)
			})
		}
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(4)).Children(func() {
			if p.Tag != "" {
				ui.Text(c, p.Tag).FontSize(t.Rem(0.72)).Bold().TextColor(t.Accent).FillWidth().MinWidth(0).SingleLine()
			}
			ui.Text(c, p.Title).FontSize(t.Rem(1)).Bold().FillWidth().Wrap()
			if p.Excerpt != "" {
				kit.Muted(c, p.Excerpt).FontSize(t.Rem(0.8)).FillWidth().Wrap()
			}
			ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1), 0, 0, 0).Children(func() {
				avatar(c, p.Author, p.Photo, 8)
				ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
					ui.Text(c, p.Author).FontSize(t.Rem(0.8)).Bold().FillWidth().MinWidth(0).SingleLine()
					meta(c, p, o.PostWords)
				})
			})
		})
	})
	if root.Clicked() {
		sound.Play(sound.Tap)
		kit.Log("BlogCard", "opened")
	}
	return root
}
