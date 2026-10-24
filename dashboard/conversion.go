package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FunnelStage is one conversion stage.
type FunnelStage struct {
	Label string
	Count int
}

// ConversionOptions word the card; Overall takes the end rate.
type ConversionOptions struct{ Title, Overall string }

// Conversion shows a funnel with each step's rate.
func Conversion(c *ui.Context, stages []FunnelStage, o ConversionOptions) ui.Element {
	if len(stages) < 2 || stages[0].Count <= 0 {
		panic("dashboard: Conversion needs two stages and a positive first count")
	}
	names := make([]string, len(stages))
	for i, stage := range stages {
		if stage.Count < 0 || i > 0 && stage.Count > stages[i-1].Count {
			panic(fmt.Sprintf("dashboard: Conversion stage %q needs a count no larger than the one before", stage.Label))
		}
		names[i] = stage.Label
	}
	unique("Conversion", names)
	def(&o.Title, "Conversion")
	def(&o.Overall, "%s overall")
	first := float64(stages[0].Count)
	rate := func(a, b int) string {
		if b == 0 {
			return "0%"
		}
		return fmt.Sprintf("%.0f%%", 100*(float64(a)/float64(b)))
	}
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "filter", o.Title, func() {
			kit.Chip(c, fmt.Sprintf(o.Overall, rate(stages[len(stages)-1].Count, stages[0].Count)), t.Accent)
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(3))
		list.Children(func() {
			for i, stage := range stages {
				c.Key(stage.Label)
				row := ui.Column(c).FillWidth().Gap(t.Space(1))
				kit.Stagger(c, list, row, i).Children(func() {
					if i > 0 {
						ui.Row(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Justify(ui.Center).TextColor(t.TextMuted).
							Children(func() {
								ui.Icon(c, icon.Must("arrow-down")).Size(12, 12)
								kit.Numeric(c, rate(stage.Count, stages[i-1].Count)).FontSize(t.Rem(0.7)).Bold().Shrink(0)
							})
					}
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						ui.Text(c, stage.Label).Bold().FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).SingleLine()
						kit.Numeric(c, number(c, stage.Count)).FontSize(t.Rem(0.85)).Shrink(0)
					})
					bar := ui.Box(c).FillWidth().Height(t.Space(6))
					w := bar.AnimateTo("w", float32(float64(stage.Count)/first), ui.Smooth)
					ink := t.Accent.Alpha(1 - 0.5*float32(i)/float32(len(stages)-1))
					bar.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
						bw := max(t.Space(1), r.W*w)
						p.Fill(ui.Rect{X: r.X + (r.W-bw)/2, Y: r.Y, W: bw, H: r.H}, ink, t.Space(1))
					}))
				})
			}
		})
	})
	return root
}
