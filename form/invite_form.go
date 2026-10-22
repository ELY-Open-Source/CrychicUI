package form

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/input"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// InviteFormOptions tune an InviteForm.
type InviteFormOptions struct {
	// Labels, placeholder, button copy and the bad-address error.
	Emails, Placeholder, Role, Send, EmailError string
}

// InviteForm collects emails as chips and a role to invite.
func InviteForm(c *ui.Context, emails *[]string, role *string, roles []string, o InviteFormOptions) (root, send ui.Element) {
	if emails == nil || role == nil || len(roles) == 0 {
		panic("form: InviteForm needs emails, a role and roles")
	}
	def(&o.Emails, "Emails")
	def(&o.Placeholder, "name@team.com, then Enter")
	def(&o.Role, "Role")
	def(&o.Send, "Send invites")
	def(&o.EmailError, "Remove the addresses that are not emails.")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(4))
	draft := ui.Local(root, "draft", func() string { return "" })
	root.Children(func() {
		bad := slices.ContainsFunc(*emails, func(e string) bool { return !emailPattern.MatchString(e) })
		FieldGroup(c, o.Emails, FieldGroupOptions{Error: when(bad, o.EmailError)}, func() {
			edited(root, input.TagInput(c, emails, draft, input.TagInputOptions{Placeholder: o.Placeholder}))
		})
		FieldGroup(c, o.Role, FieldGroupOptions{}, func() { edited(root, choice.Dropdown(c, o.Role, roles, role)) })
		send = primary(c, o.Send).Disabled(len(*emails) == 0 || bad)
	})
	if kit.Clicked(send, sound.Send) {
		kit.Log("InviteForm", "send", "count", len(*emails))
	}
	return root, send
}
