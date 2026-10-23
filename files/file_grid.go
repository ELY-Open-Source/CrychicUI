package files

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileGrid shows selectable icon tiles; cue Select.
func FileGrid(c *ui.Context, label string, list []FileItem, selected *int) ui.Element {
	validFiles("FileGrid", label, list)
	if selected == nil || *selected < -1 || *selected >= len(list) {
		panic("files: FileGrid needs a valid selection")
	}
	t := c.Theme()
	root := ui.Grid(c).Columns(3).FillWidth().Gap(t.Space(2)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		for i, f := range list {
			c.Key(i)
			tile := kit.Press(ui.ButtonBase(c).Label(f.Name).Column().Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).Checked(i == *selected))
			bg := t.Surface.Mix(t.SurfaceHover, unit(tile.Hovered()))
			if i == *selected {
				bg = t.Accent.Alpha(0.12)
			}
			tile.Background(tile.AnimateColor("bg", bg, ui.Snappy))
			tile.Children(func() {
				ui.Icon(c, icon.Must(itemGlyph(f))).Size(t.Space(8), t.Space(8)).TextColor(t.Accent)
				ui.Text(c, f.Name).FontSize(t.Rem(0.85)).SingleLine()
			})
			if tile.Clicked() && *selected != i {
				*selected = i
				root.SetChanged()
				sound.Play(sound.Select)
				kit.Log("FileGrid", "selected", "index", i)
			}
		}
	})
	return root
}
