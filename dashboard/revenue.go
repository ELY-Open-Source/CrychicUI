package dashboard

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/ELY-Open-Source/CrychicUI/typography"
	"github.com/egoist/mygo/ui"
)

// RevenueOptions word the card.
type RevenueOptions struct{ Title string }

// Revenue charts a period with its total; cue Select.
func Revenue(c *ui.Context, points []Point, currency string, period *string, periods []string, o RevenueOptions) ui.Element {
	if len(points) < 2 || currency == "" || period == nil || len(periods) == 0 {
		panic("dashboard: Revenue needs two points, a currency and periods")
	}
	unique("Revenue", periods)
	at := slices.Index(periods, *period)
	if at < 0 {
		panic(fmt.Sprintf("dashboard: Revenue period %q is not among %q", *period, periods))
	}
	total := 0.0
	for _, p := range points {
		if p.Value < 0 || !finite(p.Value) {
			panic(fmt.Sprintf("dashboard: Revenue point %q needs a finite value of zero or more", p.Label))
		}
		total += p.Value
	}
	if !finite(total) {
		panic("dashboard: Revenue total overflows")
	}
	def(&o.Title, "Revenue")
	first, last := points[0].Value, points[len(points)-1].Value
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "chart-line", o.Title, nil)
		ui.Row(c).FillWidth().Padding(t.Space(3), t.Space(3), 0, t.Space(3)).Children(func() {
			seg := ui.Row(c).FillWidth().Padding(t.Space(0.5)).Radius(999).Background(t.Surface)
			x := seg.AnimateTo("thumb", float32(at), ui.Snappy)
			seg.Draw(kit.Inside(999, func(p *ui.Painter, r ui.Rect) {
				w := (r.W - t.Space(1)) / float32(len(periods))
				p.Fill(ui.Rect{X: r.X + t.Space(0.5) + x*w, Y: r.Y + t.Space(0.5), W: w, H: r.H - t.Space(1)}, t.Background, 999)
			}))
			seg.Children(func() {
				for i, choice := range periods {
					c.Key(choice)
					on := i == at
					b := kit.Press(ui.ButtonBase(c).Label(choice).Grow(1).Basis(0).MinWidth(0).Padding(t.Space(0.75), t.Space(1)).
						Radius(999).Checked(on))
					b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Text, unit(on || b.Hovered())), ui.Snappy))
					b.Children(func() {
						ui.Text(c, choice).FontSize(t.Rem(0.75)).Bold().Grow(1).MinWidth(0).TextAlign(ui.Center).SingleLine()
					})
					if b.Clicked() && !on {
						*period = choice
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("Revenue", "period", "index", i)
					}
				}
			})
		})
		ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Children(func() {
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				typography.Money(c, total, currency).Bold().FontSize(t.Rem(1.6)).Grow(1).MinWidth(0).
					ContentTransition(ui.ContentTransition{Numeric: true})
				if first > 0 {
					change := 100 * ((last - first) / first)
					if !finite(change) {
						panic("dashboard: Revenue change overflows")
					}
					deltaChip(c, fmt.Sprintf("%+.1f%%", change))
				}
			})
			series(c, "Revenue", points, t.Accent).Height(t.Space(30))
			ui.Row(c).FillWidth().Children(func() {
				kit.Muted(c, points[0].Label).FontSize(t.Rem(0.7)).Grow(1).ContentTransition(ui.ContentTransition{})
				kit.Muted(c, points[len(points)-1].Label).FontSize(t.Rem(0.7)).Shrink(0).ContentTransition(ui.ContentTransition{})
			})
		})
	})
	return root
}
