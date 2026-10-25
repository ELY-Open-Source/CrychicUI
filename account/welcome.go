package account

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Slide is one welcome page.
type Slide struct {
	Icon, Title, Text string
}

// WelcomeOptions word the buttons and the finished state.
type WelcomeOptions struct{ Next, Start, Skip, Done, DoneNote string }

// Welcome steps through slides, then shows done; cue Tap.
func Welcome(c *ui.Context, slides []Slide, page *int, o WelcomeOptions) ui.Element {
	if len(slides) == 0 || page == nil || *page < 0 || *page > len(slides) {
		panic("account: Welcome needs slides and a page within them")
	}
	for i, s := range slides {
		if s.Icon == "" || s.Title == "" || s.Text == "" {
			panic(fmt.Sprintf("account: Welcome slide %d %q needs an icon, title and text", i, s.Title))
		}
	}
	def(&o.Next, "Next")
	def(&o.Start, "Get started")
	def(&o.Skip, "Skip")
	def(&o.Done, "Welcome aboard")
	def(&o.DoneNote, "Your workspace is ready.")
	t := c.Theme()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(6), t.Space(5), t.Space(5), t.Space(5)).AlignItems(ui.Center)
	root.Children(func() {
		mark, title, text, ink := "party-popper", o.Done, o.DoneNote, t.Success
		if *page < len(slides) {
			s := slides[*page]
			mark, title, text, ink = s.Icon, s.Title, s.Text, t.Accent
		}
		c.Key(*page)
		face := ui.Column(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center)
		kit.Arrive(face, ui.Pose{X: t.Space(6), Fade: 1}, 0).Children(func() {
			halo := ui.Box(c).Size(t.Space(18), t.Space(18)).Center().Radius(999).Background(ink.Alpha(0.1)).TextColor(ink)
			kit.Arrive(halo, ui.Pose{Scale: 0.6, Rotate: -10}, 80*time.Millisecond).Children(func() { glyph(c, mark, 32) })
			ui.Text(c, title).FontSize(t.Rem(1.15)).Bold().Wrap().TextAlign(ui.Center)
			kit.Muted(c, text).FontSize(t.Rem(0.85)).Wrap().TextAlign(ui.Center).MaxWidth(t.Space(80))
		})
		c.Key("dots")
		ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
			for i := range slides {
				d := ui.Box(c).Height(t.Space(2)).Radius(999)
				d.Width(d.AnimateTo("w", t.Space(2)+t.Space(4)*unit(i == *page), ui.Snappy))
				d.Background(d.AnimateColor("ink", t.Border.Mix(t.Accent, unit(i <= *page)), ui.Snappy))
			}
		})
		if *page == len(slides) {
			return
		}
		c.Key("actions")
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			last := *page == len(slides)-1
			if !last {
				skip := outline(c, o.Skip, o.Skip, "", t.TextMuted)
				if skip.Clicked() {
					*page = len(slides)
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("Welcome", "skipped")
				}
			}
			label, arrow := o.Next, "arrow-right"
			if last {
				label, arrow = o.Start, "rocket"
			}
			c.Key("next")
			next := primary(c, label, arrow).Grow(1)
			if next.Clicked() {
				*page++
				root.SetChanged()
				sound.Play(sound.Tap)
				kit.Log("Welcome", "page", "index", *page)
			}
		})
	})
	return root
}
