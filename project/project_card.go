package project

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ProjectCardOptions word the card; Description and Due are optional.
type ProjectCardOptions struct{ Members, One, Icon, Description, Due string }

// ProjectCard shows name, members, and progress.
func ProjectCard(c *ui.Context, name string, members []string, value float64, o ProjectCardOptions) ui.Element {
	if name == "" || len(members) == 0 || value < 0 || value > 1 || !finite(value) {
		panic(fmt.Sprintf("project: ProjectCard %q needs members and progress in [0,1]", name))
	}
	unique("ProjectCard", members)
	def(&o.Members, "%d members")
	def(&o.One, "1 member")
	def(&o.Icon, "folder-kanban")
	t := c.Theme()
	ink := t.Accent
	if value == 1 {
		ink = t.Success
	}
	count := o.One
	if len(members) > 1 {
		count = fmt.Sprintf(o.Members, len(members))
	}
	root := kit.Lift(card(c).Gap(t.Space(3)).Padding(t.Space(3)))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Start).Children(func() {
			tile(c, o.Icon, t.Accent)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, name).FontSize(t.Rem(0.9)).Bold().Wrap().FillWidth()
				if o.Description != "" {
					kit.Muted(c, o.Description).FontSize(t.Rem(0.75)).Wrap().FillWidth()
				}
			})
			kit.Chip(c, fmt.Sprintf("%d%%", percent(value)), ink)
		})
		progress(c, value, ink)
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			display.AvatarGroup(c, members, display.AvatarGroupOptions{Max: 4})
			kit.Muted(c, count).FontSize(t.Rem(0.7)).Grow(1).MinWidth(0).SingleLine()
			if o.Due != "" {
				ui.Row(c).Shrink(0).Gap(t.Space(1)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
					glyph(c, "calendar", 12)
					ui.Text(c, o.Due).FontSize(t.Rem(0.7)).Shrink(0)
				})
			}
		})
	})
	return root
}
