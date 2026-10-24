package code

import (
	"path"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FileHeaderOptions add a meta line and name the copy action.
type FileHeaderOptions struct {
	// Meta follows the language, as "128 lines · 4.2 KB".
	Meta string
	Copy string
}

// FileHeader is a path, language and actions bar.
func FileHeader(c *ui.Context, filePath, lang string, o FileHeaderOptions) (root, copy ui.Element) {
	if filePath == "" || strings.HasSuffix(filePath, "/") {
		panic("code: FileHeader needs a file path")
	}
	def(&o.Copy, "Copy path")
	t := c.Theme()
	dir, base := path.Split(filePath)
	root = ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(2), t.Space(2), t.Space(3)).
		Background(t.Surface).Border(1, t.Border).Radius(t.Radius).Role(ui.RoleGroup).Label(filePath)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(t.Radius).Background(t.Accent.Alpha(0.12)).
			TextColor(t.Accent).Children(func() { ui.Icon(c, icon.Must(glyphFor(base))).Size(16, 16) })
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
			ui.Row(c).FillWidth().Children(func() {
				if dir != "" {
					kit.Muted(c, dir).Font("monospace").FontSize(t.Rem(0.85)).Shrink(1).MinWidth(0).SingleLine()
				}
				ui.Text(c, base).Font("monospace").FontSize(t.Rem(0.85)).Bold().Shrink(0)
			})
			ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				if lang != "" {
					langChip(c, lang)
				}
				if o.Meta != "" {
					kit.Muted(c, o.Meta).FontSize(t.Rem(0.75)).Shrink(0)
				}
			})
		})
		copy = button.CopyButton(c, o.Copy, filePath)
	})
	if copy.Clicked() {
		root.SetChanged()
		kit.Log("FileHeader", "copied", "depth", strings.Count(filePath, "/"))
	}
	return root, copy
}
