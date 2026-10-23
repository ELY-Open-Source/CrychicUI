package files

import (
	"fmt"

	"github.com/egoist/mygo/ui"
)

// FileIcon chooses an icon by filename extension.
func FileIcon(c *ui.Context, name string) ui.Element {
	if name == "" {
		panic(fmt.Sprintf("files: FileIcon needs a name, got %q", name))
	}
	return fileIcon(c, name)
}
