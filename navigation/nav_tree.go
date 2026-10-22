package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// NavTree is nested links that expand; open holds folder paths.
func NavTree(c *ui.Context, label string, list []Node, chosen *string, open map[string]bool) ui.Element {
	return tree(c, "NavTree", label, list, chosen, open, nil, sound.Tap)
}
