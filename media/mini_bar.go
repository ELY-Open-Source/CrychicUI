package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// MiniBarOptions tune a MiniBar.
type MiniBarOptions struct {
	Play string
}

// MiniBar is a compact player row with progress.
func MiniBar(c *ui.Context, track Track, playing *bool, pos time.Duration, o MiniBarOptions) ui.Element {
	tracked("MiniBar", track)
	if playing == nil || pos < 0 || pos > track.Length {
		panic(fmt.Sprintf("media: MiniBar %q needs state and a position within the track", track.Title))
	}
	def(&o.Play, "Play")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Radius(t.Radius).Clip().Background(t.Surface).Border(1, t.Border).Label(track.Title)
	return kit.Arrive(root, ui.Pose{Fade: 1, Y: t.Space(3)}, 0).Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2)).Children(func() {
			cover(c, track.Cover, t.Space(10))
			ui.Column(c).Grow(1).Children(func() {
				ui.Text(c, track.Title).SingleLine().Bold().ContentTransition(ui.ContentTransition{})
				kit.Muted(c, track.Artist).SingleLine().FontSize(t.Rem(0.8))
			})
			if button.PlayButton(c, playing, button.PlayButtonOptions{Label: o.Play}).AlignSelf(ui.Center).Changed() {
				root.SetChanged()
			}
		})
		line := ui.Box(c).Height(t.Space(0.5)).Background(t.Accent)
		line.WidthPercent(100 * line.AnimateTo("done", float32(pos.Seconds()/track.Length.Seconds()), ui.Snappy))
	})
}
