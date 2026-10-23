package files

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// UploadAction is an upload or creation action.
type UploadAction uint8

const (
	UploadNone UploadAction = iota
	UploadFiles
	UploadNewFolder
	UploadNewFile
)

// UploadButtonRowOptions are the upload labels.
type UploadButtonRowOptions struct {
	Upload, New, Folder, File string
}

// UploadButtonRow shows upload and new menu; cue Open.
func UploadButtonRow(c *ui.Context, action *UploadAction, o UploadButtonRowOptions) (root, upload, newMenu ui.Element) {
	if action == nil || *action > UploadNewFile {
		panic("files: UploadButtonRow needs a valid action")
	}
	def(&o.Upload, "Upload")
	def(&o.New, "New")
	def(&o.Folder, "New folder")
	def(&o.File, "New file")
	t := c.Theme()
	root = ui.Row(c).Gap(t.Space(2)).Role(ui.RoleGroup)
	root.Children(func() {
		upload = kit.Press(ui.ButtonBase(c).Label(o.Upload).Row().Gap(t.Space(1)).Padding(t.Space(1.5), t.Space(3)).
			Radius(t.Radius).Background(t.Accent).TextColor(t.AccentText))
		upload.Children(func() {
			ui.Icon(c, icon.Must("upload")).Size(16, 16)
			ui.Text(c, o.Upload).SingleLine()
		})
		if upload.Clicked() {
			*action = UploadFiles
			root.SetChanged()
			sound.Play(sound.Open)
			kit.Log("UploadButtonRow", "upload")
		}
		newMenu = kit.Press(ui.ButtonBase(c).Label(o.New).Row().Gap(t.Space(1)).Padding(t.Space(1.5), t.Space(3)).
			Radius(t.Radius).Border(1, t.Border))
		newMenu.Children(func() {
			ui.Icon(c, icon.Must("plus")).Size(16, 16)
			ui.Text(c, o.New).SingleLine()
		})
		newMenu.Menu(func(m *ui.Menu) {
			for _, item := range []struct {
				name   string
				action UploadAction
			}{{o.Folder, UploadNewFolder}, {o.File, UploadNewFile}} {
				if m.Item(item.name).Chosen() {
					*action = item.action
					root.SetChanged()
					sound.Play(sound.Open)
					kit.Log("UploadButtonRow", "new", "action", *action)
				}
			}
		})
	})
	return root, upload, newMenu
}
