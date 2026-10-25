package social

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// InviteOptions are the invite labels.
type InviteOptions struct{ Copy, Copied, Invite, Invited string }

// Invite copies a link and invites contacts; cue Pop.
func Invite(c *ui.Context, link string, contacts []User, invited *map[string]bool, o InviteOptions) ui.Element {
	if link == "" || len(contacts) == 0 || invited == nil || *invited == nil {
		panic("social: Invite needs link, contacts, and state")
	}
	def(&o.Copy, "Copy link")
	def(&o.Copied, "Copied")
	def(&o.Invite, "Invite")
	def(&o.Invited, "Invited")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	copied := ui.Local(root, "copied", func() string { return "" })
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Text(c, link).Grow(1).SingleLine()
			label := o.Copy
			if *copied == link {
				label = o.Copied
			}
			b := button(c, label)
			b.Children(func() { ui.Text(c, label).ContentTransition(ui.ContentTransition{}) })
			if b.Clicked() {
				c.WriteClipboard(link)
				*copied = link
				root.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("Invite", "copied")
			}
		})
		for _, user := range contacts {
			validUser("Invite", user)
			c.Key(user.ID)
			ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				avatar(c, user, t.Space(7))
				ui.Text(c, user.Name).Grow(1)
				on := (*invited)[user.ID]
				label := o.Invite
				if on {
					label = o.Invited
				}
				b := button(c, user.ID+" "+label).Checked(on)
				b.Children(func() { ui.Text(c, label) })
				if b.Clicked() {
					(*invited)[user.ID] = !on
					root.SetChanged()
					sound.Play(sound.Pop)
					kit.Log("Invite", "invited", "on", !on)
				}
			})
		}
	})
	return root
}
