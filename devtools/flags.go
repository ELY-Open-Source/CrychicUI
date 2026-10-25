package devtools

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Flag is one feature flag; Rollout is a percent.
type Flag struct {
	Key, Description string
	On               bool
	Rollout          int
}

// FlagsOptions word the panel; Count and Rollout are formats.
type FlagsOptions struct{ Title, Count, Rollout string }

// Flags toggles feature flags; cue Toggle.
func Flags(c *ui.Context, flags *[]Flag, o FlagsOptions) ui.Element {
	if flags == nil || len(*flags) == 0 {
		panic("devtools: Flags needs flags")
	}
	keys := make([]string, len(*flags))
	for i, f := range *flags {
		if f.Rollout < 0 || f.Rollout > 100 {
			panic(fmt.Sprintf("devtools: Flags %q rollout %d is outside 0..100", f.Key, f.Rollout))
		}
		keys[i] = f.Key
	}
	unique("Flags", keys)
	def(&o.Title, "Feature flags")
	def(&o.Count, "%d of %d on")
	def(&o.Rollout, "%d%%")
	t := c.Theme()
	on := 0
	for _, f := range *flags {
		if f.On {
			on++
		}
	}
	root := card(c).Role(ui.RoleGroup)
	root.Children(func() {
		kit.CardHead(c, "flag", o.Title, func() { kit.Chip(c, fmt.Sprintf(o.Count, on, len(*flags)), t.Accent) })
		list := ui.Column(c).FillWidth().Padding(t.Space(1.5))
		list.Children(func() {
			for i, f := range *flags {
				c.Key(f.Key)
				row := kit.Press(ui.ButtonBase(c).Label(f.Key).FillWidth().Padding(t.Space(2.5), t.Space(2)).Radius(t.Radius).
					Role(ui.RoleSwitch).Checked(f.On))
				row.Background(row.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(row.Hovered())), ui.Snappy))
				kit.Stagger(c, list, row, i).Children(func() {
					ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
						ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.5)).Children(func() {
							ui.Text(c, f.Key).Font("monospace").FontSize(t.Rem(0.8)).Bold().FillWidth().MinWidth(0).SingleLine()
							if f.Description != "" {
								kit.Muted(c, f.Description).FontSize(t.Rem(0.75)).Wrap().FillWidth()
							}
						})
						ui.Column(c).Width(t.Space(16)).Shrink(0).Gap(t.Space(0.75)).AlignItems(ui.End).Children(func() {
							kit.Numeric(c, fmt.Sprintf(o.Rollout, f.Rollout)).FontSize(t.Rem(0.7)).TextColor(t.TextMuted).Shrink(0)
							ink := t.TextMuted
							if f.On {
								ink = t.Accent
							}
							meter(c, float32(f.Rollout)/100, ink).Height(t.Space(1))
						})
						track(c, f.On, row.Hovered())
					})
				})
				if row.Clicked() {
					(*flags)[i].On = !f.On
					root.SetChanged()
					sound.Play(sound.Toggle)
					kit.Log("Flags", "toggled", "index", i, "on", !f.On)
				}
			}
		})
	})
	return root
}
