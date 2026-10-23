// Code components for the gallery.
package code

import (
	"fmt"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

type PackageCardOptions struct {
	Copy     string
	Install  string
	Download string
}

// PackageCard is a package with install copy; cue Pop.

func PackageCard(c *ui.Context, name, version string, downloads int, o PackageCardOptions) (root, copy ui.Element) {
	if name == "" || version == "" || downloads < 0 {
		panic("code: PackageCard needs a name, version and downloads")
	}
	def(&o.Copy, "Copy install")
	def(&o.Install, "go get %s@%s")
	def(&o.Download, "%s downloads")
	t := c.Theme()
	root = kit.Card(c).AlignSelf(ui.Start).Width(t.Space(76)).Label(name)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	cmd := fmt.Sprintf(o.Install, name, version)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, iconPackage).Size(16, 16).TextColor(t.Accent)
			ui.Text(c, name).Bold().Grow(1).SingleLine()
			ui.Text(c, version).Font("monospace").FontSize(t.Rem(0.8)).TextColor(t.TextMuted).
				Padding(t.Space(0.5), t.Space(1.5)).Border(1, t.Border).Radius(999)
		})
		kit.Muted(c, fmt.Sprintf(o.Download, compact(downloads))).FontSize(t.Rem(0.8))
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(2)).
			Background(t.Surface).Radius(t.Radius).Children(func() {
			ui.Text(c, cmd).Font("monospace").FontSize(t.Rem(0.8)).SingleLine().Grow(1)
			copy = copyButton(c, o.Copy, cmd)
		})
	})
	return root, copy
}
