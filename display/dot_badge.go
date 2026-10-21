package display

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// DotBadgeOptions tune a DotBadge.
type DotBadgeOptions struct {
	// Color is the dot color; zero is Accent.
	Color ui.Color
}

// DotBadge pops a dot over a child's corner.
func DotBadge(c *ui.Context, o DotBadgeOptions, content func()) ui.Element {
	if o.Color.A == 0 {
		o.Color = c.Theme().Accent
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	return root.Children(func() {
		content()
		dot := ui.Box(c).Absolute().Top(-t.Space(0.75)).Right(-t.Space(0.75)).Size(t.Space(2), t.Space(2)).Radius(999).
			Background(o.Color).Border(1, t.Background)
		kit.Arrive(dot, ui.Pose{Scale: 0.5, Fade: 1}, 0)
	})
}
