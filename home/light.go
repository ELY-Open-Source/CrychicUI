package home

import (
	"fmt"
	"math"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// LightOptions word the light; Tones are warm, daylight and cool.
type LightOptions struct {
	Room, Off, Level, Brightness, Switch string
	Tones                                [3]string
}

// Light switches a lamp, dims and tints it; cue Toggle.
func Light(c *ui.Context, name string, on *bool, level *float64, tone *int, o LightOptions) ui.Element {
	if strings.TrimSpace(name) == "" || on == nil || level == nil || tone == nil {
		panic(fmt.Sprintf("home: Light %q needs a name, an on state, a level and a tone", name))
	}
	if !(*level >= 0.01 && *level <= 1) || *tone < 0 || *tone > 2 {
		panic(fmt.Sprintf("home: Light %q needs a level of 0.01 to 1 and a tone of 0 to 2", name))
	}
	def(&o.Off, "Off")
	def(&o.Level, "%d%%")
	def(&o.Brightness, "Brightness")
	def(&o.Switch, "Switch %s")
	for i, s := range [...]string{"Warm", "Daylight", "Cool"} {
		def(&o.Tones[i], s)
	}
	t := c.Theme()
	inks := [3]ui.Color{t.Warning, t.Warning.Mix(t.Background, 0.55), t.Accent}
	ink := inks[*tone]
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			mark := "lightbulb-off"
			if *on {
				mark = "lightbulb"
			}
			tile(c, mark, ink, t.Space(11), *on)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, name).FontSize(t.Rem(0.92)).Bold().FillWidth().MinWidth(0).SingleLine()
				if o.Room != "" {
					kit.Muted(c, o.Room).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				}
			})
			state := o.Off
			if *on {
				state = fmt.Sprintf(o.Level, int(math.Round(*level*100)))
			}
			c.Key("state")
			kit.Numeric(c, state).FontSize(t.Rem(0.85)).Bold().Shrink(0).TextColor(t.TextMuted.Mix(t.Text, unit(*on))).
				ContentTransition(ui.ContentTransition{Numeric: true})
			c.Key("switch")
			sw := toggle(c, fmt.Sprintf(o.Switch, name), *on)
			if sw.Clicked() {
				*on = !*on
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("Light", "switched", "on", *on)
			}
		})
		slider := ui.Slider(c, level, 0.01, 1).Step(0.01).Label(o.Brightness).FillWidth().Disabled(!*on)
		if slider.Changed() {
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("Light", "dimmed", "percent", int(math.Round(*level*100)))
		}
		ui.Row(c).FillWidth().Gap(t.Space(2)).Children(func() {
			for i, name := range o.Tones {
				picked := i == *tone
				b := kit.Press(ui.ButtonBase(c).Label(name).Checked(picked).Grow(1).Basis(0).MinWidth(0).Padding(t.Space(1.5), t.Space(2)).
					Radius(999).Disabled(!*on))
				b.Border(1, b.AnimateColor("edge", t.Border.Mix(t.Accent, unit(picked)), ui.Snappy))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())).Mix(t.Accent, 0.08*unit(picked)), ui.Snappy))
				b.Children(func() {
					ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						ui.Box(c).Size(t.Space(2.5), t.Space(2.5)).Shrink(0).Radius(999).Background(inks[i])
						ui.Text(c, name).FontSize(t.Rem(0.74)).MinWidth(0).SingleLine()
					})
				})
				if b.Clicked() && !picked {
					*tone = i
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Light", "tone", "index", i)
				}
			}
		})
	})
	return root
}
