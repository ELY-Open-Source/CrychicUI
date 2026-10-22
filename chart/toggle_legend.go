package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
	mc "github.com/egoist/mygo/ui/chart"
)

// ToggleLegendOptions tune a ToggleLegend.
type ToggleLegendOptions = LineOptions

// ToggleLegend draws lines whose legend hides series by hidden.
func ToggleLegend(c *ui.Context, label string, labels []string, series []Series, hidden *[]bool, at *int, o ToggleLegendOptions) ui.Element {
	if hidden == nil || len(*hidden) != len(series) {
		panic(fmt.Sprintf("chart: ToggleLegend %q needs a hidden flag per series", label))
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2))
	return root.Children(func() {
		off := *hidden
		chart := lines(c, label, labels, series, at, look{curve: mc.CurveMonotone, h: o.Height, format: o.Format,
			comp: "ToggleLegend", off: off})
		if chart.Changed() {
			root.SetChanged()
		}
		ui.Row(c).Wrap().Gap(t.Space(2)).Children(func() {
			for s, sr := range series {
				c.Key(s)
				b := kit.Press(ui.ButtonBase(c).Label(sr.Name).Checked(!off[s]).Gap(t.Space(1.5)).AlignItems(ui.Center).
					Padding(t.Space(1), t.Space(2.5)).Radius(999).Border(1, t.Border))
				dim := 1 - 0.55*b.AnimateTo("off", unit(off[s]), ui.Snappy)
				b.Opacity(dim).Children(func() {
					ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Radius(999).Border(2, tone(t, s)).
						Background(tone(t, s).Alpha(1 - unit(off[s])))
					ui.Text(c, sr.Name).FontSize(t.Rem(0.8))
				})
				if kit.Clicked(b, sound.Toggle) {
					(*hidden)[s] = !off[s]
					root.SetChanged()
					kit.Log("ToggleLegend", "toggled", "series", s, "hidden", (*hidden)[s])
				}
			}
		})
	})
}
