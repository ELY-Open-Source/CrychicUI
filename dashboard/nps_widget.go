package dashboard

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// NPSWidgetOptions word the card; Score takes the NPS.
type NPSWidgetOptions struct{ Title, Score, Promoters, Passives, Detractors string }

// NPSWidget shows the score on a scale over three groups.
func NPSWidget(c *ui.Context, promoters, passives, detractors int, o NPSWidgetOptions) ui.Element {
	if promoters < 0 || passives < 0 || detractors < 0 {
		panic("dashboard: NPSWidget needs counts of zero or more")
	}
	total := float64(promoters) + float64(passives) + float64(detractors)
	if total == 0 {
		panic("dashboard: NPSWidget needs one response")
	}
	def(&o.Title, "Net Promoter Score")
	def(&o.Score, "NPS %d")
	def(&o.Promoters, "Promoters")
	def(&o.Passives, "Passives")
	def(&o.Detractors, "Detractors")
	score := int(math.Round(100 * (float64(promoters) - float64(detractors)) / total))
	t := c.Theme()
	ink := t.Danger
	switch {
	case score >= 50:
		ink = t.Success
	case score >= 0:
		ink = t.Warning
	}
	groups := []struct {
		name  string
		count int
		ink   ui.Color
	}{{o.Promoters, promoters, t.Success}, {o.Passives, passives, t.Warning}, {o.Detractors, detractors, t.Danger}}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "smile", o.Title, nil)
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			kit.Numeric(c, fmt.Sprintf(o.Score, score)).Bold().FontSize(t.Rem(1.8)).TextColor(ink).Shrink(0).
				ContentTransition(ui.ContentTransition{Numeric: true})
			scale := ui.Box(c).FillWidth().Height(t.Space(4))
			at := scale.AnimateTo("score", float32(score+100)/200, ui.Smooth)
			scale.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
				y := r.Y + r.H/2
				seg := []ui.Color{t.Danger, t.Warning, t.Success}
				for i, col := range seg {
					w := r.W / 3
					p.Fill(ui.Rect{X: r.X + w*float32(i) + t.Space(0.25), Y: y - t.Space(0.5), W: w - t.Space(0.5), H: t.Space(1)}, col.Alpha(0.35), 999)
				}
				x, d := r.X+t.Space(1.5)+(r.W-t.Space(3))*at, t.Space(3.5)
				p.Fill(ui.Rect{X: x - d/2, Y: y - d/2, W: d, H: d}, t.Background, 999)
				p.Fill(ui.Rect{X: x - d/3, Y: y - d/3, W: 2 * d / 3, H: 2 * d / 3}, ink, 999)
			}))
			list := ui.Column(c).FillWidth().Gap(t.Space(2))
			list.Children(func() {
				for i, g := range groups {
					c.Key(i)
					share := float64(g.count) / total
					line := ui.Column(c).FillWidth().Gap(t.Space(1))
					kit.Stagger(c, list, line, i).Children(func() {
						ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
							ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background(g.ink)
							ui.Text(c, g.name).FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
							kit.Numeric(c, number(c, g.count)).FontSize(t.Rem(0.8)).TextColor(t.TextMuted).Shrink(0)
							kit.Numeric(c, fmt.Sprintf("%.0f%%", 100*share)).Bold().FontSize(t.Rem(0.85)).Width(t.Space(10)).TextAlign(ui.End).Shrink(0)
						})
						progress(c, share, g.ink)
					})
				}
			})
		})
	})
	return root
}
