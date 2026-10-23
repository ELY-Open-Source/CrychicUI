package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CompressOptions are the size change labels.
type CompressOptions struct{ Smaller, Larger string }

// Compress shows sizes before and after compression.
func Compress(c *ui.Context, label string, before, after int64, o CompressOptions) ui.Element {
	if label == "" || before < 0 || after < 0 {
		panic("files: Compress needs nonnegative sizes")
	}
	def(&o.Smaller, "%s smaller")
	def(&o.Larger, "%s larger")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).
		Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		format, delta, ink := o.Smaller, before-after, t.Success
		if after > before {
			format, delta, ink = o.Larger, after-before, t.Warning
		}
		ui.Row(c).FillWidth().AlignItems(ui.Center).Children(func() {
			ui.Text(c, sizeLabel(before)).Grow(1)
			kit.Numeric(c, fmt.Sprintf(format, sizeLabel(delta))).TextColor(ink)
			ui.Text(c, sizeLabel(after))
		})
		fraction := 0.0
		if before > 0 {
			fraction = min(1, float64(after)/float64(before))
		} else if after > 0 {
			fraction = 1
		}
		progressBar(c, fraction, t.Accent)
	})
	return root
}
