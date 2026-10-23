package dashboard

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RevenueOptions are the revenue labels.
type RevenueOptions struct{ Title string }

// Revenue shows a chart and period select; cue Select.
func Revenue(c *ui.Context, points []Point, period *string, periods []string, o RevenueOptions) ui.Element {
	if len(points) < 2 || period == nil || len(periods) == 0 {
		panic("dashboard: Revenue needs points and periods")
	}
	def(&o.Title, "Revenue")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.Title).Bold().Grow(1)
			for _, choice := range periods {
				if choice == "" {
					panic("dashboard: Revenue period must not be empty")
				}
				b := kit.Press(ui.ButtonBase(c).Label(choice).Padding(t.Space(0.5), t.Space(1.5)).Radius(999).
					Border(1, t.Border).Checked(*period == choice))
				b.Background(b.AnimateColor("bg", t.Surface.Mix(t.Accent.Alpha(0.12), unit(*period == choice)), ui.Snappy))
				b.Children(func() { ui.Text(c, choice).FontSize(t.Rem(0.75)) })
				if b.Clicked() && *period != choice {
					*period = choice
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Revenue", "period", "period", *period)
				}
			}
		})
		series(c, points, t.Accent)
	})
	return root
}
