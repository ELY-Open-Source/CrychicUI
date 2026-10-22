package color

import (
	"slices"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// RecentColorsOptions tune RecentColors.
type RecentColorsOptions struct {
	// Max caps the swatches shown; 8 when zero.
	Max int
	// Empty shows before any color is chosen.
	Empty string
}

// recentKey tells apart repeated colors by their order.
type recentKey struct {
	col ui.Color
	n   int
}

// RecentColors shows colors newest first; a new one slides in.
func RecentColors(c *ui.Context, label string, colors []ui.Color, value *ui.Color, o RecentColorsOptions) ui.Element {
	if o.Max == 0 {
		o.Max = 8
	}
	if label == "" || value == nil || o.Max < 1 {
		panic("color: RecentColors needs a label, a color and Max >= 1")
	}
	def(&o.Empty, "No recent colors")
	t := c.Theme()
	size := t.Space(8)
	shown := colors[:min(len(colors), o.Max)]
	slide := ui.ElementTransition{Enter: &ui.Motion{X: -size}, Exit: &ui.Motion{Scale: 0.6}}
	root := ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1.5)).MinHeight(size).AlignItems(ui.Center).Label(label)
	return root.Children(func() {
		if len(shown) == 0 {
			kit.Muted(c, o.Empty).Key("empty").Transition(kit.Enter)
			return
		}
		for i, col := range shown {
			c.Key(recentKey{col, i - slices.Index(shown, col)})
			b := kit.Press(ui.ButtonBase(c).Label(hexOf(col)).Tooltip(hexOf(col)).Size(size, size).Padding(t.Space(0.5)).Radius(t.Radius))
			ring := t.Accent.Alpha(0)
			if col == *value {
				ring = t.Accent
			}
			b.Border(2, b.AnimateColor("ring", ring, ui.Snappy)).Transition(slide).Children(func() {
				swatch(c, col, size-t.Space(1)-4, size-t.Space(1)-4)
			})
			if kit.Clicked(b, sound.Select) {
				*value = col
				root.SetChanged()
				kit.Log("RecentColors", "picked", "index", i)
			}
		}
	})
}
