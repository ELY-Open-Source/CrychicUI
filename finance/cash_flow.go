package finance

import (
	"fmt"
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FlowMonth is one month's incoming and outgoing cash.
type FlowMonth struct {
	Label   string
	In, Out float64
}

// CashFlowOptions word the card; Net takes the signed net.
type CashFlowOptions struct{ Title, In, Out, Net string }

// CashFlow charts monthly in and out; hover shows a month.
func CashFlow(c *ui.Context, months []FlowMonth, currency string, o CashFlowOptions) ui.Element {
	if len(months) == 0 || len(months) > 24 || currency == "" {
		panic("finance: CashFlow needs 1–24 months and a currency")
	}
	peak, in, out := 0.0, 0.0, 0.0
	for i, month := range months {
		if month.Label == "" || month.In < 0 || month.Out < 0 || !finite(month.In) || !finite(month.Out) ||
			slices.IndexFunc(months, func(m FlowMonth) bool { return m.Label == month.Label }) != i {
			panic(fmt.Sprintf("finance: CashFlow month %q needs a unique label and amounts of zero or more", month.Label))
		}
		peak = max(peak, month.In, month.Out)
		in += month.In
		out += month.Out
	}
	if !finite(in) || !finite(out) || !finite(in-out) {
		panic("finance: CashFlow totals overflow")
	}
	if peak == 0 {
		peak = 1
	}
	def(&o.Title, "Cash flow")
	def(&o.In, "In")
	def(&o.Out, "Out")
	def(&o.Net, "Net %s")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "arrow-up-down", o.Title, func() {
			net := in - out
			sign := ""
			if net > 0 {
				sign = "+"
			}
			kit.Chip(c, fmt.Sprintf(o.Net, sign+moneyText(c, net, currency)), changeInk(t, net))
		})
		ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Children(func() {
			hover := -1
			chart := ui.Row(c).FillWidth().Height(t.Space(34)).Gap(t.Space(1))
			grow := chart.Keyframes("crychic.grow", []ui.Keyframe{{Value: 0}, {At: 800 * time.Millisecond, Value: 1, Ease: ui.EaseOut}})
			chart.Children(func() {
				for i, month := range months {
					c.Key(month.Label)
					col := ui.Column(c).Grow(1).Basis(0).MinWidth(0).FillHeight().Gap(t.Space(1)).Padding(t.Space(1), 0).Radius(t.Radius)
					if col.Hovered() {
						hover = i
					}
					col.Background(col.AnimateColor("bg", t.SurfaceHover.Alpha(unit(col.Hovered())), ui.Snappy))
					col.Children(func() {
						bars := ui.Box(c).FillWidth().Grow(1)
						bars.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
							w := min(t.Space(3), r.W/3)
							x := r.X + r.W/2 - w - t.Space(0.25)
							for k, v := range []float64{month.In, month.Out} {
								h := max(t.Space(0.5), r.H*float32(v/peak)*grow)
								ink := []ui.Color{t.Success, t.Danger}[k]
								p.Fill(ui.Rect{X: x + float32(k)*(w+t.Space(0.5)), Y: r.Y + r.H - h, W: w, H: h}, ink, t.Space(0.75))
							}
						}))
						kit.Muted(c, month.Label).FontSize(t.Rem(0.7)).FillWidth().TextAlign(ui.Center).SingleLine()
					})
				}
			})
			shownIn, shownOut := in, out
			if hover >= 0 {
				shownIn, shownOut = months[hover].In, months[hover].Out
			}
			ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
				for k, v := range []float64{shownIn, shownOut} {
					c.Key(k)
					ui.Row(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).
						Radius(t.Radius).Background(t.Surface).Children(func() {
						ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background([]ui.Color{t.Success, t.Danger}[k])
						kit.Muted(c, []string{o.In, o.Out}[k]).FontSize(t.Rem(0.75)).Shrink(0)
						money(c, v, currency).Bold().FontSize(t.Rem(0.85)).Grow(1).MinWidth(0).TextAlign(ui.End).SingleLine().
							ContentTransition(ui.ContentTransition{Numeric: true})
					})
				}
			})
		})
	})
	return root
}
