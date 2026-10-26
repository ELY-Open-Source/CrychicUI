package marketing

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Quote is one testimonial; zero Rating hides the stars.
type Quote struct {
	Text, Name, Role string
	Photo            *ui.Bitmap
	Rating           int
}

// TestimonialsOptions word the carousel; Dot takes the quote number.
type TestimonialsOptions struct{ Title, Prev, Next, Dot string }

// Testimonials pages through quotes; cue Swipe.
func Testimonials(c *ui.Context, quotes []Quote, at *int, o TestimonialsOptions) ui.Element {
	def(&o.Prev, "Previous quote")
	def(&o.Next, "Next quote")
	def(&o.Dot, "Quote %d")
	if len(quotes) == 0 || at == nil || *at < 0 || *at >= len(quotes) {
		panic(fmt.Sprintf("marketing: Testimonials %q needs quotes and a position among them", o.Title))
	}
	for _, q := range quotes {
		if strings.TrimSpace(q.Text) == "" || strings.TrimSpace(q.Name) == "" || q.Rating < 0 || q.Rating > 5 {
			panic(fmt.Sprintf("marketing: Testimonials quote by %q needs text, a name and a rating of 0 to 5", q.Name))
		}
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(6), t.Space(5)).Role(ui.RoleGroup)
	move := func(to int) {
		*at = to
		root.SetChanged()
		sound.Play(sound.Swipe)
		kit.Log("Testimonials", "moved", "index", to)
	}
	root.Children(func() {
		if o.Title != "" {
			ui.Text(c, o.Title).FontSize(t.Rem(1.1)).Bold().FillWidth().TextAlign(ui.Center).Wrap()
		}
		q := quotes[*at]
		kit.Slide(c, *at, func() {
			ui.Column(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
				glyph(c, "quote", 28).TextColor(t.Accent.Alpha(0.6))
				if q.Rating > 0 {
					ui.Row(c).Gap(t.Space(0.5)).TextColor(t.Warning).Children(func() {
						for range q.Rating {
							glyph(c, "star", 14)
						}
					})
				}
				ui.Text(c, q.Text).FontSize(t.Rem(1)).FillWidth().TextAlign(ui.Center).Wrap()
				ui.Row(c).Gap(t.Space(2.5)).AlignItems(ui.Center).Children(func() {
					avatar(c, q.Name, q.Photo, 9)
					ui.Column(c).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
						ui.Text(c, q.Name).FontSize(t.Rem(0.84)).Bold().Shrink(0).MaxWidth(t.Space(60)).SingleLine()
						if q.Role != "" {
							kit.Muted(c, q.Role).FontSize(t.Rem(0.72)).Shrink(0).MaxWidth(t.Space(60)).SingleLine()
						}
					})
				})
			})
		})
		c.Key("controls")
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Justify(ui.Center).Children(func() {
			c.Key("prev")
			prev := arrow(c, o.Prev, "chevron-left").Disabled(*at == 0)
			if prev.Clicked() && *at > 0 {
				move(*at - 1)
			}
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				for i := range quotes {
					on := i == *at
					d := ui.ButtonBase(c).Label(fmt.Sprintf(o.Dot, i+1)).Checked(on).Height(t.Space(2)).Radius(999)
					d.Width(t.Space(2) + t.Space(4)*d.AnimateTo("wide", unit(on), ui.Snappy))
					d.Background(d.AnimateColor("bg", t.Border.Mix(t.TextMuted, unit(d.Hovered())).Mix(t.Accent, unit(on)), ui.Snappy))
					if d.Clicked() && !on {
						move(i)
					}
				}
			})
			c.Key("next")
			next := arrow(c, o.Next, "chevron-right").Disabled(*at == len(quotes)-1)
			if next.Clicked() && *at < len(quotes)-1 {
				move(*at + 1)
			}
		})
	})
	return root
}
