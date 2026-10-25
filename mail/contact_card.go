package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ContactCardOptions word the actions and field names.
type ContactCardOptions struct{ Call, Message, Email, Phone string }

// ContactCard shows avatar, name, fields and actions; cue Tap.
func ContactCard(c *ui.Context, person Contact, action *string, o ContactCardOptions) ui.Element {
	validContact("ContactCard", person)
	if action == nil {
		panic("mail: ContactCard needs action state")
	}
	def(&o.Call, "Call")
	def(&o.Message, "Message")
	def(&o.Email, "Email")
	def(&o.Phone, "Phone")
	t := c.Theme()
	root := card(c)
	root.Children(func() {
		ui.Column(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(4), t.Space(3), t.Space(3), t.Space(3)).Children(func() {
			avatar(c, person.Name, 16)
			ui.Text(c, person.Name).FontSize(t.Rem(1)).Bold().Padding(t.Space(1), 0, 0, 0)
			if person.Role != "" {
				kit.Muted(c, person.Role).FontSize(t.Rem(0.78))
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(6)).Justify(ui.Center).Padding(0, t.Space(3), t.Space(3), t.Space(3)).Children(func() {
			for _, item := range []struct{ label, mark, id string }{{o.Call, "phone", "call"}, {o.Message, "message-square", "message"}} {
				b := roundAction(c, item.label, item.mark, t.Accent.Alpha(0.12), t.Accent)
				if b.Clicked() {
					*action = item.id
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("ContactCard", "action", "id", item.id)
				}
			}
		})
		fields := ui.Column(c).FillWidth().Padding(t.Space(1), 0).BorderWidth(1, 0, 0, 0).BorderColor(t.Border)
		fields.Children(func() {
			for _, f := range []struct{ name, mark, value string }{{o.Email, "mail", person.Email}, {o.Phone, "phone", person.Phone}} {
				if f.value == "" {
					continue
				}
				ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Children(func() {
					glyph(c, f.mark, 14).TextColor(t.TextMuted)
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
						kit.Muted(c, f.name).FontSize(t.Rem(0.68))
						ui.Text(c, f.value).FontSize(t.Rem(0.82)).FillWidth().MinWidth(0).SingleLine()
					})
				})
			}
		})
	})
	return root
}
