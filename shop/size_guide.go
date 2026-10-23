package shop

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SizeRow is one size guide row.
type SizeRow struct {
	Size, Chest, Waist string
}

// SizeGuideOptions are the size guide labels.
type SizeGuideOptions struct{ Guide, Size, Chest, Waist string }

// SizeGuide shows a size table in a popover; cue Open.
func SizeGuide(c *ui.Context, rows []SizeRow, open *bool, o SizeGuideOptions) ui.Element {
	if len(rows) == 0 || open == nil {
		panic("shop: SizeGuide needs rows and state")
	}
	def(&o.Guide, "Size guide")
	def(&o.Size, "Size")
	def(&o.Chest, "Chest")
	def(&o.Waist, "Waist")
	t := c.Theme()
	root := ui.Box(c).AlignSelf(ui.Start)
	root.Children(func() {
		button := kit.Press(ui.ButtonBase(c).Label(o.Guide).Padding(t.Space(1), t.Space(2)).Radius(t.Radius).
			Border(1, t.Border).Checked(*open))
		button.Children(func() { ui.Text(c, o.Guide).SingleLine() })
		if button.Clicked() {
			*open = !*open
			root.SetChanged()
			if *open {
				sound.Play(sound.Open)
			} else {
				sound.Play(sound.Close)
			}
			kit.Log("SizeGuide", "toggled", "open", *open)
		}
		if ui.PopoverBase(c, button, open, func(panel ui.Element) {
			panel.Padding(t.Space(3)).Radius(t.Radius).Border(1, t.Border).
				Background(t.Surface).Shadow(0, t.Space(2), t.Space(6), 0, kit.Shade(0.15)).Transition(kit.Enter)
			panel.Children(func() {
				ui.Grid(c).Columns(3).Gap(t.Space(1)).Children(func() {
					ui.Text(c, o.Size).Bold()
					ui.Text(c, o.Chest).Bold()
					ui.Text(c, o.Waist).Bold()
					for _, row := range rows {
						if row.Size == "" || row.Chest == "" || row.Waist == "" {
							panic("shop: SizeGuide row must not be empty")
						}
						ui.Text(c, row.Size)
						ui.Text(c, row.Chest)
						ui.Text(c, row.Waist)
					}
				})
			})
		}).Dismissed() {
			root.SetChanged()
			sound.Play(sound.Close)
			kit.Log("SizeGuide", "dismissed")
		}
	})
	return root
}
