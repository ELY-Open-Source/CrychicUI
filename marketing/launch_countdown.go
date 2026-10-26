package marketing

import (
	"fmt"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// LaunchCountdownOptions word the countdown; Units name days to seconds.
type LaunchCountdownOptions struct {
	Title, Text, Live string
	Units             [4]string
}

// LaunchCountdown counts down to a launch.
func LaunchCountdown(c *ui.Context, launch time.Time, o LaunchCountdownOptions) ui.Element {
	def(&o.Live, "We're live")
	for i, d := range [4]string{"Days", "Hours", "Minutes", "Seconds"} {
		def(&o.Units[i], d)
	}
	if launch.IsZero() {
		panic(fmt.Sprintf("marketing: LaunchCountdown %q needs a launch time", o.Title))
	}
	t := c.Theme()
	secs := countdown(c, launch)
	root := card(c).Gap(t.Space(4)).Padding(t.Space(6), t.Space(5)).AlignItems(ui.Center).Role(ui.RoleGroup)
	root.Children(func() {
		kit.Arrive(tile(c, "rocket", t.Accent, t.Space(12)).Radius(999), ui.Pose{Y: t.Space(3), Scale: 0.7, Fade: 1}, 0)
		heading(c, o.Title, o.Text)
		if secs == 0 {
			c.Key("live")
			live := ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(4)).Radius(999).
				Background(t.Success.Alpha(0.14)).TextColor(t.Success)
			kit.Arrive(live, ui.Pose{Scale: 0.6, Fade: 1}, 0).Children(func() {
				glyph(c, "circle-check", 18)
				ui.Text(c, o.Live).FontSize(t.Rem(1)).Bold().Shrink(0).MaxWidth(t.Space(60)).SingleLine()
			})
			return
		}
		values := [4]string{number(c, secs/86400), fmt.Sprintf("%02d", secs/3600%24), fmt.Sprintf("%02d", secs/60%60), fmt.Sprintf("%02d", secs%60)}
		c.Key("tiles")
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			for i, v := range values {
				cell := ui.Column(c).Grow(1).Basis(0).MinWidth(0).Gap(t.Space(1)).AlignItems(ui.Center).Padding(t.Space(3), t.Space(1)).
					Radius(t.Radius).Background(t.SurfaceHover)
				kit.Arrive(cell, ui.Pose{Y: t.Space(2), Fade: 1}, time.Duration(i)*70*time.Millisecond).Children(func() {
					kit.Numeric(c, v).FontSize(t.Rem(1.6)).Bold().FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine().
						ContentTransition(ui.ContentTransition{Numeric: true})
					kit.Muted(c, o.Units[i]).FontSize(t.Rem(0.7)).FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
				})
			}
		})
	})
	return root
}
