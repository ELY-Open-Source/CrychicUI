package project

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// StandUpPerson is one stand-up report.
type StandUpPerson struct {
	Name, Yesterday, Today, Blockers string
}

// StandUpOptions are the stand-up labels.
type StandUpOptions struct{ Yesterday, Today, Blockers string }

// StandUp shows yesterday, today, and blockers per person.
func StandUp(c *ui.Context, people []StandUpPerson, o StandUpOptions) ui.Element {
	if len(people) == 0 {
		panic("project: StandUp needs people")
	}
	def(&o.Yesterday, "Yesterday")
	def(&o.Today, "Today")
	def(&o.Blockers, "Blockers")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		for i, person := range people {
			if person.Name == "" || person.Yesterday == "" || person.Today == "" {
				panic("project: StandUp has an invalid person")
			}
			c.Key(i)
			ui.Column(c).Gap(t.Space(1)).Padding(t.Space(2)).Radius(t.Radius).Border(1, t.Border).Children(func() {
				ui.Text(c, person.Name).Bold()
				ui.Text(c, o.Yesterday+": "+person.Yesterday).FontSize(t.Rem(0.8))
				ui.Text(c, o.Today+": "+person.Today).FontSize(t.Rem(0.8))
				if person.Blockers != "" {
					ui.Text(c, o.Blockers+": "+person.Blockers).TextColor(t.Danger).FontSize(t.Rem(0.8))
				}
			})
		}
	})
	return root
}
