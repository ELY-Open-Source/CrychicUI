package layout

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Reveal rises and fades content in once it shows.
func Reveal(c *ui.Context, content func()) ui.Element {
	root := ui.Box(c)
	shown := ui.Local(root, "shown", func() bool { return false })
	if root.InView() {
		*shown = true
	}
	if !*shown {
		kit.Settle(c, root, []float32{0}, func() []float32 { return []float32{unit(root.InView())} })
	}
	return root.Children(func() {
		b := ui.Box(c)
		p := b.AnimateTo("hidden", unit(!*shown), ui.Smooth)
		b.Translate(0, 12*p).Opacity(1 - p).Children(content)
	})
}
