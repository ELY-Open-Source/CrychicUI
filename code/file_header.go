package code

import (
	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FileHeaderOptions name the copy action.
type FileHeaderOptions struct {
	Copy string
}

// FileHeader is a path, language and actions bar.
func FileHeader(c *ui.Context, path, lang string, o FileHeaderOptions) (root, copy ui.Element) {
	if path == "" {
		panic("code: FileHeader needs a path")
	}
	def(&o.Copy, "Copy path")
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).
		Background(t.Surface).Border(1, t.Border).Radius(t.Radius).Role(ui.RoleGroup).Label(path)
	root.Children(func() {
		ui.Icon(c, icon.Must("file-code")).Size(16, 16).TextColor(t.TextMuted)
		ui.Text(c, path).Font("monospace").FontSize(t.Rem(0.9)).Grow(1).SingleLine()
		if lang != "" {
			kit.Muted(c, lang).FontSize(t.Rem(0.75)).Padding(t.Space(0.5), t.Space(1.5)).Radius(999).Background(t.Background)
		}
		copy = button.CopyButton(c, o.Copy, path)
	})
	if copy.Clicked() {
		root.SetChanged()
		kit.Log("FileHeader", "copied")
	}
	return root, copy
}
