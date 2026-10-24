package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TemperatureOptions name the slider, its ends and its middle.
type TemperatureOptions struct {
	Label, Precise, Balanced, Creative string
}

// Temperature is a slider from Precise to Creative.
func Temperature(c *ui.Context, value *float64, o TemperatureOptions) ui.Element {
	if value == nil || !(*value >= 0 && *value <= 1) {
		panic("agent: Temperature needs a value from 0 to 1")
	}
	def(&o.Label, "Temperature")
	def(&o.Precise, "Precise")
	def(&o.Balanced, "Balanced")
	def(&o.Creative, "Creative")
	t := c.Theme()
	mood := o.Balanced
	switch {
	case *value < 0.34:
		mood = o.Precise
	case *value > 0.66:
		mood = o.Creative
	}
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Background(t.Background).
		Border(1, t.Border).Role(ui.RoleGroup).Label(o.Label)
	return root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, o.Label).Bold().Grow(1).SingleLine()
			kit.Muted(c, mood).FontSize(t.Rem(0.8)).ContentTransition(ui.ContentTransition{})
			kit.Numeric(c, fmt.Sprintf("%.2f", *value)).Bold().TextColor(t.Accent).Padding(t.Space(0.5), t.Space(2)).
				Radius(t.Radius).Background(t.Accent.Alpha(0.12))
		})
		s := ui.Slider(c, value, 0, 1).Label(o.Label).FillWidth()
		released(s, "Temperature", *value)
		if s.Changed() {
			root.SetChanged()
		}
		ui.Row(c).FillWidth().Children(func() {
			kit.Muted(c, o.Precise).FontSize(t.Rem(0.75)).Grow(1)
			kit.Muted(c, o.Creative).FontSize(t.Rem(0.75))
		})
	})
}
