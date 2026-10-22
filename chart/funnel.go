package chart

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FunnelOptions tune a Funnel.
type FunnelOptions struct {
	// Format prints a value; Rate takes a stage's percent kept.
	Format, Rate string
}

// Funnel narrows stages by value, each with its rate kept.
func Funnel(c *ui.Context, label string, names []string, vs []float64, at *int, o FunnelOptions) ui.Element {
	values("Funnel", label, names, vs)
	for _, v := range vs {
		if v < 0 {
			panic(fmt.Sprintf("chart: Funnel %q values must not be negative", label))
		}
	}
	if vs[0] == 0 {
		panic(fmt.Sprintf("chart: Funnel %q first stage must be above zero", label))
	}
	chosen("Funnel", label, at, len(names))
	def(&o.Format, "%g")
	def(&o.Rate, "%.0f%%")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1.5)).Role(ui.RoleList).Label(label)
	g := grow(root)
	return root.Children(func() {
		for i, v := range vs {
			c.Key(names[i])
			rate := ""
			if i > 0 && vs[i-1] > 0 {
				rate = fmt.Sprintf(o.Rate, 100*v/vs[i-1])
			}
			tip := fmt.Sprintf("%s: "+o.Format, names[i], v)
			ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Grow(1).Center().Children(func() {
					b := ui.Box(c).Height(t.Space(9)).Radius(t.Radius).Center().Focusable().Label(tip).Tooltip(tip)
					ink := tone(t, 0).Mix(t.Background, 0.12*float32(i))
					if *at >= 0 && *at != i {
						ink = ink.Alpha(0.4)
					}
					if b.Hovered() {
						ink = ink.Mix(t.Text, 0.15)
					}
					b.Background(b.AnimateColor("ink", ink, ui.Snappy)).
						WidthPercent(100 * g * b.AnimateTo("w", float32(v/vs[0]), ui.Snappy)).Children(func() {
						ui.Text(c, names[i]).FontSize(t.Rem(0.8)).Bold().TextColor(t.AccentText).NoWrap()
					})
					if b.Clicked() {
						pick(root, at, i, "Funnel")
					}
				})
				kit.Numeric(c, rate).TextColor(t.TextMuted).Width(t.Space(12)).Shrink(0).TextAlign(ui.End)
			})
		}
	})
}
