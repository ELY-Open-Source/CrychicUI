package code

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Package is a published package.
type Package struct {
	Name, Version, Description, License string
	Downloads                           int
}

// PackageCardOptions are the copy knobs of PackageCard.
type PackageCardOptions struct {
	Copy string
	// Install words the command, "%s" taking name then version.
	Install string
	// Download words the count, "%s" taking it compacted.
	Download string
}

// PackageCard shows a package, stats and install copy; cue Pop.
func PackageCard(c *ui.Context, p Package, o PackageCardOptions) (root, copy ui.Element) {
	if p.Name == "" || p.Version == "" || p.Downloads < 0 {
		panic(fmt.Sprintf("code: PackageCard %q needs a name, a version and downloads", p.Name))
	}
	def(&o.Copy, "Copy install")
	def(&o.Install, "go get %s@%s")
	def(&o.Download, "%s downloads")
	t := c.Theme()
	cmd := fmt.Sprintf(o.Install, p.Name, p.Version)
	root = kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3)).Label(p.Name)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).Children(func() {
			ui.Box(c).Size(t.Space(10), t.Space(10)).Shrink(0).Center().Radius(t.Radius).Background(t.Accent.Alpha(0.12)).TextColor(t.Accent).
				Children(func() { ui.Icon(c, iconPackage).Size(20, 20) })
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					ui.Text(c, p.Name).Font("monospace").FontSize(t.Rem(0.9)).Bold().Grow(1).MinWidth(0).SingleLine()
					chip(c, p.Version, t.Accent)
				})
				if p.Description != "" {
					kit.Muted(c, p.Description).FontSize(t.Rem(0.85)).Wrap()
				}
			})
		})
		ui.Row(c).Gap(t.Space(4)).AlignItems(ui.Center).Children(func() {
			stat := func(glyph, text string) {
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(t.TextMuted).Children(func() {
					ui.Icon(c, icon.Must(glyph)).Size(14, 14)
					ui.Text(c, text).FontSize(t.Rem(0.8)).Shrink(0)
				})
			}
			stat("download", fmt.Sprintf(o.Download, compact(p.Downloads)))
			if p.License != "" {
				stat("scale", p.License)
			}
		})
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1), t.Space(1), t.Space(1), t.Space(3)).
			Background(t.Surface).Border(1, t.Border).Radius(t.Radius).Children(func() {
			ui.Text(c, "$").Font("monospace").FontSize(t.Rem(0.8)).Bold().TextColor(t.Success)
			ui.Text(c, cmd).Font("monospace").FontSize(t.Rem(0.8)).Grow(1).MinWidth(0).SingleLine()
			copy = copyButton(c, o.Copy, cmd)
		})
	})
	if copy.Clicked() {
		kit.Log("PackageCard", "copied")
	}
	return root, copy
}
