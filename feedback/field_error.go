package feedback

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// FieldError is error text that shakes in; empty shows nothing.
func FieldError(c *ui.Context, message string) ui.Element {
	t := c.Theme()
	root := ui.Column(c).FillWidth()
	shown := ui.Local(root, "shown", func() string { return "" })
	runs := ui.Local(root, "runs", func() uint64 { return 0 })
	if message != *shown {
		*shown = message
		*runs++
	}
	return root.Children(func() {
		if message == "" {
			return
		}
		c.Key("error")
		row := ui.Row(c).Gap(t.Space(1.5)).Role(ui.RoleStatus).Label(message).TextColor(t.Danger)
		kit.Arrive(row, ui.Pose{Fade: 1}, 0).Transition(collapse)
		row.Translate(row.Keyframes("shake", shake, ui.KeyframeOptions{Run: *runs}), 0).Children(func() {
			ui.Icon(c, icon.Must("circle-alert")).Size(16, 16)
			ui.Text(c, message).FontSize(t.Rem(0.85))
		})
	})
}
