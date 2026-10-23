package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DepthOptions are the depth labels.
type DepthOptions struct{ Bids, Asks string }

// Depth shows cumulative bid and ask areas.
func Depth(c *ui.Context, bids, asks []BookLevel, currency string, o DepthOptions) ui.Element {
	if len(bids) == 0 || len(asks) == 0 || currency == "" {
		panic("finance: Depth needs both sides and currency")
	}
	def(&o.Bids, "Bids")
	def(&o.Asks, "Asks")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, side := range []struct {
			name string
			rows []BookLevel
			ink  ui.Color
		}{{o.Bids, bids, t.Success}, {o.Asks, asks, t.Danger}} {
			ui.Column(c).Grow(1).Basis(0).Gap(t.Space(1)).Children(func() {
				ui.Text(c, side.name).Bold()
				cumulative := make([]float64, len(side.rows))
				total := 0.0
				for i, row := range side.rows {
					if row.Price <= 0 || row.Amount <= 0 || !finite(row.Price) || !finite(row.Amount) {
						panic("finance: Depth has an invalid level")
					}
					total += row.Amount
					if !finite(total) || total <= 0 {
						panic("finance: Depth total overflow")
					}
					cumulative[i] = total
				}
				box := ui.Box(c).FillWidth().Height(t.Space(24))
				kit.Arrive(box, ui.Pose{Fade: 1}, 0)
				box.Draw(kit.Inside(0, func(p *ui.Painter, r ui.Rect) {
					var area ui.Path
					area.MoveTo(r.X, r.Y+r.H)
					for i, value := range cumulative {
						x := r.X + r.W*float32(i)/float32(max(1, len(cumulative)-1))
						y := r.Y + r.H*(1-float32(value/total))
						area.LineTo(x, y)
					}
					area.LineTo(r.X+r.W, r.Y+r.H).Close()
					p.FillPath(&area, side.ink.Alpha(0.25))
					p.StrokePath(&area, 1.5, side.ink)
				}))
			})
		}
	})
	return root
}
