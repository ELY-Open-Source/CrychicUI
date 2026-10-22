package media

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// PodcastCardOptions tune a PodcastCard.
type PodcastCardOptions struct {
	// Minutes writes the length; "%d min" by default.
	Minutes, Play string
}

// PodcastCard shows an episode with its show, length and play.
func PodcastCard(c *ui.Context, episode Track, playing *bool, o PodcastCardOptions) ui.Element {
	tracked("PodcastCard", episode)
	if playing == nil {
		panic(fmt.Sprintf("media: PodcastCard %q needs state", episode.Title))
	}
	def(&o.Minutes, "%d min")
	def(&o.Play, "Play")
	t := c.Theme()
	root := kit.Lift(kit.Card(c).Row().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(3)).Label(episode.Title))
	return kit.Arrive(root, ui.Pose{Fade: 1, Y: t.Space(2)}, 0).Children(func() {
		cover(c, episode.Cover, t.Space(18))
		ui.Column(c).Grow(1).Gap(t.Space(1)).Children(func() {
			kit.Muted(c, episode.Artist).SingleLine().FontSize(t.Rem(0.75)).Bold()
			ui.Text(c, episode.Title).Bold().MaxLines(2)
			kit.Muted(c, fmt.Sprintf(o.Minutes, int(episode.Length.Minutes()+0.5))).FontSize(t.Rem(0.8))
		})
		if button.PlayButton(c, playing, button.PlayButtonOptions{Label: o.Play}).AlignSelf(ui.Center).Changed() {
			root.SetChanged()
		}
	})
}
