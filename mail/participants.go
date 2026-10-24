package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Participants lists call participants with mute marks.
func Participants(c *ui.Context, participants []Participant) ui.Element {
	if len(participants) == 0 {
		panic("mail: Participants needs people")
	}
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleList)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for _, person := range participants {
			if person.ID == "" || person.Name == "" {
				panic("mail: invalid participant")
			}
			c.Key(person.ID)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1)).Children(func() {
				ui.Box(c).Size(t.Space(6), t.Space(6)).Radius(999).Background(t.Accent.Alpha(0.15)).Center().Children(func() {
					ui.Text(c, initial(person.Name))
				})
				ui.Text(c, person.Name).Grow(1)
				if person.Muted {
					ui.Icon(c, icon.Must("mic-off")).Size(16, 16).TextColor(t.Danger)
				}
			})
		}
	})
	return root
}
