package media

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/feedback"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// NowPlayingOptions tune a NowPlaying.
type NowPlayingOptions struct{}

// NowPlaying shows a track's cover, title, artist and progress.
func NowPlaying(c *ui.Context, track Track, pos time.Duration, o NowPlayingOptions) ui.Element {
	tracked("NowPlaying", track)
	if pos < 0 || pos > track.Length {
		panic(fmt.Sprintf("media: NowPlaying %q needs a position within the track", track.Title))
	}
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(4)).AlignItems(ui.Center).Padding(t.Space(4)).Radius(t.Radius).
		Border(1, t.Border).Label(track.Title)
	return kit.Arrive(root, ui.Pose{Fade: 1, Y: t.Space(2)}, 0).Children(func() {
		c.Key(track.Title)
		kit.Arrive(cover(c, track.Cover, t.Space(20)), ui.Pose{Fade: 1, Scale: 0.9}, 0)
		ui.Column(c).Grow(1).Gap(t.Space(1)).Children(func() {
			ui.Text(c, track.Title).SingleLine().Bold().FontSize(t.Rem(1.1)).ContentTransition(ui.ContentTransition{})
			kit.Muted(c, track.Artist).SingleLine().ContentTransition(ui.ContentTransition{})
			feedback.ProgressBar(c, track.Title, pos.Seconds()/track.Length.Seconds()).Margin(t.Space(1), 0, 0, 0)
			ui.Row(c).FillWidth().Justify(ui.SpaceBetween).Children(func() {
				kit.Numeric(c, clock(pos)).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
				kit.Numeric(c, clock(track.Length)).FontSize(t.Rem(0.75)).TextColor(t.TextMuted)
			})
		})
	})
}
