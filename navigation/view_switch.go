package navigation

import (
	"github.com/ELY-Open-Source/CrychicUI/choice"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ViewSwitchOptions name the three views.
type ViewSwitchOptions struct {
	List, Grid, Board string
}

// ViewSwitch picks a list, grid or board view by icon.
func ViewSwitch(c *ui.Context, label string, view *int, o ViewSwitchOptions) ui.Element {
	def(&o.List, "List")
	def(&o.Grid, "Grid")
	def(&o.Board, "Board")
	items := []choice.SegmentItem{{Label: o.List, Icon: "list"}, {Label: o.Grid, Icon: "layout-grid"}, {Label: o.Board, Icon: "kanban"}}
	root := choice.IconSegments(c, label, items, view)
	if root.Changed() {
		kit.Log("ViewSwitch", "switched", "index", *view)
	}
	return root
}
