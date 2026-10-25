package account

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// AvatarEditorOptions word the card.
type AvatarEditorOptions struct{ Title, Hint, Change, Remove string }

// AvatarEditor changes or removes an avatar; cue Pop.
func AvatarEditor(c *ui.Context, image *ui.ImageSource, change *bool, o AvatarEditorOptions) ui.Element {
	if image == nil || change == nil {
		panic("account: AvatarEditor needs image and change state")
	}
	def(&o.Title, "Profile photo")
	def(&o.Hint, "PNG or JPG, at least 400 × 400 px.")
	def(&o.Change, "Change")
	def(&o.Remove, "Remove")
	t := c.Theme()
	root := card(c).Row().Wrap().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		frame := ui.Box(c).Size(t.Space(18), t.Space(18)).Shrink(0)
		frame.Children(func() {
			c.Key(*image != nil)
			face := ui.Box(c).Fill().Radius(999).Clip().Center().Background(t.Accent.Alpha(0.12)).TextColor(t.Accent).
				Border(t.Space(0.75), t.Background)
			kit.Arrive(face, ui.Pose{Scale: 0.8, Fade: 1}, 0).Children(func() {
				if *image != nil {
					ui.Image(c, *image).Fill().Fit(ui.Cover)
				} else {
					glyph(c, "user-round", 32)
				}
			})
			ui.Box(c).Absolute().Right(0).Bottom(0).Size(t.Space(7), t.Space(7)).Center().Radius(999).
				Background(t.Accent).TextColor(t.AccentText).Border(t.Space(0.5), t.Background).
				Children(func() { glyph(c, "camera", 14) })
		})
		ui.Column(c).Grow(1).MinWidth(t.Space(40)).Gap(t.Space(2.5)).Children(func() {
			ui.Column(c).FillWidth().Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, o.Title).FontSize(t.Rem(0.95)).Bold().Wrap().FillWidth()
				kit.Muted(c, o.Hint).FontSize(t.Rem(0.78)).Wrap().FillWidth()
			})
			ui.Row(c).Wrap().Gap(t.Space(2)).Children(func() {
				changeButton := outline(c, o.Change, o.Change, "upload", t.Text)
				if changeButton.Clicked() {
					*change = true
					root.SetChanged()
					sound.Play(sound.Pop)
					kit.Log("AvatarEditor", "change")
				}
				remove := outline(c, o.Remove, o.Remove, "trash-2", t.Danger).Disabled(*image == nil)
				if remove.Clicked() && *image != nil {
					*image = nil
					root.SetChanged()
					sound.Play(sound.Pop)
					kit.Log("AvatarEditor", "removed")
				}
			})
		})
	})
	return root
}
