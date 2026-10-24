package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ShareCountOptions are the share count labels.
type ShareCountOptions struct{ Label string }

// ShareCount rolls a share count.
func ShareCount(c *ui.Context, count int, o ShareCountOptions) ui.Element {
	if count < 0 {
		panic("social: ShareCount must not be negative")
	}
	def(&o.Label, "%d shares")
	t := c.Theme()
	root := ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center)
	kit.Arrive(root, ui.Pose{Scale: 0.9, Fade: 1}, 0)
	root.Children(func() {
		ui.Icon(c, icon.Must("share-2")).Size(16, 16).TextColor(t.Accent)
		kit.Numeric(c, fmt.Sprintf(o.Label, count))
	})
	return root
}
