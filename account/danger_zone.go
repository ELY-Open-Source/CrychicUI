package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DangerZoneOptions word the card; Confirm and Label are formats.
type DangerZoneOptions struct {
	Title, Transfer, TransferNote, TransferAction string
	Delete, DeleteNote, DeleteAction, Confirm     string
	Forever, Cancel, Label                        string
}

// DangerZone transfers or deletes an account; cue Warning, Delete.
func DangerZone(c *ui.Context, account string, action *string, o DangerZoneOptions) ui.Element {
	if account == "" || action == nil {
		panic(fmt.Sprintf("account: DangerZone %q needs an account and action", account))
	}
	def(&o.Title, "Danger zone")
	def(&o.Transfer, "Transfer ownership")
	def(&o.TransferNote, "Hand this workspace to another admin.")
	def(&o.TransferAction, "Transfer")
	def(&o.Delete, "Delete account")
	def(&o.DeleteNote, "Removes all data. This cannot be undone.")
	def(&o.DeleteAction, "Delete")
	def(&o.Confirm, "Type %s to confirm")
	def(&o.Forever, "Delete forever")
	def(&o.Cancel, "Cancel")
	def(&o.Label, "%s %s")
	t := c.Theme()
	root := card(c).Border(1, t.Danger.Alpha(0.5)).Role(ui.RoleGroup)
	armed := ui.Local(root, "armed", func() bool { return false })
	typed := ui.Local(root, "typed", func() string { return "" })
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(4)).
			BorderWidth(0, 0, 1, 0).BorderColor(t.Danger.Alpha(0.3)).Background(t.Danger.Alpha(0.05)).TextColor(t.Danger).
			Children(func() {
				glyph(c, "shield-alert", 16)
				ui.Text(c, o.Title).FontSize(t.Rem(0.85)).Bold().Shrink(0)
			})
		row := func(title, note string, action func()) {
			ui.Row(c).FillWidth().Wrap().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(4)).Children(func() {
				ui.Column(c).Grow(1).MinWidth(t.Space(40)).Gap(t.Space(0.5)).Children(func() {
					ui.Text(c, title).FontSize(t.Rem(0.88)).Bold().Wrap().FillWidth()
					kit.Muted(c, note).FontSize(t.Rem(0.78)).Wrap().FillWidth()
				})
				action()
			})
		}
		row(o.Transfer, o.TransferNote, func() {
			b := outline(c, fmt.Sprintf(o.Label, o.TransferAction, account), o.TransferAction, "arrow-left-right", t.Text)
			if b.Clicked() {
				*action = "transfer"
				root.SetChanged()
				sound.Play(sound.Warning)
				kit.Log("DangerZone", "transfer")
			}
		})
		ui.Box(c).FillWidth().Height(1).Background(t.Border)
		row(o.Delete, o.DeleteNote, func() {
			b := outline(c, fmt.Sprintf(o.Label, o.DeleteAction, account), o.DeleteAction, "trash-2", t.Danger).Disabled(*armed)
			if b.Clicked() && !*armed {
				*armed, *typed = true, ""
				sound.Play(sound.Warning)
				kit.Log("DangerZone", "armed")
			}
		})
		if !*armed {
			return
		}
		c.Key("confirm")
		panel := ui.Column(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3), t.Space(4), t.Space(4), t.Space(4)).
			Background(t.Danger.Alpha(0.05))
		kit.Arrive(panel, ui.Pose{Y: -t.Space(2), Fade: 1}, 0).Children(func() {
			field(c, fmt.Sprintf(o.Confirm, account), "keyboard", typed, false)
			ui.Row(c).FillWidth().Gap(t.Space(2)).Justify(ui.End).Children(func() {
				cancel := outline(c, o.Cancel, o.Cancel, "", t.Text)
				if cancel.Clicked() {
					*armed, *typed = false, ""
					sound.Play(sound.Close)
					kit.Log("DangerZone", "cancelled")
				}
				ok := *typed == account
				forever := kit.Press(ui.ButtonBase(c).Label(o.Forever).Padding(t.Space(2), t.Space(3.5)).Radius(t.Radius).Disabled(!ok))
				forever.Background(forever.AnimateColor("bg", t.Danger.Mix(t.Text, 0.15*unit(forever.Hovered())), ui.Snappy)).
					TextColor(t.AccentText).Children(func() { content(c, o.Forever, "trash-2") })
				if forever.Clicked() && ok {
					*action, *armed, *typed = "delete", false, ""
					root.SetChanged()
					sound.Play(sound.Delete)
					kit.Log("DangerZone", "deleted")
				}
			})
		})
	})
	return root
}
