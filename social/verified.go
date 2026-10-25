package social

import (
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
	root := ui.Box(c).Size(t.Space(6), t.Space(6)).Shrink(0).Center().Radius(999).Label(label).Focusable().Tooltip(o.Tooltip).
		TextColor(t.Accent)
	hover := unit(root.Hovered() || root.Focused())
	root.Background(root.AnimateColor("bg", t.Accent.Alpha(0).Mix(t.Accent.Alpha(0.16), hover), ui.Snappy))
	kit.Arrive(root, ui.Pose{Scale: 0.4, Rotate: -45, Fade: 1}, 0)
	root.Children(func() { glyph(c, "badge-check", 18) })
	return root
}
