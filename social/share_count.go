package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ShareCountOptions: One and Label format the grouped count.
type ShareCountOptions struct{ One, Label string }

// ShareCount rolls a share count.
func ShareCount(c *ui.Context, count int, o ShareCountOptions) ui.Element {
	if count < 0 {
		panic("social: ShareCount must not be negative")
	}
	def(&o.One, "%s share")
	def(&o.Label, "%s shares")
	format := o.Label
	if count == 1 {
		format = o.One
	}
	t := c.Theme()
	root := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).AlignSelf(ui.Start).Padding(t.Space(1), t.Space(3), t.Space(1), t.Space(1)).
		Radius(999).Background(t.Accent.Alpha(0.08))
	kit.Arrive(root, ui.Pose{Scale: 0.9, Fade: 1}, 0)
	root.Children(func() {
		ui.Box(c).Size(t.Space(7), t.Space(7)).Shrink(0).Center().Radius(999).Background(t.Accent).TextColor(t.AccentText).
			Children(func() { glyph(c, "share-2", 14) })
		kit.Numeric(c, fmt.Sprintf(format, number(c, count))).FontSize(t.Rem(0.85)).Bold().TextColor(t.Accent).Shrink(0)
	})
	return root
}
