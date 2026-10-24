package shop

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SizeRow is one size guide row.
type SizeRow struct {
	Size, Chest, Waist string
}

// SizeGuideOptions word the button and column heads.
type SizeGuideOptions struct{ Guide, Size, Chest, Waist string }

// SizeGuide opens a striped size table popover; cue Open.
func SizeGuide(c *ui.Context, rows []SizeRow, open *bool, o SizeGuideOptions) ui.Element {
	if len(rows) == 0 || open == nil {
		panic("shop: SizeGuide needs rows and open state")
	}
	for _, r := range rows {
		if r.Size == "" || r.Chest == "" || r.Waist == "" {
			panic(fmt.Sprintf("shop: SizeGuide row %q needs every cell", r.Size))
		}
	}
	def(&o.Guide, "Size guide")
	def(&o.Size, "Size")
	def(&o.Chest, "Chest")
	def(&o.Waist, "Waist")
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	root.Children(func() {
		b := kit.Press(ui.ButtonBase(c).Label(o.Guide).Radius(999).Expanded(*open))
		b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered() || *open)), ui.Snappy))
		b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(*open)), ui.Snappy))
		b.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Padding(t.Space(1.5), t.Space(3)).Children(func() {
				ui.Icon(c, icon.Must("ruler")).Size(14, 14).TextColor(t.TextMuted)
				ui.Text(c, o.Guide).FontSize(t.Rem(0.85)).Shrink(0)
			})
		})
		if b.Clicked() {
			*open = !*open
			root.SetChanged()
			if *open {
				sound.Play(sound.Open)
			} else {
				sound.Play(sound.Close)
			}
			kit.Log("SizeGuide", "toggled", "open", *open)
		}
		if ui.PopoverBase(c, b, open, func(panel ui.Element) {
			panel.Margin(t.Space(1), 0, 0, 0).Width(t.Space(72)).Padding(t.Space(2)).Radius(t.Radius).
				Background(t.Background).Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.16)).TransformOrigin(0.5, 0)
			kit.Arrive(panel, ui.Pose{Scale: 0.95, Y: -t.Space(1), Fade: 1}, 0)
			cells := func(bg ui.Color, bold bool, values ...string) ui.Element {
				return ui.Row(c).FillWidth().Radius(t.Space(1)).Background(bg).Children(func() {
					for i, v := range values {
						c.Key(i)
						cell := ui.Text(c, v).FontSize(t.Rem(0.8)).Grow(1).Basis(0).MinWidth(0).Padding(t.Space(1.5), t.Space(2))
						if bold {
							cell.Bold().TextColor(t.TextMuted)
						}
					}
				})
			}
			c.Key("head")
			cells(t.Surface, true, o.Size, o.Chest, o.Waist)
			list := ui.Scroll(c).Key("rows").FillWidth().MaxHeight(t.Space(64))
			list.Children(func() {
				for i, r := range rows {
					c.Key(i)
					kit.Stagger(c, list, cells(t.Surface.Alpha(0.5*unit(i%2 == 1)), false, r.Size, r.Chest, r.Waist), i)
				}
			})
		}).Dismissed() {
			root.SetChanged()
			sound.Play(sound.Close)
			kit.Log("SizeGuide", "dismissed")
		}
	})
	return root
}
