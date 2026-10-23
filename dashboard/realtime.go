package dashboard

import (
	"fmt"
	"math"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RealtimeOptions are the realtime labels.
type RealtimeOptions struct{ Label string }

// Realtime shows a live visitor count that ticks.
func Realtime(c *ui.Context, count *int, o RealtimeOptions) ui.Element {
	if count == nil || *count < 0 {
		panic("dashboard: Realtime needs a count")
	}
	def(&o.Label, "Visitors now")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	last := ui.Local(root, "last", func() time.Time { return c.Now() })
	if c.Now().Sub(*last) >= time.Second {
		if *count < math.MaxInt {
			*count++
			root.SetChanged()
			kit.Log("Realtime", "ticked", "count", *count)
		}
		*last = c.Now()
	}
	c.After(max(0, time.Second-c.Now().Sub(*last)))
	root.Children(func() {
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
			ui.Box(c).Size(t.Space(2), t.Space(2)).Radius(999).Background(t.Success)
			kit.Muted(c, o.Label).FontSize(t.Rem(0.75))
		})
		kit.Numeric(c, fmt.Sprint(*count)).Bold().FontSize(t.Rem(1.5))
	})
	return root
}
