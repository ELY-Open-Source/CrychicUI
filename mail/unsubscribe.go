package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UnsubscribeOptions are the unsubscribe labels.
type UnsubscribeOptions struct{ Title, Confirm string }

// Unsubscribe confirms a selected reason; cue Success.
func Unsubscribe(c *ui.Context, reasons []string, reason *string, confirmed *bool, o UnsubscribeOptions) (root, confirm ui.Element) {
	if len(reasons) == 0 || reason == nil || confirmed == nil {
		panic("mail: Unsubscribe needs reasons and state")
	}
	def(&o.Title, "Why are you leaving?")
	def(&o.Confirm, "Unsubscribe")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	validReason := false
	root.Children(func() {
		ui.Text(c, o.Title).Bold()
		for i, item := range reasons {
			if item == "" {
				panic("mail: unsubscribe reason must not be empty")
			}
			validReason = validReason || *reason == item
			c.Key(item)
			b := button(c, item).FillWidth().Checked(*reason == item).Disabled(*confirmed)
			b.Children(func() { ui.Text(c, item) })
			if b.Clicked() && !*confirmed && *reason != item {
				*reason = item
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("Unsubscribe", "reason", "index", i)
			}
		}
		confirm = button(c, o.Confirm).Background(t.Danger).TextColor(t.AccentText).Disabled(!validReason || *confirmed)
		confirm.Children(func() { ui.Text(c, o.Confirm) })
	})
	if confirm.Clicked() && validReason && !*confirmed {
		*confirmed = true
		confirm.SetChanged()
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("Unsubscribe", "confirmed")
	}
	return root, confirm
}
