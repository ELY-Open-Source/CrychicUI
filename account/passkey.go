package account

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PasskeyOptions are the passkey prompt labels.
type PasskeyOptions struct{ Title, Description, Add, Added string }

// Passkey adds a passkey; cue Success.
func Passkey(c *ui.Context, added *bool, o PasskeyOptions) ui.Element {
	if added == nil {
		panic("account: Passkey needs state")
	}
	def(&o.Title, "Add a passkey")
	def(&o.Description, "Use this device to sign in")
	def(&o.Add, "Add passkey")
	def(&o.Added, "Added")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Text(c, o.Title).Bold()
		ui.Text(c, o.Description).Wrap()
		label := o.Add
		if *added {
			label = o.Added
		}
		b := primary(c, label).Disabled(*added)
		b.Children(func() { ui.Text(c, label) })
		if b.Clicked() && !*added {
			*added = true
			root.SetChanged()
			sound.Play(sound.Success)
			kit.Log("Passkey", "added")
		}
	})
	return root
}
