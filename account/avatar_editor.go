package account

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AvatarEditorOptions are the avatar actions.
type AvatarEditorOptions struct{ Change, Remove string }

// AvatarEditor changes or removes an avatar; cue Pop.
func AvatarEditor(c *ui.Context, image *ui.ImageSource, change *bool, o AvatarEditorOptions) ui.Element {
	if image == nil || change == nil {
		panic("account: AvatarEditor needs image and change state")
	}
	def(&o.Change, "Change")
	def(&o.Remove, "Remove")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		box := ui.Box(c).Size(t.Space(16), t.Space(16)).Radius(999).Clip().Background(t.Accent.Alpha(0.15)).Center()
		box.Children(func() {
			if *image != nil {
				ui.Image(c, *image).Fill().Fit(ui.Cover)
			} else {
				ui.Icon(c, icon.Must("user")).Size(t.Space(8), t.Space(8)).TextColor(t.Accent)
			}
		})
		changeButton := button(c, o.Change)
		changeButton.Children(func() { ui.Text(c, o.Change) })
		if changeButton.Clicked() {
			*change = true
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("AvatarEditor", "change")
		}
		remove := button(c, o.Remove).Disabled(*image == nil)
		remove.Children(func() { ui.Text(c, o.Remove) })
		if remove.Clicked() && *image != nil {
			*image = nil
			root.SetChanged()
			sound.Play(sound.Pop)
			kit.Log("AvatarEditor", "removed")
		}
	})
	return root
}
