package account

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/display"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Member is one team member; an Owner cannot be removed.
type Member struct {
	ID, Name, Email, Role string
	Owner                 bool
}

// TeamMembersOptions word the card; Label takes action and name.
type TeamMembersOptions struct{ Title, Remove, Empty, Label string }

// TeamMembers lists members with roles and remove; cue Delete.
func TeamMembers(c *ui.Context, members *[]Member, o TeamMembersOptions) ui.Element {
	if members == nil {
		panic("account: TeamMembers needs members")
	}
	ids := make([]string, len(*members))
	for i, m := range *members {
		if m.Name == "" || m.Email == "" || m.Role == "" {
			panic(fmt.Sprintf("account: TeamMembers %q needs a name, email and role", m.ID))
		}
		ids[i] = m.ID
	}
	unique("TeamMembers", ids)
	def(&o.Title, "Team")
	def(&o.Remove, "Remove")
	def(&o.Empty, "No team members yet")
	def(&o.Label, "%s %s")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	remove := -1
	root.Children(func() {
		kit.CardHead(c, "users", o.Title, func() { kit.Chip(c, fmt.Sprint(len(*members)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			if len(*members) == 0 {
				empty(c, "users", o.Empty)
			}
			for i, m := range *members {
				c.Key(m.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
					Transition(rowExit)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					display.Avatar(c, m.Name, nil).Size(t.Space(9), t.Space(9)).FontSize(t.Rem(0.75)).Shrink(0)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, m.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
						kit.Muted(c, m.Email).FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine()
					})
					ink := t.TextMuted
					if m.Owner {
						ink = t.Warning
					}
					kit.Chip(c, m.Role, ink)
					if m.Owner {
						return
					}
					b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Label, o.Remove, m.Name)).Size(t.Space(8), t.Space(8)).
						Shrink(0).Center().Radius(999))
					b.Background(b.AnimateColor("bg", t.Danger.Alpha(0).Mix(t.Danger.Alpha(0.12), unit(b.Hovered())), ui.Snappy))
					b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Danger, unit(b.Hovered())), ui.Snappy))
					b.Children(func() { glyph(c, "user-x", 16) })
					if b.Clicked() {
						remove = i
					}
				})
			}
		})
	})
	if remove >= 0 {
		*members = slices.Delete(slices.Clone(*members), remove, remove+1)
		root.SetChanged()
		sound.Play(sound.Delete)
		kit.Log("TeamMembers", "removed", "count", len(*members))
	}
	return root
}
