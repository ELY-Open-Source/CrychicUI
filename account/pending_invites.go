package account

import (
	"fmt"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Invite is one pending team invitation; Resent marks a resend.
type Invite struct {
	ID, Email, Role, Sent string
	Resent                bool
}

// PendingInvitesOptions word the card; Label takes action and email.
type PendingInvitesOptions struct{ Title, Resend, Resent, Revoke, Empty, Label string }

// PendingInvites resends and revokes invites; cue Send, Delete.
func PendingInvites(c *ui.Context, invites *[]Invite, o PendingInvitesOptions) ui.Element {
	if invites == nil {
		panic("account: PendingInvites needs invites")
	}
	ids := make([]string, len(*invites))
	for i, inv := range *invites {
		if !validEmail(inv.Email) || inv.Role == "" || inv.Sent == "" {
			panic(fmt.Sprintf("account: PendingInvites %q needs an email, role and sent time", inv.ID))
		}
		ids[i] = inv.ID
	}
	unique("PendingInvites", ids)
	def(&o.Title, "Pending invites")
	def(&o.Resend, "Resend")
	def(&o.Resent, "Sent")
	def(&o.Revoke, "Revoke")
	def(&o.Empty, "No pending invites")
	def(&o.Label, "%s %s")
	t := c.Theme()
	root := card(c).Role(ui.RoleList)
	revoke := -1
	root.Children(func() {
		kit.CardHead(c, "mail", o.Title, func() { kit.Chip(c, fmt.Sprint(len(*invites)), t.Accent) })
		list := ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(1.5))
		list.Children(func() {
			if len(*invites) == 0 {
				empty(c, "mail-check", o.Empty)
			}
			for i, inv := range *invites {
				c.Key(inv.ID)
				row := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).
					Transition(rowExit)
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					tile(c, "mail", t.Accent, t.Space(9))
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
						ui.Text(c, inv.Email).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
						ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
							kit.Muted(c, inv.Role).FontSize(t.Rem(0.75)).MinWidth(0).SingleLine()
							dot(c)
							kit.Muted(c, inv.Sent).FontSize(t.Rem(0.75)).MinWidth(0).SingleLine()
						})
					})
					done := inv.Resent
					text, mark := o.Resend, "rotate-cw"
					if done {
						text, mark = o.Resent, "check"
					}
					send := outline(c, fmt.Sprintf(o.Label, o.Resend, inv.Email), text, mark, t.Text).Disabled(done)
					if send.Clicked() && !done {
						(*invites)[i].Resent = true
						root.SetChanged()
						sound.Play(sound.Send)
						kit.Log("PendingInvites", "resent", "index", i)
					}
					cut := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Label, o.Revoke, inv.Email)).Size(t.Space(8), t.Space(8)).
						Shrink(0).Center().Radius(999))
					cut.Background(cut.AnimateColor("bg", t.Danger.Alpha(0).Mix(t.Danger.Alpha(0.12), unit(cut.Hovered())), ui.Snappy))
					cut.TextColor(cut.AnimateColor("ink", t.TextMuted.Mix(t.Danger, unit(cut.Hovered())), ui.Snappy))
					cut.Children(func() { glyph(c, "x", 16) })
					if cut.Clicked() {
						revoke = i
					}
				})
			}
		})
	})
	if revoke >= 0 {
		*invites = slices.Delete(slices.Clone(*invites), revoke, revoke+1)
		root.SetChanged()
		sound.Play(sound.Delete)
		kit.Log("PendingInvites", "revoked", "count", len(*invites))
	}
	return root
}
