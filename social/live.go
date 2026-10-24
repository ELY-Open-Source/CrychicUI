package social

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LiveData is one live stream card.
type LiveData struct {
	ID, Title, Streamer string
	Viewers             int
	Cover               ui.ImageSource
}

// LiveOptions are the live card labels.
type LiveOptions struct{ Live, Viewers string }

// Live shows a live stream card with viewers.
func Live(c *ui.Context, live LiveData, o LiveOptions) ui.Element {
	if live.ID == "" || live.Title == "" || live.Streamer == "" || live.Viewers < 0 || live.Cover == nil {
		panic("social: invalid live stream")
	}
	def(&o.Live, "LIVE")
	def(&o.Viewers, "%d viewers")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2)).Radius(t.Radius).Border(1, t.Border)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	root.Children(func() {
		ui.Box(c).FillWidth().AspectRatio(1.8).Clip().Children(func() {
			ui.Image(c, live.Cover).Fill().Fit(ui.Cover)
			ui.Text(c, o.Live).Absolute().Left(t.Space(2)).Top(t.Space(2)).Padding(t.Space(0.5), t.Space(1)).
				Radius(t.Space(0.5)).Background(t.Danger).TextColor(t.AccentText)
		})
		ui.Column(c).Gap(t.Space(1)).Padding(t.Space(2)).Children(func() {
			ui.Text(c, live.Title).Bold()
			ui.Row(c).Children(func() {
				kit.Muted(c, live.Streamer).Grow(1)
				kit.Numeric(c, fmt.Sprintf(o.Viewers, live.Viewers)).TextColor(t.TextMuted)
			})
		})
	})
	return root
}
