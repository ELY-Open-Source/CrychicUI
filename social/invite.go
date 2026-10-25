package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// InviteOptions word the card; Label takes action, name.
type InviteOptions struct{ Title, Copy, Copied, Invite, Invited, Label string }

// Invite copies a link and invites contacts; cue Pop.
func Invite(c *ui.Context, link string, contacts []User, invited *map[string]bool, o InviteOptions) (root, copy ui.Element) {
	if link == "" || len(contacts) == 0 || invited == nil || *invited == nil {
		panic("social: Invite needs link, contacts, and state")
	}
	ids := make([]string, len(contacts))
	for i, user := range contacts {
		validUser("Invite", user)
		ids[i] = user.ID
	}
	unique("Invite", ids)
	def(&o.Title, "Invite people")
	def(&o.Copy, "Copy link")
	def(&o.Copied, "Copied")
	def(&o.Invite, "Invite")
	def(&o.Invited, "Invited")
	def(&o.Label, "%s %s")
	t := c.Theme()
	root = card(c).Role(ui.RoleGroup)
	copied := ui.Local(root, "copied", func() string { return "" })
	root.Children(func() {
		kit.CardHead(c, "user-plus", o.Title, nil)
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(3), t.Space(1), t.Space(3)).Children(func() {
			ui.Row(c).Grow(1).MinWidth(0).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(2.5)).Radius(t.Radius).
				Background(t.Surface).TextColor(t.TextMuted).Children(func() {
				glyph(c, "link", 14)
				ui.Text(c, link).FontSize(t.Rem(0.8)).TextColor(t.Text).Grow(1).MinWidth(0).SingleLine()
			})
			done := *copied == link
			label, mark := o.Copy, "copy"
			if done {
				label, mark = o.Copied, "check"
			}
			copy = kit.Press(ui.ButtonBase(c).Label(label).Row().Gap(t.Space(1.5)).AlignItems(ui.Center).Shrink(0).
				Padding(t.Space(2), t.Space(3)).Radius(t.Radius))
			copy.Background(copy.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(copy.Hovered())).Mix(t.Success, unit(done)), ui.Snappy))
			copy.TextColor(t.AccentText)
			copy.Children(func() {
				c.Key(done)
				kit.Arrive(glyph(c, mark, 14), ui.Pose{Scale: 0.5, Fade: 1}, 0)
				ui.Text(c, label).FontSize(t.Rem(0.8)).Bold().Shrink(0).ContentTransition(ui.ContentTransition{})
			})
			if copy.Clicked() {
				c.WriteClipboard(link)
				*copied = link
				copy.SetChanged()
				sound.Play(sound.Pop)
				kit.Log("Invite", "copied")
			}
		})
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			for i, user := range contacts {
				c.Key(user.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(1.5)).Radius(t.Radius)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					avatar(c, user, t.Space(9))
					ui.Column(c).Grow(1).MinWidth(0).Children(func() {
						ui.Text(c, user.Name).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
						kit.Muted(c, user.Handle).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
					})
					on := (*invited)[user.ID]
					text := o.Invite
					if on {
						text = o.Invited
					}
					b := follow(c, fmt.Sprintf(o.Label, text, user.Name), text, "send", on)
					if b.Clicked() {
						(*invited)[user.ID] = !on
						root.SetChanged()
						sound.Play(sound.Pop)
						kit.Log("Invite", "invited", "index", i, "on", !on)
					}
				})
			}
		})
	})
	return root, copy
}
