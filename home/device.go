package home

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// DeviceOptions word the tile; Switch takes the name.
type DeviceOptions struct{ Room, On, Off, Switch string }

// Device is a tile that switches one device; cue Toggle.
func Device(c *ui.Context, name, mark string, on *bool, o DeviceOptions) ui.Element {
	if strings.TrimSpace(name) == "" || mark == "" || on == nil {
		panic(fmt.Sprintf("home: Device %q needs a name, an icon and an on state", name))
	}
	def(&o.On, "On")
	def(&o.Off, "Off")
	def(&o.Switch, "Switch %s")
	t := c.Theme()
	b := kit.Press(ui.ButtonBase(c).Label(fmt.Sprintf(o.Switch, name)).Checked(*on).Column().MinWidth(0).Gap(t.Space(3)).
		Padding(t.Space(3.5)).Radius(t.Radius * 1.5).AlignItems(ui.Stretch))
	kit.Arrive(b, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	b.Border(1, b.AnimateColor("edge", t.Border.Mix(t.Accent, 0.5*unit(*on)), ui.Snappy))
	b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, 0.6*unit(b.Hovered())).Mix(t.Accent, 0.07*unit(*on)), ui.Snappy))
	b.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Start).Children(func() {
			tile(c, mark, t.Accent, t.Space(10), *on)
			ui.Spacer(c)
			track := ui.Box(c).Size(t.Space(8), t.Space(4.5)).Shrink(0).Radius(999)
			track.Background(track.AnimateColor("track", t.SurfaceHover.Mix(t.Accent, unit(*on)), ui.Snappy))
			x := track.AnimateTo("knob", unit(*on), ui.Snappy)
			pad := t.Space(0.5)
			track.Draw(func(p *ui.Painter, r ui.Rect) {
				d := r.H - 2*pad
				p.Fill(ui.Rect{X: r.X + pad + x*(r.W-r.H), Y: r.Y + pad, W: d, H: d}, t.Background, d/2)
			})
		})
		ui.Column(c).FillWidth().Gap(t.Space(0.5)).Children(func() {
			ui.Text(c, name).FontSize(t.Rem(0.88)).Bold().FillWidth().MinWidth(0).SingleLine()
			ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
				state := o.Off
				if *on {
					state = o.On
				}
				c.Key(state)
				kit.Arrive(ui.Text(c, state).FontSize(t.Rem(0.72)).Bold().TextColor(t.TextMuted.Mix(t.Accent, unit(*on))).Shrink(0),
					ui.Pose{Y: t.Space(1), Fade: 1}, 0)
				if o.Room != "" {
					kit.Muted(c, o.Room).FontSize(t.Rem(0.72)).Grow(1).MinWidth(0).SingleLine()
				}
			})
		})
	})
	if b.Clicked() {
		*on = !*on
		b.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("Device", "switched", "on", *on)
	}
	return b
}
