package display

import (
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CountBadge shows a rolling count over a child's corner.
func CountBadge(c *ui.Context, count int, content func()) ui.Element {
	if count < 0 {
		panic("display: CountBadge count must be non-negative")
	}
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	return root.Children(func() {
		content()
		text := strconv.Itoa(count)
		if count > 99 {
			text = "99+"
		}
		badge := ui.Row(c).Absolute().Top(-t.Space(1.5)).Right(-t.Space(1.5)).Padding(t.Space(0.5), t.Space(1)).Radius(999).
			Background(t.Danger).TextColor(t.AccentText).Border(1, t.Background)
		badge.Children(func() { kit.Numeric(c, text).FontSize(t.Rem(0.75)) })
	})
}
