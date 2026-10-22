package form

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// RequiredLegendOptions tune a RequiredLegend.
type RequiredLegendOptions struct {
	// Text follows the mark; "Required" when empty.
	Text string
}

// RequiredLegend explains the required mark; it fades in.
func RequiredLegend(c *ui.Context, o RequiredLegendOptions) ui.Element {
	def(&o.Text, "Required")
	t := c.Theme()
	row := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Label("* " + o.Text)
	return kit.Arrive(row, ui.Pose{Fade: 1}, 0).Children(func() {
		ui.Text(c, "*").TextColor(t.Danger).FontWeight(600)
		kit.Muted(c, o.Text)
	})
}
