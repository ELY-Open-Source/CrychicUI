package mail

import (
	"time"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ScreenShareOptions word the bar; Detail is optional.
type ScreenShareOptions struct{ Sharing, Detail, Stop string }

// ScreenShare stops the active share; cue Close.
func ScreenShare(c *ui.Context, sharing *bool, o ScreenShareOptions) ui.Element {
	if sharing == nil {
		panic("mail: ScreenShare needs state")
	}
	def(&o.Sharing, "You are sharing")
	def(&o.Stop, "Stop")
	t := c.Theme()
	root := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(2), t.Space(2), t.Space(3)).
		Radius(t.Radius).Background(t.Accent).TextColor(t.AccentText).Shadow(0, 1, 3, 0, kit.Shade(0.08)).Role(ui.RoleGroup)
	shown := root.AnimateTo("sharing", unit(*sharing), ui.Snappy)
	root.Opacity(shown).Translate(0, t.Space(2)*(1-shown))
	root.Children(func() {
		dot := ui.Box(c).Size(t.Space(2), t.Space(2)).Shrink(0).Radius(999)
		beat := float32(0)
		if *sharing {
			beat = dot.Loop("beat", 1200*time.Millisecond, ui.EaseInOut)
		}
		dot.Background(t.Danger.Mix(t.AccentText, 0.2)).Opacity(0.4 + 0.6*(1-beat))
		glyph(c, "monitor-up", 16)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
			ui.Text(c, o.Sharing).FontSize(t.Rem(0.85)).Bold().FillWidth().MinWidth(0).SingleLine()
			if o.Detail != "" {
				ui.Text(c, o.Detail).FontSize(t.Rem(0.72)).TextColor(t.AccentText.Alpha(0.8)).FillWidth().MinWidth(0).SingleLine()
			}
		})
		stop := kit.Press(ui.ButtonBase(c).Label(o.Stop).Padding(t.Space(1.5), t.Space(3)).Radius(999).Shrink(0).
			TextColor(t.Danger).Disabled(!*sharing))
		stop.Background(stop.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(stop.Hovered())), ui.Snappy))
		stop.Children(func() {
			ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				glyph(c, "square", 12)
				ui.Text(c, o.Stop).FontSize(t.Rem(0.8)).Bold().Shrink(0)
			})
		})
		if stop.Clicked() && *sharing {
			*sharing = false
			root.SetChanged()
			sound.Play(sound.Close)
			kit.Log("ScreenShare", "stopped")
		}
	})
	return root
}
