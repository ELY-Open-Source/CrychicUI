package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ServerLoadOptions are the server load labels.
type ServerLoadOptions struct{ CPU, Memory string }

// ServerLoad shows CPU and memory bars that move.
func ServerLoad(c *ui.Context, cpu, memory float64, o ServerLoadOptions) ui.Element {
	if cpu < 0 || cpu > 100 || memory < 0 || memory > 100 || !finite(cpu) || !finite(memory) {
		panic("dashboard: ServerLoad has invalid values")
	}
	def(&o.CPU, "CPU")
	def(&o.Memory, "Memory")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, row := range []struct {
			name  string
			value float64
		}{{o.CPU, cpu}, {o.Memory, memory}} {
			ui.Column(c).FillWidth().Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Children(func() {
					ui.Text(c, row.name).Grow(1)
					kit.Numeric(c, fmt.Sprintf("%.0f%%", row.value)).TextColor(t.TextMuted)
				})
				progress(c, row.value/100, t.Accent)
			})
		}
	})
	return root
}
