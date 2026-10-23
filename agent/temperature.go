package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TemperatureOptions name the ends.
type TemperatureOptions struct {
	Precise, Creative string
}

// Temperature is a slider from Precise to Creative.
func Temperature(c *ui.Context, value *float64, o TemperatureOptions) ui.Element {
	if value == nil || !(*value >= 0 && *value <= 1) {
		panic(fmt.Sprintf("agent: Temperature needs a value from 0 to 1, got %v", *value))
	}
	def(&o.Precise, "Precise")
	def(&o.Creative, "Creative")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Role(ui.RoleGroup).Label("Temperature")
	return root.Children(func() {
		kit.Muted(c, o.Precise).FontSize(t.Rem(0.8)).Shrink(0)
		s := ui.Slider(c, value, 0, 1).Grow(1)
		released(s, "Temperature", *value)
		if s.Changed() {
			root.SetChanged()
		}
		kit.Muted(c, o.Creative).FontSize(t.Rem(0.8)).Shrink(0)
	})
}
