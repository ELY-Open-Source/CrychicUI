package schedule

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Attendee is one attendee and response.
type Attendee struct {
	Name     string
	Response RSVPResponse
}

// AttendeesOptions word the list; Going takes a count.
type AttendeesOptions struct{ Title, Going, Yes, Maybe, No, Waiting string }

// Attendees shows avatars with responses.
func Attendees(c *ui.Context, people []Attendee, o AttendeesOptions) ui.Element {
	if len(people) == 0 {
		panic("schedule: Attendees needs people")
	}
	def(&o.Title, "Attendees")
	def(&o.Going, "%d going")
	def(&o.Yes, "Going")
	def(&o.Maybe, "Maybe")
	def(&o.No, "Declined")
	def(&o.Waiting, "Awaiting")
	t := c.Theme()
	going := 0
	for _, person := range people {
		if person.Name == "" || person.Response > No {
			panic(fmt.Sprintf("schedule: Attendees has an invalid person %q", person.Name))
		}
		if person.Response == Yes {
			going++
		}
	}
	root := card(c)
	root.Children(func() {
		kit.CardHead(c, "users", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Going, going), t.Success) })
		list := ui.Column(c).FillWidth().Padding(t.Space(1), 0)
		list.Children(func() {
			for i, person := range people {
				name, text, ink := "clock", o.Waiting, t.TextMuted
				switch person.Response {
				case Yes:
					name, text, ink = "circle-check", o.Yes, t.Success
				case Maybe:
					name, text, ink = "circle-help", o.Maybe, t.Warning
				case No:
					name, text, ink = "circle-x", o.No, t.Danger
				}
				c.Key(i)
				row := ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(3))
				row.Background(row.AnimateColor("bg", t.SurfaceHover.Alpha(0.6*unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					display.Avatar(c, person.Name, nil).Size(t.Space(8), t.Space(8))
					ui.Text(c, person.Name).FontSize(t.Rem(0.85)).Bold().Grow(1).MinWidth(0).SingleLine()
					ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(2)).Radius(999).
						Background(ink.Alpha(0.12)).TextColor(ink).Children(func() {
						glyph(c, name, 12)
						ui.Text(c, text).FontSize(t.Rem(0.7)).Bold().Shrink(0)
					})
				})
			}
		})
	})
	return root
}
