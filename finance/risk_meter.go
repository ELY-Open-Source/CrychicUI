package finance

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RiskMeterOptions are the risk labels.
type RiskMeterOptions struct{ Low, High string }

// RiskMeter shows risk one through five on a scale.
func RiskMeter(c *ui.Context, level int, o RiskMeterOptions) ui.Element {
	if level < 1 || level > 5 {
		panic("finance: RiskMeter level must be 1–5")
	}
	def(&o.Low, "Low risk")
	def(&o.High, "High risk")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(1)).Children(func() {
			for i := 1; i <= 5; i++ {
				ink := t.Success
				if i >= 4 {
					ink = t.Danger
				} else if i == 3 {
					ink = t.Warning
				}
				segment := ui.Box(c).Grow(1).Basis(0).Height(t.Space(2)).Radius(999)
				bg := t.SurfaceHover
				if i <= level {
					bg = ink
				}
				segment.Background(segment.AnimateColor("risk", bg, ui.Snappy))
			}
		})
		ui.Row(c).FillWidth().Children(func() {
			kit.Muted(c, o.Low).FontSize(t.Rem(0.75)).Grow(1)
			kit.Muted(c, o.High).FontSize(t.Rem(0.75))
		})
	})
	return root
}
