package home

import (
	"fmt"
	"strings"
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// CameraOptions word the tile; Format lays out the frame time.
type CameraOptions struct{ Room, Live, Offline, Format string }

// Camera shows a live frame; nil frame is offline.
func Camera(c *ui.Context, name string, frame *ui.Bitmap, at time.Time, o CameraOptions) ui.Element {
	if strings.TrimSpace(name) == "" || frame != nil && at.IsZero() {
		panic(fmt.Sprintf("home: Camera %q needs a name and a frame time", name))
	}
	def(&o.Live, "LIVE")
	def(&o.Offline, "Offline")
	def(&o.Format, "15:04:05")
	t := c.Theme()
	root := card(c).Clip().Role(ui.RoleGroup)
	root.Children(func() {
		view := ui.Box(c).FillWidth().AspectRatio(16.0 / 9).Clip().Background(t.Text.Mix(t.Background, 0.15))
		view.Children(func() {
			if frame == nil {
				ui.Column(c).Fill().Center().Gap(t.Space(1.5)).TextColor(t.Background.Alpha(0.7)).Children(func() {
					glyph(c, "video-off", 24)
					ui.Text(c, o.Offline).FontSize(t.Rem(0.78)).Bold().Shrink(0)
				})
				return
			}
			kit.Arrive(ui.Image(c, frame).Fill().Fit(ui.Cover), ui.Pose{Scale: 1.04, Fade: 1}, 0)
			ui.Row(c).Absolute().Left(t.Space(2.5)).Top(t.Space(2.5)).Gap(t.Space(1.5)).AlignItems(ui.Center).
				Padding(t.Space(0.5), t.Space(2)).Radius(t.Radius).Background(t.Danger).TextColor(t.AccentText).Children(func() {
				dot := ui.Box(c).Size(t.Space(1.5), t.Space(1.5)).Shrink(0).Radius(999).Background(t.AccentText)
				beat := dot.Loop("beat", 1200*time.Millisecond, ui.EaseInOut)
				dot.Opacity(0.35 + 0.65*(1-beat))
				ui.Text(c, o.Live).FontSize(t.Rem(0.68)).Bold().Shrink(0)
			})
			ui.Row(c).Absolute().Right(t.Space(2.5)).Bottom(t.Space(2.5)).Padding(t.Space(0.5), t.Space(2)).Radius(t.Radius).
				Background(t.Text.Alpha(0.55)).TextColor(t.Background).Children(func() {
				kit.Numeric(c, at.Format(o.Format)).Font("monospace").FontSize(t.Rem(0.68)).Bold().Shrink(0)
			})
		})
		ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(3)).Children(func() {
			mark := "cctv"
			if frame == nil {
				mark = "video-off"
			}
			tile(c, mark, t.Accent, t.Space(9), frame != nil)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
				ui.Text(c, name).FontSize(t.Rem(0.88)).Bold().FillWidth().MinWidth(0).SingleLine()
				if o.Room != "" {
					kit.Muted(c, o.Room).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				}
			})
		})
	})
	return root
}
