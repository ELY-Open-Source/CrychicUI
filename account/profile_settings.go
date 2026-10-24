package account

import (
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Profile is one editable account profile.
type Profile struct {
	Name, Handle, Bio string
}

// ProfileSettingsOptions are the profile labels.
type ProfileSettingsOptions struct{ Name, Handle, Bio, Save string }

// ProfileSettings saves name, handle, and bio; cue Success.
func ProfileSettings(c *ui.Context, profile *Profile, o ProfileSettingsOptions) (root, save ui.Element) {
	if profile == nil {
		panic("account: ProfileSettings needs a profile")
	}
	def(&o.Name, "Name")
	def(&o.Handle, "Handle")
	def(&o.Bio, "Bio")
	def(&o.Save, "Save")
	t := c.Theme()
	root = ui.Column(c).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		name := ui.TextInput(c, &profile.Name).Label(o.Name).Placeholder(o.Name).FillWidth()
		handle := ui.TextInput(c, &profile.Handle).Label(o.Handle).Placeholder(o.Handle).FillWidth()
		bio := ui.TextArea(c, &profile.Bio).Label(o.Bio).Placeholder(o.Bio).FillWidth().Lines(3, 6)
		if name.Changed() || handle.Changed() || bio.Changed() {
			root.SetChanged()
		}
		valid := strings.TrimSpace(profile.Name) != "" && strings.TrimSpace(profile.Handle) != ""
		save = primary(c, o.Save).Disabled(!valid)
		save.Children(func() { ui.Text(c, o.Save) })
		if save.Clicked() && valid {
			save.SetChanged()
			root.SetChanged()
			sound.Play(sound.Success)
			kit.Log("ProfileSettings", "saved", "bioRunes", len([]rune(profile.Bio)))
		}
	})
	return root, save
}
