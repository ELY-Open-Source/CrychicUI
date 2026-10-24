package mail

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ScreenShareOptions are the screen sharing labels.
type ScreenShareOptions struct{ Sharing, Stop string }

// ScreenShare stops the active share; cue Close.
func ScreenShare(c *ui.Context, sharing *bool, o ScreenShareOptions) ui.Element {
	if sharing == nil {
		panic("mail: ScreenShare needs state")
	}
	def(&o.Sharing, "You are sharing")
	def(&o.Stop, "Stop")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2)).Radius(t.Radius).Background(t.Accent)
	root.Opacity(root.AnimateTo("sharing", float32(unit(*sharing)), ui.Snappy))
	root.Children(func() {
		ui.Text(c, o.Sharing).TextColor(t.AccentText).Grow(1)
		stop := button(c, o.Stop).Background(t.Background).Disabled(!*sharing)
		stop.Children(func() { ui.Text(c, o.Stop) })
		if stop.Clicked() && *sharing {
			*sharing = false
			root.SetChanged()
			sound.Play(sound.Close)
			kit.Log("ScreenShare", "stopped")
		}
	})
	return root
}
