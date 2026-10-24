package files

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/egoist/mygo/ui"
)

// FileIconOptions size a FileIcon.
type FileIconOptions struct {
	// Size is the tile's side; zero gives the default.
	Size float32
}

// FileIcon is a tile picked by extension; it swaps.
func FileIcon(c *ui.Context, name string, o FileIconOptions) ui.Element {
	if name == "" {
		panic(fmt.Sprintf("files: FileIcon needs a name, got %q", name))
	}
	t := c.Theme()
	dimension("FileIcon", &o.Size, t.Space(12))
	f := FileItem{Name: name}
	ink := kindInk(t, f)
	ext := strings.ToUpper(strings.TrimPrefix(filepath.Ext(name), "."))
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(1)).AlignItems(ui.Center).Role(ui.RoleImage).Label(name)
	return root.Children(func() {
		fileTile(c, f, o.Size)
		if ext != "" {
			ui.Text(c, ext).FontSize(t.Rem(0.6)).Bold().TextColor(ink).Padding(0, t.Space(1.5)).Radius(999).
				Background(ink.Alpha(0.12)).Shrink(0).ContentTransition(ui.ContentTransition{})
		}
	})
}
