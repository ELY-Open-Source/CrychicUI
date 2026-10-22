package color

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// HarmonyOptions tune a Harmony.
type HarmonyOptions struct {
	// Names title the complement, triad and analogous rows.
	Names [3]string
}

// Harmony offers base's complement, triad and analogous colors.
func Harmony(c *ui.Context, label string, base ui.Color, value *ui.Color, o HarmonyOptions) ui.Element {
	if label == "" || value == nil {
		panic("color: Harmony needs a label and a color")
	}
	if o.Names == ([3]string{}) {
		o.Names = [3]string{"Complement", "Triad", "Analogous"}
	}
	t := c.Theme()
	l, ch, h := oklch(base)
	turn := func(by float64) ui.Color { return ui.Oklch(float32(l), float32(ch), float32(h+by)).SRGB() }
	groups := [3][]ui.Color{{turn(180)}, {turn(120), turn(240)}, {turn(-30), turn(30)}}
	size := t.Space(9)
	root := ui.Column(c).AlignSelf(ui.Start).Gap(t.Space(2)).Label(label)
	return root.Children(func() {
		for g, cols := range groups {
			c.Key(g)
			ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
				kit.Muted(c, o.Names[g]).Width(t.Space(24))
				swatch(c, base, size, size)
				for i, col := range cols {
					c.Key(i)
					b := kit.Press(ui.ButtonBase(c).Label(hexOf(col)).Tooltip(hexOf(col)).Size(size, size).Padding(t.Space(0.5)).Radius(t.Radius))
					ring := t.Accent.Alpha(0)
					if col == *value {
						ring = t.Accent
					}
					b.Border(2, b.AnimateColor("ring", ring, ui.Snappy)).Children(func() {
						swatch(c, col, size-t.Space(1)-4, size-t.Space(1)-4)
					})
					if kit.Clicked(b, sound.Select) {
						*value = col
						root.SetChanged()
						kit.Log("Harmony", "picked", "row", g, "index", i)
					}
				}
			})
		}
	})
}
