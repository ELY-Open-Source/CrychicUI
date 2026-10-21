package input

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UnitFieldOptions tune a UnitField.
type UnitFieldOptions struct {
	// Min, Max and Step bound the number.
	Min, Max, Step float64
}

// UnitField pairs a number with a unit select.
func UnitField(c *ui.Context, label string, value *float64, unit *string, units []string, o UnitFieldOptions) ui.Element {
	if value == nil || !finite(*value) || unit == nil || !slices.Contains(units, *unit) {
		panic("input: UnitField needs a finite value and selected unit")
	}
	if o.Step == 0 {
		o.Step = 1
	}
	if o.Max == 0 {
		o.Max = 1000000
	}
	if len(units) == 0 || !finite(o.Min) || !finite(o.Max) || !finite(o.Step) || o.Step <= 0 || o.Min > o.Max {
		panic("input: UnitField needs units and a valid range")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1))
	return root.Children(func() {
		ui.Text(c, label).FontWeight(500)
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Radius(t.Radius).Background(t.Background).Border(1, t.Border).Children(func() {
			in := ui.NumberInput(c, value, o.Min, o.Max, o.Step).Grow(1).Padding(t.Space(2), 0, t.Space(2), t.Space(3))
			changed(root, in)
			menu := ui.MenuButton(c, "", func(m *ui.Menu) {
				for _, choice := range units {
					if m.Item(choice).Checked(choice == *unit).Chosen() {
						*unit = choice
						root.SetChanged()
						sound.Play(sound.Select)
						kit.Log("UnitField", "unit")
					}
				}
			})
			menu.Children(func() { ui.Text(c, *unit).ContentTransition(ui.ContentTransition{}) })
		})
	})
}
