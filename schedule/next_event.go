package schedule

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// NextEventOptions are the next event labels.
type NextEventOptions struct{ Label string }

// NextEvent shows the next event with countdown.
func NextEvent(c *ui.Context, event Event, o NextEventOptions) ui.Element {
	validEvent("NextEvent", event)
	def(&o.Label, "Next event")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup)
	left := max(0, event.Start.Sub(c.Now()))
	if left > 0 {
		c.After(min(time.Second, left))
	}
	root.Children(func() {
		kit.Muted(c, o.Label).FontSize(t.Rem(0.75))
		ui.Text(c, event.Title).Bold()
		kit.Numeric(c, left.Round(time.Second).String()).TextColor(t.Accent)
	})
	return root
}
