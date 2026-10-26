package home

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// PlugOptions word the plug; zero Today hides the day's use.
type PlugOptions struct {
	Room, Watts, Used, Switch string
	Today                     float64
}

// Plug switches an outlet, showing its draw; cue Toggle.
func Plug(c *ui.Context, name string, on *bool, watts float64, o PlugOptions) ui.Element {
	if strings.TrimSpace(name) == "" || on == nil || !(watts >= 0 && watts <= 1e5) || !(o.Today >= 0 && o.Today <= 1e6) {
		panic(fmt.Sprintf("home: Plug %q needs an on state, 0 to 1e5 watts and 0 to 1e6 kWh today", name))
	}
	def(&o.Watts, "%s W")
	def(&o.Used, "%s kWh today")
	def(&o.Switch, "Switch %s")
	t := c.Theme()
	draw := 0.0
	if *on {
		draw = watts
	}
	root := card(c).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			mark := "plug"
			if *on {
				mark = "plug-zap"
			}
			tile(c, mark, t.Warning, t.Space(11), *on)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, name).FontSize(t.Rem(0.92)).Bold().FillWidth().MinWidth(0).SingleLine()
				ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
					kit.Numeric(c, fmt.Sprintf(o.Watts, strconv.FormatFloat(draw, 'f', 1, 64))).FontSize(t.Rem(0.8)).Bold().
						TextColor(t.TextMuted.Mix(t.Warning, unit(*on))).MinWidth(0).SingleLine().
						ContentTransition(ui.ContentTransition{Numeric: true})
					if o.Today > 0 {
						kit.Muted(c, fmt.Sprintf(o.Used, strconv.FormatFloat(o.Today, 'f', 2, 64))).FontSize(t.Rem(0.72)).Grow(1).MinWidth(0).SingleLine()
					}
				})
				if o.Room != "" {
					kit.Muted(c, o.Room).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				}
			})
			c.Key("switch")
			sw := toggle(c, fmt.Sprintf(o.Switch, name), *on)
			if sw.Clicked() {
				*on = !*on
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("Plug", "switched", "on", *on)
			}
		})
	})
	return root
}
