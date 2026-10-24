package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BannerStat is one figure in the day's summary.
type BannerStat struct{ Label, Value string }

// WelcomeBannerOptions add an icon and up to four stats.
type WelcomeBannerOptions struct {
	Icon  string
	Stats []BannerStat
}

// WelcomeBanner shows a greeting with the day's summary and stats.
func WelcomeBanner(c *ui.Context, greeting, summary string, o WelcomeBannerOptions) ui.Element {
	if greeting == "" || summary == "" || len(o.Stats) > 4 {
		panic(fmt.Sprintf("dashboard: WelcomeBanner %q needs a greeting, a summary and at most four stats", greeting))
	}
	for _, s := range o.Stats {
		if s.Label == "" || s.Value == "" {
			panic(fmt.Sprintf("dashboard: WelcomeBanner %q stat needs a label and a value", greeting))
		}
	}
	def(&o.Icon, "sun")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).Radius(t.Radius).Background(t.Accent).TextColor(t.AccentText).
		Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(3), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(11), t.Space(11)).Shrink(0).Center().Radius(999).Background(t.AccentText.Alpha(0.18)).
				Children(func() { ui.Icon(c, icon.Must(o.Icon)).Size(20, 20) })
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, greeting).Bold().FontSize(t.Rem(1.4)).FillWidth().Wrap()
				ui.Text(c, summary).FontSize(t.Rem(0.85)).TextColor(t.AccentText.Alpha(0.85)).FillWidth().Wrap()
			})
		})
		if len(o.Stats) > 0 {
			stats := ui.Row(c).FillWidth().Gap(t.Space(2))
			stats.Children(func() {
				for i, s := range o.Stats {
					c.Key(i)
					box := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(0.25)).Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius)
					box.Background(box.AnimateColor("bg", t.AccentText.Alpha(0.14+0.08*unit(box.Hovered())), ui.Snappy))
					kit.Stagger(c, stats, box, i).Children(func() {
						kit.Numeric(c, s.Value).Bold().FontSize(t.Rem(1.1)).FillWidth().SingleLine()
						ui.Text(c, s.Label).FontSize(t.Rem(0.7)).TextColor(t.AccentText.Alpha(0.8)).FillWidth().SingleLine()
					})
				}
			})
		}
	})
	return root
}
