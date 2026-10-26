package home

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SensorKind is what a contact sensor guards.
type SensorKind int

const (
	DoorSensor SensorKind = iota
	WindowSensor
)

// Sensor is one contact sensor; Since is shown as given.
type Sensor struct {
	Name, Room, Since string
	Kind              SensorKind
	Open              bool
}

// SensorsOptions word the list; Count takes the open count.
type SensorsOptions struct{ Title, Open, Closed, AllClosed, Count, Empty string }

// Sensors lists door and window sensors with their state.
func Sensors(c *ui.Context, sensors []Sensor, o SensorsOptions) ui.Element {
	def(&o.Title, "Sensors")
	def(&o.Open, "Open")
	def(&o.Closed, "Closed")
	def(&o.AllClosed, "All closed")
	def(&o.Count, "%d open")
	def(&o.Empty, "No sensors yet")
	names := make([]string, len(sensors))
	open := 0
	for i, s := range sensors {
		if s.Kind < DoorSensor || s.Kind > WindowSensor {
			panic(fmt.Sprintf("home: Sensors sensor %q has an unknown kind", s.Name))
		}
		names[i] = s.Name
		if s.Open {
			open++
		}
	}
	unique("Sensors", names)
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	root.Children(func() {
		kit.CardHead(c, "shield-check", o.Title, func() {
			if len(sensors) == 0 {
				return
			}
			if open == 0 {
				kit.Chip(c, o.AllClosed, t.Success)
			} else {
				kit.Chip(c, fmt.Sprintf(o.Count, open), t.Warning)
			}
		})
		list := ui.Column(c).FillWidth().Padding(t.Space(1.5), 0)
		list.Children(func() {
			if len(sensors) == 0 {
				kit.Muted(c, o.Empty).FontSize(t.Rem(0.8)).FillWidth().TextAlign(ui.Center).Padding(t.Space(5), t.Space(3))
				return
			}
			for i, s := range sensors {
				c.Key(strconv.Quote(s.Name))
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Role(ui.RoleListItem)
				kit.Stagger(c, list, row, i).Children(func() {
					mark := "door-closed"
					switch {
					case s.Kind == DoorSensor && s.Open:
						mark = "door-open"
					case s.Kind == WindowSensor:
						mark = "app-window"
					}
					tile(c, mark, t.Warning, t.Space(9), s.Open)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
						ui.Text(c, s.Name).FontSize(t.Rem(0.84)).Bold().FillWidth().MinWidth(0).SingleLine()
						ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
							if s.Room != "" {
								kit.Muted(c, s.Room).FontSize(t.Rem(0.7)).Shrink(0).MaxWidth(t.Space(28)).SingleLine()
							}
							if s.Since != "" {
								kit.Muted(c, s.Since).FontSize(t.Rem(0.7)).Grow(1).MinWidth(0).SingleLine()
							}
						})
					})
					state, ink := o.Closed, t.Success
					if s.Open {
						state, ink = o.Open, t.Warning
					}
					c.Key(state)
					kit.Arrive(kit.Chip(c, state, ink), ui.Pose{Scale: 0.8, Fade: 1}, 0)
				})
			}
		})
	})
	return root
}
