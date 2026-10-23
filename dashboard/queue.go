package dashboard

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// QueueOptions are the queue labels.
type QueueOptions struct{ Waiting, Running string }

// Queue shows jobs waiting and running.
func Queue(c *ui.Context, waiting, running int, o QueueOptions) ui.Element {
	if waiting < 0 || running < 0 {
		panic("dashboard: Queue counts must not be negative")
	}
	def(&o.Waiting, "Waiting")
	def(&o.Running, "Running")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).Role(ui.RoleGroup)
	root.Children(func() {
		for _, row := range []struct {
			name  string
			value int
			ink   ui.Color
		}{{o.Waiting, waiting, t.Warning}, {o.Running, running, t.Accent}} {
			ui.Column(c).Grow(1).Basis(0).Gap(t.Space(1)).Children(func() {
				kit.Muted(c, row.name).FontSize(t.Rem(0.75))
				kit.Numeric(c, fmt.Sprint(row.value)).Bold().TextColor(row.ink)
			})
		}
	})
	return root
}
