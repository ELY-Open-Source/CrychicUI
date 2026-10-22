package navigation

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RecentNavOptions tune a RecentNav.
type RecentNavOptions struct {
	Back, Forward, Recent string
}

// seen is a page RecentNav saw shown.
type seen struct{ loc, title string }

// RecentNav is back, forward and a menu of recent pages.
func RecentNav(c *ui.Context, label string, r *ui.Router, open *bool, o RecentNavOptions) ui.Element {
	if label == "" || r == nil || open == nil {
		panic(fmt.Sprintf("navigation: RecentNav %q needs a router and open", label))
	}
	def(&o.Back, "Back")
	def(&o.Forward, "Forward")
	def(&o.Recent, "Recent pages")
	t := c.Theme()
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).AlignItems(ui.Center).FocusGroup(ui.Horizontal).Label(label)
	recent := ui.Local(root, "recent", func() []seen { return nil })
	now := seen{r.Location(), r.Title()}
	if now.title == "" {
		now.title = r.Path()
	}
	// Newest first; a page seen again moves up.
	if len(*recent) == 0 || (*recent)[0] != now {
		rest := slices.DeleteFunc(*recent, func(s seen) bool { return s.loc == now.loc })
		*recent = append([]seen{now}, rest[:min(len(rest), 7)]...)
	}
	was, list := *open, slices.Clone(*recent)
	moved := func(event string) {
		root.SetChanged()
		kit.Log("RecentNav", event)
	}
	root.Children(func() {
		back, forth := "chevron-left", "chevron-right"
		if c.Direction() == ui.RightToLeft {
			back, forth = forth, back
		}
		if arrowButton(c, o.Back, back, !r.CanGoBack(), sound.Tap) {
			r.Back()
			moved("back")
		}
		if arrowButton(c, o.Forward, forth, !r.CanGoForward(), sound.Tap) {
			r.Forward()
			moved("forward")
		}
		trigger := kit.Press(ui.ButtonBase(c).Label(o.Recent).Gap(t.Space(1.5)).Padding(t.Space(1.5), t.Space(2.5)).Radius(t.Radius)).Expanded(*open)
		trigger.Background(trigger.AnimateColor("bg", t.SurfaceHover.Alpha(unit(*open || trigger.Hovered())), ui.Snappy)).Children(func() {
			ui.Text(c, now.title).Bold().SingleLine().ContentTransition(ui.ContentTransition{})
			chev := ui.Icon(c, icon.Must("chevron-down")).Size(16, 16)
			chev.Rotate(180 * chev.AnimateTo("turn", unit(*open), ui.Snappy))
		})
		popup(c, trigger, open, func() {
			ui.Column(c).Width(t.Space(56)).Gap(t.Space(0.5)).Children(func() {
				for _, s := range list {
					c.Key(s.loc)
					if kit.Clicked(checked(c, s.title, s.loc == now.loc, nil), sound.Tap) {
						*open = false
						if s.loc != now.loc {
							r.Push(s.loc)
							moved("opened recent")
						}
					}
				}
			})
		})
	})
	reopened(root, was, *open, "RecentNav")
	return root
}
