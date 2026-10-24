package files

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileAction identifies a chosen file menu action.
type FileAction uint8

// File menu actions.
const (
	FileActionNone FileAction = iota
	FileActionOpen
	FileActionRename
	FileActionDelete
)

// FileCardOptions word the menu and dates.
type FileCardOptions struct {
	Menu, Open, Rename, Delete string
	FormatDate                 func(t time.Time) string
}

// FileCard shows a preview band, details and menu; it lifts.
func FileCard(c *ui.Context, item FileItem, action *FileAction, o FileCardOptions) (root, menu ui.Element) {
	validFile("FileCard", item)
	if action == nil || *action > FileActionDelete {
		panic(fmt.Sprintf("files: FileCard %q needs a known bound action", item.Name))
	}
	def(&o.Menu, "File actions")
	def(&o.Open, "Open")
	def(&o.Rename, "Rename")
	def(&o.Delete, "Delete")
	if o.FormatDate == nil {
		o.FormatDate = shortDate
	}
	t := c.Theme()
	ink := kindInk(t, item)
	root = kit.Lift(kit.Card(c).FillWidth().Clip().Role(ui.RoleGroup).Label(item.Name))
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Column(c).FillWidth().Height(t.Space(24)).Center().Background(ink.Alpha(0.08)).Children(func() {
			fileTile(c, item, t.Space(14))
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2.5), t.Space(2), t.Space(2.5), t.Space(3)).Children(func() {
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, item.Name).Bold().FontSize(t.Rem(0.85)).FillWidth().SingleLine()
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					sizeText(c, item.Size).FontSize(t.Rem(0.75))
					kit.Muted(c, "·").FontSize(t.Rem(0.75))
					kit.Muted(c, o.FormatDate(item.Modified)).FontSize(t.Rem(0.75)).Shrink(0)
				})
			})
			menu = kit.Press(ui.ButtonBase(c).Label(o.Menu).Tooltip(o.Menu).Size(t.Space(7), t.Space(7)).Shrink(0).Center().Radius(999))
			menu.Background(menu.AnimateColor("bg", t.SurfaceHover.Alpha(unit(menu.Hovered())), ui.Snappy)).
				Children(func() { ui.Icon(c, icon.Must("ellipsis")).Size(16, 16) })
			menu.Menu(func(m *ui.Menu) {
				for _, it := range []struct {
					name   string
					action FileAction
				}{{o.Open, FileActionOpen}, {o.Rename, FileActionRename}, {o.Delete, FileActionDelete}} {
					if m.Item(it.name).Chosen() {
						*action = it.action
						root.SetChanged()
						sound.Play(sound.Tap)
						kit.Log("FileCard", "menu", "action", int(it.action))
					}
				}
			})
		})
	})
	return root, menu
}
