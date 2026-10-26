package desktop

import "github.com/egoist/mygo/ui"

// VolumeHUDOptions mute the sound and name the overlay.
type VolumeHUDOptions struct {
	Muted bool
	Label string
}

// VolumeHUD shows the volume, fading after a change.
func VolumeHUD(c *ui.Context, level float32, o VolumeHUDOptions) ui.Element {
	def(&o.Label, "Volume")
	share("VolumeHUD", o.Label, level)
	mark := "volume-2"
	switch {
	case o.Muted || level == 0:
		mark, level = "volume-x", 0
	case level < 0.5:
		mark = "volume-1"
	}
	return hud(c, o.Label, mark, level)
}
