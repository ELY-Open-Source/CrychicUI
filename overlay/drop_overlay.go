package overlay

import (
	"slices"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DropOverlayOptions tune a DropOverlay.
type DropOverlayOptions struct {
	// Text shows on the veil.
	Text string
}

// DropOverlay veils content while files drag over it.
func DropOverlay(c *ui.Context, files *[]string, o DropOverlayOptions, content func()) ui.Element {
	root := ui.Box(c)
	return dropOverlay(c, root, files, o, content, root.FileDragOver(), root.DroppedFiles())
}

// dropOverlay builds DropOverlay from what the window reports.
func dropOverlay(c *ui.Context, root ui.Element, files *[]string, o DropOverlayOptions, content func(), over bool, got []string) ui.Element {
	if files == nil || content == nil {
		panic("overlay: DropOverlay needs files and content")
	}
	def(&o.Text, "Drop files here")
	t := c.Theme()
	if len(got) > 0 {
		*files = slices.Clone(got)
		root.SetChanged()
		sound.Play(sound.Drop)
		kit.Log("DropOverlay", "dropped", "count", len(got))
	}
	cover(c, root, over, 0, content, func() {
		ui.Column(c).FillWidth().FillHeight().Center().Gap(t.Space(2)).Padding(t.Space(4)).Radius(t.Radius).
			Background(t.Background.Mix(t.Accent, 0.1)).Border(1, t.Accent).Children(func() {
			up := ui.Icon(c, icon.Must("upload")).Size(16, 16).TextColor(t.Accent)
			up.Translate(0, -t.Space(1)*up.Loop("bob", 1200*time.Millisecond, ui.Bounce(ui.EaseInOut)))
			ui.Text(c, o.Text).Bold().TextColor(t.Accent)
		})
	})
	moved(root, over, "DropOverlay")
	return root
}
