package marketing

import (
	"fmt"
	"strconv"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// Bento is one tile; Wide spans two columns.
type Bento struct {
	Icon, Title, Text string
	Wide              bool
	Media             *ui.Bitmap
}

// BentoFeaturesOptions word the section.
type BentoFeaturesOptions struct{ Title, Text string }

// BentoFeatures lays features out in tiles of two widths.
func BentoFeatures(c *ui.Context, tiles []Bento, o BentoFeaturesOptions) ui.Element {
	if len(tiles) == 0 {
		panic(fmt.Sprintf("marketing: BentoFeatures %q needs tiles", o.Title))
	}
	titles := make([]string, len(tiles))
	for i, b := range tiles {
		if b.Icon == "" {
			panic(fmt.Sprintf("marketing: BentoFeatures tile %q needs an icon", b.Title))
		}
		titles[i] = b.Title
	}
	unique("BentoFeatures", titles)
	t := c.Theme()
	inks := [4]ui.Color{t.Accent, t.Success, t.Warning, t.Danger}
	root := ui.Column(c).FillWidth().Gap(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		heading(c, o.Title, o.Text)
		grid := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(3))
		grid.Children(func() {
			for i, b := range tiles {
				ink := inks[i%4]
				c.Key(strconv.Quote(b.Title))
				cell := ui.Column(c).MinWidth(0).Gap(t.Space(2)).Radius(t.Radius*1.5).Clip().Border(1, t.Border)
				if b.Wide {
					cell.ColumnSpan(2)
				}
				cell.Background(cell.AnimateColor("bg", t.Background.Mix(ink, 0.05+0.05*unit(cell.Hovered())), ui.Snappy))
				kit.Stagger(c, grid, cell, i).Children(func() {
					if b.Media != nil {
						ui.Box(c).FillWidth().AspectRatio(3).Clip().Children(func() {
							kit.Arrive(ui.Image(c, b.Media).Fill().Fit(ui.Cover), ui.Pose{Scale: 1.06, Fade: 1}, 0)
						})
					}
					ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(4)).Children(func() {
						tile(c, b.Icon, ink, t.Space(9))
						ui.Text(c, b.Title).FontSize(t.Rem(0.92)).Bold().FillWidth().MinWidth(0).Wrap()
						kit.Muted(c, b.Text).FontSize(t.Rem(0.8)).FillWidth().MinWidth(0).Wrap()
					})
				})
			}
		})
	})
	return root
}
