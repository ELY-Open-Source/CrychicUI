package overlay

import (
	"fmt"
	"math"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// MiniPlayerOptions tune a MiniPlayer.
type MiniPlayerOptions struct {
	Play, Pause string
}

// MiniPlayer is a floating play bar with progress.
func MiniPlayer(c *ui.Context, title, artist string, playing *bool, progress float32, o MiniPlayerOptions) ui.Element {
	if title == "" || playing == nil || math.IsNaN(float64(progress)) || progress < 0 || progress > 1 {
		panic(fmt.Sprintf("overlay: MiniPlayer %q needs playing and progress in 0..1", title))
	}
	def(&o.Play, "Play")
	def(&o.Pause, "Pause")
	t := c.Theme()
	root := ui.Column(c).FillWidth().MaxWidth(t.Space(110)).Radius(t.Radius*1.5).Background(t.Background).Border(1, t.Border).
		Shadow(0, 8, 24, 0, kit.Shade(0.16)).Clip().Role(ui.RoleGroup).Label(title)
	kit.Arrive(root, ui.Pose{Y: t.Space(6), Fade: 1}, 0).Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Children(func() {
			ui.Box(c).Size(t.Space(10), t.Space(10)).Center().Radius(t.Radius).Background(t.Accent.Alpha(0.12)).Children(func() {
				ui.Icon(c, icon.Must("music")).Size(16, 16).TextColor(t.Accent)
			})
			ui.Column(c).Grow(1).Children(func() {
				ui.Text(c, title).Bold().SingleLine()
				if artist != "" {
					kit.Muted(c, artist).FontSize(t.Rem(0.8)).SingleLine()
				}
			})
			name, glyph := o.Play, "play"
			if *playing {
				name, glyph = o.Pause, "pause"
			}
			b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(10), t.Space(10)).Center().Radius(999))
			b.Background(b.AnimateColor("bg", t.Accent.Mix(t.AccentText, 0.12*unit(b.Hovered())), ui.Snappy)).TextColor(t.AccentText).
				Children(func() {
					ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
				})
			if b.Clicked() {
				*playing = !*playing
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("MiniPlayer", "playing", *playing)
			}
		})
		ui.Box(c).FillWidth().Height(t.Space(0.5)).Background(t.Border).Children(func() {
			bar := ui.Box(c).FillHeight().Background(t.Accent)
			bar.WidthPercent(100 * bar.AnimateTo("p", progress, ui.Smooth))
		})
	})
	return root
}
