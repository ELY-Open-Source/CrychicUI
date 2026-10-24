package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileGridOptions word a folder's tile.
type FileGridOptions struct {
	Folder string
}

// FileGrid shows files as selectable icon tiles; cue Select.
func FileGrid(c *ui.Context, label string, list []FileItem, selected *int, o FileGridOptions) ui.Element {
	validFiles("FileGrid", label, list)
	unique("FileGrid", label, list)
	if selected == nil || *selected < -1 || *selected >= len(list) {
		panic(fmt.Sprintf("files: FileGrid %q needs a selection or -1", label))
	}
	def(&o.Folder, "Folder")
	t := c.Theme()
	root := ui.Grid(c).Columns(3).FillWidth().Gap(t.Space(2)).FocusGroup(ui.Horizontal).Role(ui.RoleList).Label(label)
	root.Children(func() {
		for i, f := range list {
			c.Key(fileKey(f))
			on := i == *selected
			tile := kit.Press(ui.ButtonBase(c).Label(f.Name).Column().Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(t.Space(3), t.Space(2)).Radius(t.Radius).Checked(on))
			kit.Lift(tile)
			kit.Stagger(c, root, tile, i)
			tile.Background(tile.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(tile.Hovered())).Mix(t.Accent.Alpha(0.1), unit(on)), ui.Snappy))
			tile.Border(1, tile.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
			tile.Children(func() {
				fileTile(c, f, t.Space(12))
				ui.Text(c, f.Name).FontSize(t.Rem(0.8)).Bold().TextAlign(ui.Center).FillWidth().SingleLine()
				if f.Dir {
					kit.Muted(c, o.Folder).FontSize(t.Rem(0.7))
				} else {
					sizeText(c, f.Size).FontSize(t.Rem(0.7))
				}
			})
			if tile.Clicked() && !on {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("FileGrid", "selected", "index", i)
			}
		}
	})
	return root
}
