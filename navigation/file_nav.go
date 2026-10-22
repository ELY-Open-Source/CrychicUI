package navigation

import (
	"path"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// FileNav is a folder tree with icons by file kind.
func FileNav(c *ui.Context, label string, list []Node, chosen *string, open map[string]bool) ui.Element {
	return tree(c, "FileNav", label, list, chosen, open, fileIcon, sound.Select)
}

// fileIcon picks an icon by folder or extension.
func fileIcon(n Node, open bool) string {
	switch {
	case n.Children != nil && open:
		return "folder-open"
	case n.Children != nil:
		return "folder"
	}
	switch strings.ToLower(path.Ext(n.Label)) {
	case ".go", ".js", ".ts", ".py", ".rs", ".c", ".css", ".html", ".json":
		return "file-code"
	case ".md", ".txt":
		return "file-text"
	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp":
		return "file-image"
	}
	return "file"
}
