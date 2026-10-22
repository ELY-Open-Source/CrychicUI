package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Feature is one page of a WhatsNew.
type Feature struct{ Icon, Title, Text string }

// WhatsNewOptions tune a WhatsNew.
type WhatsNewOptions struct {
	Back, Next, Done string
}

// WhatsNew pages through new features in a modal.
func WhatsNew(c *ui.Context, title string, features []Feature, page *int, open *bool, o WhatsNewOptions) ui.Element {
	if title == "" || len(features) == 0 || page == nil || open == nil || !valid(*page, len(features)) {
		panic(fmt.Sprintf("overlay: WhatsNew %q needs features, a page and open", title))
	}
	for _, f := range features {
		if f.Title == "" {
			panic(fmt.Sprintf("overlay: WhatsNew %q features need titles", title))
		}
		icon.Must(f.Icon)
	}
	def(&o.Back, "Back")
	def(&o.Next, "Next")
	def(&o.Done, "Done")
	t := c.Theme()
	root := ui.Box(c)
	shown, last := *page, len(features)-1
	turn := func(to int) {
		*page = to
		root.SetChanged()
		sound.Play(sound.Swipe)
		kit.Log("WhatsNew", "page", "index", to)
	}
	if dialog(c, open, t.Space(110), func() {
		kit.Muted(c, title).Bold().FontSize(t.Rem(0.8))
		kit.Slide(c, shown, func() {
			f := features[shown]
			ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(2), 0).Children(func() {
				ui.Box(c).Size(t.Space(12), t.Space(12)).Center().Radius(t.Radius).Background(t.Accent.Alpha(0.12)).Children(func() {
					ui.Icon(c, icon.Must(f.Icon)).Size(16, 16).TextColor(t.Accent)
				})
				ui.Text(c, f.Title).Bold().FontSize(t.Rem(1.2))
				kit.Muted(c, f.Text)
			})
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Row(c).Grow(1).Gap(t.Space(1.5)).Children(func() {
				for i := range features {
					c.Key(i)
					dot := ui.Box(c)
					k := dot.AnimateTo("on", unit(i == shown), ui.Snappy)
					dot.Size(t.Space(2)+k*t.Space(4), t.Space(2)).Radius(999).Background(t.Border.Mix(t.Accent, k))
				}
			})
			if quiet(c.Key("back"), o.Back).Disabled(shown == 0).Clicked() {
				turn(shown - 1)
			}
			if shown < last {
				if primary(c.Key("next"), o.Next, t.Accent, t.AccentText).Clicked() {
					turn(shown + 1)
				}
				return
			}
			if button.Solid(c.Key("done"), o.Done).Clicked() {
				*open = false
				root.SetChanged()
			}
		})
	}).Dismissed() {
		root.SetChanged()
	}
	moved(root, *open, "WhatsNew")
	return root
}
