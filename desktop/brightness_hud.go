package desktop

import "github.com/egoist/mygo/ui"

// BrightnessHUDOptions name the overlay.
type BrightnessHUDOptions struct{ Label string }

// BrightnessHUD shows brightness, fading after a change.
func BrightnessHUD(c *ui.Context, level float32, o BrightnessHUDOptions) ui.Element {
	def(&o.Label, "Brightness")
	share("BrightnessHUD", o.Label, level)
	mark := "sun-medium"
	if level < 0.5 {
		mark = "sun"
	}
	return hud(c, o.Label, mark, level)
}
