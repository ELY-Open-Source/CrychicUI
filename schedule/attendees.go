package schedule

import (
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Attendee is one attendee and response.
type Attendee struct {
	Name, Response string
}

// Attendees shows avatars with responses.
func Attendees(c *ui.Context, people []Attendee) ui.Element {
	if len(people) == 0 {
		panic("schedule: Attendees needs people")
	}
	t := c.Theme()
	root := ui.Row(c).Wrap().Gap(t.Space(2)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, person := range people {
			if person.Name == "" || person.Response == "" {
				panic("schedule: Attendees has an invalid person")
			}
			c.Key(i)
			ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
				ui.Box(c).Size(t.Space(7), t.Space(7)).Radius(999).Background(t.Accent.Alpha(0.15)).Center().Children(func() {
					first, _ := utf8.DecodeRuneInString(person.Name)
					ui.Text(c, string(first)).TextColor(t.Accent)
				})
				kit.Muted(c, person.Response).FontSize(t.Rem(0.75))
			})
		}
	})
	return root
}
