package files

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UploadAction is an upload or creation action.
type UploadAction uint8

// Upload row actions.
const (
	UploadNone UploadAction = iota
	UploadFiles
	UploadNewFolder
	UploadNewFile
)

// UploadButtonRowOptions word the buttons and menu.
type UploadButtonRowOptions struct {
	Upload, New, Folder, File string
}

// UploadButtonRow is Upload and a New menu; cue Open.
func UploadButtonRow(c *ui.Context, action *UploadAction, o UploadButtonRowOptions) (root, upload, newMenu ui.Element) {
	if action == nil || *action > UploadNewFile {
		panic("files: UploadButtonRow needs a known bound action")
	}
	def(&o.Upload, "Upload")
	def(&o.New, "New")
	def(&o.Folder, "New folder")
	def(&o.File, "New file")
	t := c.Theme()
	root = ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(2)).FocusGroup(ui.Horizontal).Role(ui.RoleToolbar).Label(o.Upload + " · " + o.New)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		upload = kit.Press(ui.ButtonBase(c).Label(o.Upload).Padding(t.Space(1.5), t.Space(3)).Radius(t.Radius))
		upload.Background(upload.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(upload.Hovered())), ui.Snappy)).TextColor(t.AccentText)
		upload.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				lift := upload.AnimateTo("lift", unit(upload.Hovered()), ui.Snappy)
				ui.Icon(c, icon.Must("upload")).Size(16, 16).Translate(0, -t.Space(0.5)*lift)
				ui.Text(c, o.Upload).Bold().Shrink(0)
			})
		})
		if upload.Clicked() {
			*action = UploadFiles
			root.SetChanged()
			sound.Play(sound.Open)
			kit.Log("UploadButtonRow", "upload")
		}
		newMenu = kit.Press(ui.ButtonBase(c).Label(o.New).Padding(t.Space(1.5), t.Space(2.5), t.Space(1.5), t.Space(3)).Radius(t.Radius).Border(1, t.Border))
		newMenu.Background(newMenu.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(newMenu.Hovered())), ui.Snappy))
		newMenu.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				ui.Icon(c, icon.Must("plus")).Size(16, 16)
				ui.Text(c, o.New).Bold().Shrink(0)
				ui.Icon(c, icon.Must("chevron-down")).Size(14, 14).TextColor(t.TextMuted)
			})
		})
		newMenu.Menu(func(m *ui.Menu) {
			for _, it := range []struct {
				name   string
				action UploadAction
			}{{o.Folder, UploadNewFolder}, {o.File, UploadNewFile}} {
				if m.Item(it.name).Chosen() {
					*action = it.action
					root.SetChanged()
					sound.Play(sound.Open)
					kit.Log("UploadButtonRow", "new", "action", int(it.action))
				}
			}
		})
	})
	return root, upload, newMenu
}
