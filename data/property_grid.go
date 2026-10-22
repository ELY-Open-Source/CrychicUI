package data

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Property is a named *string, *bool or *float64.
type Property struct {
	Name  string
	Value any
	// Min, Max and Step bound a number.
	Min, Max, Step float64
}

// PropertyGrid edits typed properties, one row each.
func PropertyGrid(c *ui.Context, label string, props []Property) ui.Element {
	if label == "" || len(props) == 0 {
		panic(fmt.Sprintf("data: PropertyGrid %q needs properties", label))
	}
	var names []string
	for _, p := range props {
		names = append(names, p.Name)
		if v, ok := p.Value.(*float64); ok {
			bad := !finite64(p.Min) || !finite64(p.Max) || !finite64(p.Step) || !finite64(*v) || p.Max <= p.Min || p.Step < 0
			if bad || *v < p.Min || *v > p.Max {
				panic(fmt.Sprintf("data: PropertyGrid %q number %q needs finite Min < Max, Step and value", label, p.Name))
			}
		}
	}
	unique("PropertyGrid", label, names)
	t := c.Theme()
	root := frame(c, label)
	root.Children(func() {
		for i, p := range props {
			c.Key(i)
			row := hovered(c, false).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(3))
			kit.Arrive(row, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(i)*30*time.Millisecond).Children(func() {
				kit.Muted(c, p.Name).Width(t.Space(36)).Shrink(0).SingleLine()
				var e ui.Element
				switch v := p.Value.(type) {
				case *string:
					e = ui.TextInput(c, v).Label(p.Name).Grow(1)
				case *bool:
					e = ui.Switch(c, v).Label(p.Name)
				case *float64:
					step := p.Step
					if step == 0 {
						step = 1
					}
					e = ui.NumberInput(c, v, p.Min, p.Max, step).Label(p.Name).Grow(1)
				default:
					panic(fmt.Sprintf("data: PropertyGrid %q cannot edit %T", label, p.Value))
				}
				if e.Changed() {
					root.SetChanged()
					sound.Play(sound.Tick)
					kit.Log("PropertyGrid", "edited", "index", i)
				}
			})
		}
	})
	return root
}
