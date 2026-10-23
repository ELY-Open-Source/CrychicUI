package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LatencyOptions are the latency labels.
type LatencyOptions struct{ P50, P95, P99 string }

// Latency shows p50, p95, and p99 values.
func Latency(c *ui.Context, p50, p95, p99 float64, o LatencyOptions) ui.Element {
	if p50 < 0 || p95 < p50 || p99 < p95 || !finite(p50) || !finite(p95) || !finite(p99) {
		panic("dashboard: Latency has invalid values")
	}
	def(&o.P50, "p50")
	def(&o.P95, "p95")
	def(&o.P99, "p99")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, row := range []struct {
			name  string
			value float64
		}{{o.P50, p50}, {o.P95, p95}, {o.P99, p99}} {
			ui.Column(c).Grow(1).Basis(0).Children(func() {
				kit.Muted(c, row.name).FontSize(t.Rem(0.75))
				kit.Numeric(c, fmt.Sprintf("%.0f ms", row.value)).Bold()
			})
		}
	})
	return root
}
