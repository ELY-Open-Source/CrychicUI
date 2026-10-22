package color

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PaletteOptions tune a Palette.
type PaletteOptions struct {
	// Steps counts shades from light to dark; 9 when zero.
	Steps int
}

// Palette shows shades of base's hue; they grow in.
func Palette(c *ui.Context, base ui.Color, o PaletteOptions) ui.Element {
	if o.Steps == 0 {
		o.Steps = 9
	}
	if o.Steps < 2 || o.Steps > 20 {
		panic(fmt.Sprintf("color: Palette needs 2 to 20 steps, not %d", o.Steps))
	}
	t := c.Theme()
	_, ch, h := oklch(base)
	return ui.Row(c).AlignSelf(ui.Start).Gap(t.Space(1)).Role(ui.RoleList).Label(hexOf(base)).Children(func() {
		for i := range o.Steps {
			l := 0.95 - 0.7*float64(i)/float64(o.Steps-1)
			shade := ui.Oklch(float32(l), float32(ch), float32(h)).SRGB()
			c.Key(i)
			kit.Arrive(swatch(c, shade, t.Space(9), t.Space(12)), ui.Pose{Scale: 0.6, Fade: 1}, time.Duration(i)*40*time.Millisecond)
		}
	})
}
