package social

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LiveData is one live stream card.
type LiveData struct {
	ID, Title, Streamer string
	Viewers             int
	Cover               ui.ImageSource
}

// LiveOptions word the badge; Viewers formats the grouped count.
type LiveOptions struct{ Live, Viewers string }

// Live shows a live stream card with viewers.
func Live(c *ui.Context, live LiveData, o LiveOptions) ui.Element {
	if live.ID == "" || live.Title == "" || live.Streamer == "" || live.Viewers < 0 || live.Cover == nil {
		panic(fmt.Sprintf("social: Live %q needs id, streamer, viewers and cover", live.Title))
	}
	def(&o.Live, "LIVE")
	def(&o.Viewers, "%s watching")
	t := c.Theme()
	root := card(c)
	kit.Lift(root)
	root.Children(func() {
		ui.Box(c).FillWidth().AspectRatio(16.0 / 9).Clip().Background(t.SurfaceHover).Children(func() {
			ui.Image(c, live.Cover).Fill().Fit(ui.Cover)
			ui.Row(c).Absolute().Left(t.Space(2.5)).Top(t.Space(2.5)).Gap(t.Space(1.5)).AlignItems(ui.Center).
				Padding(t.Space(0.5), t.Space(2)).Radius(t.Radius).Background(t.Danger).TextColor(t.AccentText).Children(func() {
				dot := ui.Box(c).Size(t.Space(1.5), t.Space(1.5)).Shrink(0).Radius(999).Background(t.AccentText)
				beat := dot.Loop("beat", 1200*time.Millisecond, ui.EaseInOut)
				dot.Opacity(0.35 + 0.65*(1-beat))
				ui.Text(c, o.Live).FontSize(t.Rem(0.7)).Bold().Shrink(0)
			})
			ui.Row(c).Absolute().Right(t.Space(2.5)).Top(t.Space(2.5)).Gap(t.Space(1)).AlignItems(ui.Center).
				Padding(t.Space(0.5), t.Space(2)).Radius(t.Radius).Background(t.Text.Alpha(0.55)).TextColor(t.Background).Children(func() {
				glyph(c, "eye", 12)
				kit.Numeric(c, fmt.Sprintf(o.Viewers, number(c, live.Viewers))).FontSize(t.Rem(0.7)).Bold().Shrink(0)
			})
		})
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(3)).Children(func() {
			ui.Box(c).Shrink(0).Padding(t.Space(0.5)).Radius(999).Border(t.Space(0.5), t.Danger).Children(func() {
				avatar(c, User{ID: live.Streamer, Name: live.Streamer}, t.Space(9))
			})
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
				ui.Text(c, live.Title).FontSize(t.Rem(0.9)).Bold().FillWidth().MinWidth(0).SingleLine()
				kit.Muted(c, live.Streamer).FontSize(t.Rem(0.75)).FillWidth().MinWidth(0).SingleLine()
			})
		})
	})
	return root
}
