package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// SegmentNav is a segmented header pushing r's pages.
func SegmentNav(c *ui.Context, label string, r *ui.Router, pages []Page) ui.Element {
	at := routed("SegmentNav", label, r, pages)
	names := make([]string, len(pages))
	for i, p := range pages {
		names[i] = p.Name
	}
	root := choice.Segments(c, label, names, &at)
	if root.Changed() {
		r.Push(pages[at].Path)
		kit.Log("SegmentNav", "pushed", "index", at)
	}
	return root
}
