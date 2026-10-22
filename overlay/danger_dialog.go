package overlay

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DangerDialogOptions tune a DangerDialog.
type DangerDialogOptions struct {
	Cancel, Confirm string
}

// DangerDialog confirms a destructive action in the danger color.
func DangerDialog(c *ui.Context, title, text string, open *bool, o DangerDialogOptions) (root, confirm ui.Element) {
	if title == "" || open == nil {
		panic(fmt.Sprintf("overlay: DangerDialog %q needs open", title))
	}
	def(&o.Cancel, "Cancel")
	def(&o.Confirm, "Delete")
	t := c.Theme()
	root = ui.Box(c)
	if dialog(c, open, t.Space(100), func() {
		ui.Row(c).Gap(t.Space(3)).AlignItems(ui.Start).Children(func() {
			badge := ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(999).Background(t.Danger.Alpha(0.12))
			kit.Arrive(badge, ui.Pose{Scale: 0.6, Fade: 1}, 0).Children(func() {
				ui.Icon(c, icon.Must("triangle-alert")).Size(16, 16).TextColor(t.Danger)
			})
			ui.Column(c).Grow(1).Gap(t.Space(1.5)).Children(func() {
				ui.Text(c, title).Bold().FontSize(t.Rem(1.1))
				kit.Muted(c, text)
			})
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).Justify(ui.End).Children(func() {
			if button.Ghost(c, o.Cancel).Clicked() {
				*open = false
				root.SetChanged()
			}
			confirm = button.Destructive(c, o.Confirm)
			if confirm.Clicked() {
				*open = false
				kit.Log("DangerDialog", "confirmed")
			}
		})
	}).Dismissed() {
		root.SetChanged()
	}
	if opened, _ := moved(root, *open, "DangerDialog"); opened {
		sound.Play(sound.Warning)
	}
	return root, confirm
}
