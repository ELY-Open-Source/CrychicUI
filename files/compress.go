package files

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CompressOptions word the change, "%s" taking the size.
type CompressOptions struct {
	Smaller, Larger, Before, After string
}

// Compress shows sizes before and after compression as two bars.
func Compress(c *ui.Context, label string, before, after int64, o CompressOptions) ui.Element {
	if label == "" || before < 0 || after < 0 {
		panic(fmt.Sprintf("files: Compress %q needs sizes not below zero", label))
	}
	def(&o.Smaller, "%s smaller")
	def(&o.Larger, "%s larger")
	def(&o.Before, "Before")
	def(&o.After, "After")
	t := c.Theme()
	format, delta, ink := o.Smaller, before-after, t.Success
	if after > before {
		format, delta, ink = o.Larger, after-before, t.Warning
	}
	biggest := max(before, after, 1)
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	kit.Arrive(root, ui.Pose{Y: t.Space(1), Fade: 1}, 0)
	root.Children(func() {
		kit.CardHead(c, "minimize-2", label, func() {
			kit.Chip(c, fmt.Sprintf(format, sizeLabel(delta)), ink)
		})
		ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3)).Children(func() {
			bar := func(key, name string, bytes int64, ink ui.Color) {
				c.Key(key)
				ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
					kit.Muted(c, name).FontSize(t.Rem(0.8)).Width(t.Space(14)).Shrink(0)
					ui.Box(c).Grow(1).Children(func() { progressBar(c, float64(bytes)/float64(biggest), ink) })
					sizeText(c, bytes).Width(t.Space(16)).Shrink(0).TextAlign(ui.End)
				})
			}
			bar("before", o.Before, before, t.TextMuted)
			bar("after", o.After, after, ink)
		})
	})
	return root
}

// sizeLabel is a byte size in decimal units, as Bytes.
func sizeLabel(n int64) string {
	units := [...]string{"B", "KB", "MB", "GB", "TB"}
	v, u := float64(n), 0
	for v >= 1000 && u < len(units)-1 {
		v /= 1000
		u++
	}
	if u == 0 {
		return fmt.Sprintf("%d B", n)
	}
	return fmt.Sprintf("%.1f %s", v, units[u])
}
