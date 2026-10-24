package social

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// VerifiedOptions are the verification labels.
type VerifiedOptions struct{ Tooltip string }

// Verified shows a check badge with tooltip.
func Verified(c *ui.Context, label string, o VerifiedOptions) ui.Element {
	if label == "" {
		panic("social: Verified needs a label")
	}
	def(&o.Tooltip, "Verified account")
	t := c.Theme()
	root := ui.Box(c).Size(t.Space(6), t.Space(6)).Center().Label(label).Focusable().Tooltip(o.Tooltip)
	kit.Arrive(root, ui.Pose{Scale: 0.8, Fade: 1}, 0)
	root.Children(func() { ui.Icon(c, icon.Must("badge-check")).Size(t.Space(5), t.Space(5)).TextColor(t.Accent) })
	return root
}
