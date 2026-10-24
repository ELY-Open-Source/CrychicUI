package mail

import (
	"net/mail"
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RecipientsOptions are the recipient group labels.
type RecipientsOptions struct{ To, Cc, Bcc string }

// Recipients edits To, Cc, and Bcc chips; cue Pop.
func Recipients(c *ui.Context, to, cc, bcc *[]string, o RecipientsOptions) ui.Element {
	if to == nil || cc == nil || bcc == nil {
		panic("mail: Recipients needs all groups")
	}
	def(&o.To, "To")
	def(&o.Cc, "Cc")
	def(&o.Bcc, "Bcc")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	removeGroup, removeAt := -1, -1
	root.Children(func() {
		groups := []struct {
			label  string
			people *[]string
		}{{o.To, to}, {o.Cc, cc}, {o.Bcc, bcc}}
		for i, group := range groups {
			c.Key(i)
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(1)).AlignItems(ui.Center).Children(func() {
				kit.Muted(c, group.label).Width(t.Space(10))
				for j, address := range *group.people {
					if _, err := mail.ParseAddress(address); err != nil {
						panic("mail: invalid recipient")
					}
					c.Key(address)
					chip := button(c, address).Row().Gap(t.Space(1)).AlignItems(ui.Center).Radius(999).Background(t.Surface)
					chip.Children(func() {
						ui.Text(c, address).SingleLine()
						ui.Icon(c, icon.Must("x")).Size(16, 16)
					})
					if chip.Clicked() {
						removeGroup, removeAt = i, j
					}
				}
			})
		}
	})
	if removeAt >= 0 {
		people := []*[]string{to, cc, bcc}[removeGroup]
		*people = slices.Delete(*people, removeAt, removeAt+1)
		root.SetChanged()
		sound.Play(sound.Pop)
		kit.Log("Recipients", "removed", "group", removeGroup, "count", len(*people))
	}
	return root
}
