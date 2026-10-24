package schedule

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// TimeOffStatus is where a leave request stands.
type TimeOffStatus uint8

const (
	Pending TimeOffStatus = iota
	Approved
	Declined
)

// TimeOffRequest is one leave request; Kind is optional.
type TimeOffRequest struct {
	Name, Kind, Range string
	Status            TimeOffStatus
}

// TimeOffOptions word the list; Waiting takes a count.
type TimeOffOptions struct{ Title, Waiting, Pending, Approved, Declined string }

// TimeOff shows requests with status.
func TimeOff(c *ui.Context, requests []TimeOffRequest, o TimeOffOptions) ui.Element {
	if len(requests) == 0 {
		panic("schedule: TimeOff needs requests")
	}
	def(&o.Title, "Time off")
	def(&o.Waiting, "%d pending")
	def(&o.Pending, "Pending")
	def(&o.Approved, "Approved")
	def(&o.Declined, "Declined")
	t := c.Theme()
	waiting := 0
	for _, request := range requests {
		if request.Name == "" || request.Range == "" || request.Status > Declined {
			panic(fmt.Sprintf("schedule: TimeOff has an invalid request for %q", request.Name))
		}
		if request.Status == Pending {
			waiting++
		}
	}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "plane", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Waiting, waiting), t.Warning) })
		list := ui.Column(c).FillWidth().Padding(t.Space(1), 0)
		list.Children(func() {
			for i, request := range requests {
				name, text, ink := "clock", o.Pending, t.Warning
				switch request.Status {
				case Approved:
					name, text, ink = "circle-check", o.Approved, t.Success
				case Declined:
					name, text, ink = "circle-x", o.Declined, t.Danger
				}
				line := request.Range
				if request.Kind != "" {
					line = request.Kind + " · " + request.Range
				}
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(3))
				row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(0.6*unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					display.Avatar(c, request.Name, nil).Size(t.Space(8), t.Space(8))
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						ui.Text(c, request.Name).FontSize(t.Rem(0.85)).Bold().SingleLine().MinWidth(0)
						kit.Muted(c, line).FontSize(t.Rem(0.7)).SingleLine().MinWidth(0)
					})
					chip := ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999)
					chip.Background(chip.AnimateColor("bg", ink.Alpha(0.12), ui.Snappy)).TextColor(chip.AnimateColor("ink", ink, ui.Snappy))
					chip.Children(func() {
						c.Key(name)
						kit.Arrive(glyph(c, name, 12), ui.Pose{Scale: 0.5, Fade: 1}, 0)
						ui.Text(c, text).FontSize(t.Rem(0.7)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
					})
				})
			}
		})
	})
	return root
}
