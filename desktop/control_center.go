package desktop

import (
	"math"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Controls are the quick settings ControlCenter edits.
type Controls struct {
	Wifi, Bluetooth, Focus, Dark bool
	Volume, Brightness           float64
}

// ControlCenterOptions name the four toggles, the sliders and states.
type ControlCenterOptions struct {
	Toggles                     [4]string
	Volume, Brightness, On, Off string
}

// ControlCenter has quick toggles and sliders; cue Toggle.
func ControlCenter(c *ui.Context, s *Controls, o ControlCenterOptions) ui.Element {
	for i, d := range [4]string{"Wi-Fi", "Bluetooth", "Focus", "Dark mode"} {
		def(&o.Toggles[i], d)
	}
	def(&o.Volume, "Sound")
	def(&o.Brightness, "Display")
	def(&o.On, "On")
	def(&o.Off, "Off")
	if s == nil {
		panic("desktop: ControlCenter needs controls")
	}
	share("ControlCenter", o.Volume, float32(s.Volume))
	share("ControlCenter", o.Brightness, float32(s.Brightness))
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(2.5)).Padding(t.Space(3)).Radius(t.Radius*2).Background(t.SurfaceHover).
		Border(1, t.Border).Shadow(0, 8, 24, 0, kit.Shade(0.12)).Role(ui.RoleGroup)
	kit.Arrive(root, ui.Pose{Scale: 0.96, Fade: 1}, 0)
	root.Children(func() {
		grid := ui.Grid(c).Columns(2).FillWidth().Gap(t.Space(2))
		grid.Children(func() {
			for i, p := range [4]*bool{&s.Wifi, &s.Bluetooth, &s.Focus, &s.Dark} {
				mark := [4]string{"wifi", "bluetooth", "moon", "sun-medium"}[i]
				on := *p
				b := kit.Press(ui.ButtonBase(c).Label(o.Toggles[i]).Role(ui.RoleSwitch).Checked(on).Padding(t.Space(2)).Radius(t.Radius * 1.5).MinWidth(0))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.Border, 0.4*unit(b.Hovered())), ui.Snappy))
				kit.Stagger(c, grid, b, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
						dot := ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999)
						dot.Background(dot.AnimateColor("dot", t.SurfaceHover.Mix(t.Accent, unit(on)), ui.Snappy))
						dot.Children(func() {
							glyph(c, mark, 16).TextColor(dot.AnimateColor("ink", t.TextMuted.Mix(t.AccentText, unit(on)), ui.Snappy))
						})
						ui.Column(c).Grow(1).MinWidth(0).Children(func() {
							ui.Text(c, o.Toggles[i]).FontSize(t.Rem(0.78)).Bold().FillWidth().MinWidth(0).SingleLine()
							state := o.Off
							if on {
								state = o.On
							}
							kit.Muted(c, state).FontSize(t.Rem(0.7)).FillWidth().MinWidth(0).SingleLine().ContentTransition(ui.ContentTransition{})
						})
					})
				})
				if b.Clicked() {
					*p = !on
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("ControlCenter", "toggled", "index", i, "on", !on)
				}
			}
		})
		for i, v := range [2]*float64{&s.Volume, &s.Brightness} {
			label, low, high := o.Volume, "volume-x", "volume-2"
			if i == 1 {
				label, low, high = o.Brightness, "sun", "sun-medium"
			}
			c.Key(i)
			ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(2.5), t.Space(3)).Radius(t.Radius * 1.5).Background(t.Background).Children(func() {
				kit.Muted(c, label).FontSize(t.Rem(0.72)).Bold().FillWidth().MinWidth(0).SingleLine()
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					glyph(c, low, 14).TextColor(t.TextMuted)
					sl := ui.Slider(c, v, 0, 1).Step(0.01).Label(label).Grow(1).MinWidth(0)
					if sl.Changed() {
						root.SetChanged()
						sound.Play(sound.Tick)
						kit.Log("ControlCenter", "slid", "index", i, "percent", int(math.Round(*v*100)))
					}
					glyph(c, high, 14).TextColor(t.TextMuted)
				})
			})
		}
	})
	return root
}
