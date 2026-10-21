package layout

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PageHeaderOptions tune a PageHeader.
type PageHeaderOptions struct {
	// Subtitle shows under the title when set.
	Subtitle string
	// Actions fill the header's end when set.
	Actions func()
}

// PageHeader shows a title, subtitle and actions; title slides in.
func PageHeader(c *ui.Context, title string, o PageHeaderOptions) ui.Element {
	t := c.Theme()
	return ui.Row(c).FillWidth().Gap(t.Space(4)).Padding(t.Space(4), 0).Children(func() {
		ui.Column(c).Grow(1).Basis(0).Gap(t.Space(1)).Children(func() {
			kit.Arrive(ui.Text(c, title).Bold().FontSize(t.Rem(1.6)), ui.Pose{X: -8, Fade: 1}, 0)
			if o.Subtitle != "" {
				kit.Arrive(kit.Muted(c, o.Subtitle), ui.Pose{X: -8, Fade: 1}, 60*time.Millisecond)
			}
		})
		if o.Actions != nil {
			ui.Row(c).Gap(t.Space(2)).Children(o.Actions)
		}
	})
}
