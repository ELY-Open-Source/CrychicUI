package data

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// EmptyTableOptions tune an EmptyTable.
type EmptyTableOptions struct {
	// Icon is a Lucide name; "inbox" when empty.
	Icon string
}

// EmptyTable is a table head over a message and action.
func EmptyTable(c *ui.Context, label string, cols []Column, text, action string, o EmptyTableOptions) (root, act ui.Element) {
	columns("EmptyTable", label, cols)
	if text == "" || action == "" {
		panic(fmt.Sprintf("data: EmptyTable %q needs text and an action", label))
	}
	def(&o.Icon, "inbox")
	icon.Must(o.Icon)
	t := c.Theme()
	root = frame(c, label)
	root.Children(func() {
		head(c, root, "EmptyTable", cols, nil, nil)
		box := ui.Column(c).FillWidth().AlignItems(ui.Center).Gap(t.Space(3)).Padding(t.Space(8), t.Space(4))
		kit.Arrive(box, ui.Pose{Y: t.Space(2), Fade: 1}, 0).Children(func() {
			ui.Icon(c, icon.Must(o.Icon)).Size(t.Space(6), t.Space(6)).TextColor(t.TextMuted)
			kit.Muted(c, text)
			act = button.Solid(c, action).AlignSelf(ui.Center)
		})
	})
	if act.Clicked() {
		kit.Log("EmptyTable", "action")
	}
	return root, act
}
