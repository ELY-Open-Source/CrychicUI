package files

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileAction identifies a chosen file menu action.
type FileAction uint8

const (
	FileActionNone FileAction = iota
	FileActionOpen
	FileActionRename
	FileActionDelete
)

// FileCardOptions are the file menu labels.
type FileCardOptions struct {
	Menu, Open, Rename, Delete string
}

// FileCard shows file details and a menu; it lifts.
func FileCard(c *ui.Context, item FileItem, action *FileAction, o FileCardOptions) (root, menu ui.Element) {
	validFile("FileCard", item)
	if action == nil || *action > FileActionDelete {
		panic("files: FileCard needs a valid bound action")
	}
	def(&o.Menu, "File actions")
	def(&o.Open, "Open")
	def(&o.Rename, "Rename")
	def(&o.Delete, "Delete")
	t := c.Theme()
	root = kit.Lift(ui.Row(c).AlignSelf(ui.Start).Width(t.Space(56)).Gap(t.Space(2)).AlignItems(ui.Center).
		Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).Background(t.Surface))
	root.Children(func() {
		ui.Icon(c, icon.Must(itemGlyph(item))).Size(t.Space(8), t.Space(8)).TextColor(t.Accent)
		ui.Column(c).Grow(1).Children(func() {
			ui.Text(c, item.Name).Bold().SingleLine()
			sizeText(c, item.Size)
		})
		menu = kit.Press(ui.ButtonBase(c).Label(o.Menu).Tooltip(o.Menu).Size(t.Space(7), t.Space(7)).Center().Radius(999))
		menu.Children(func() { ui.Icon(c, icon.Must("more-horizontal")).Size(16, 16) })
		menu.Menu(func(m *ui.Menu) {
			for _, item := range []struct {
				name   string
				action FileAction
			}{{o.Open, FileActionOpen}, {o.Rename, FileActionRename}, {o.Delete, FileActionDelete}} {
				if m.Item(item.name).Chosen() {
					*action = item.action
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("FileCard", "menu", "action", item.action)
				}
			}
		})
	})
	return root, menu
}
