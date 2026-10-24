package code

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// EnvVar is one environment variable.
type EnvVar struct {
	Key, Value string
}

// EnvOptions are the copy knobs of Env.
type EnvOptions struct {
	Show, Hide, Empty string
}

// mask hides a value without telling its length.
const mask = "••••••••"

// Env shows KEY=value rows, values masked until shown; cue Toggle.
func Env(c *ui.Context, label string, vars []EnvVar, shown *bool, o EnvOptions) ui.Element {
	if label == "" || shown == nil {
		panic(fmt.Sprintf("code: Env %q needs a label and a bound toggle", label))
	}
	seen := map[string]bool{}
	for _, v := range vars {
		if v.Key == "" || seen[v.Key] {
			panic(fmt.Sprintf("code: Env %q needs a unique key in every variable", label))
		}
		seen[v.Key] = true
	}
	def(&o.Show, "Show values")
	def(&o.Hide, "Hide values")
	def(&o.Empty, "No variables")
	t := c.Theme()
	root := kit.Card(c).FillWidth().Role(ui.RoleGroup).Label(label)
	var toggle ui.Element
	root.Children(func() {
		panelHead(c, "key-round", label, func() {
			name, glyph := o.Show, "eye"
			if *shown {
				name, glyph = o.Hide, "eye-off"
			}
			toggle = round(c, name, glyph).Size(t.Space(7), t.Space(7))
			toggle.Background(toggle.AnimateColor("on", t.SurfaceHover.Alpha(unit(toggle.Hovered())).Mix(t.Accent.Alpha(0.14), unit(*shown)), ui.Snappy))
		})
		ui.Column(c).FillWidth().Padding(t.Space(1), 0).Children(func() {
			if len(vars) == 0 {
				kit.Muted(c, o.Empty).FontSize(t.Rem(0.85)).Padding(t.Space(3)).AlignSelf(ui.Center)
			}
			for i, v := range vars {
				c.Key(v.Key)
				row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(1.25), t.Space(3)).
					BorderWidth(0, 0, 1, 0).BorderColor(t.Border.Alpha(unit(i < len(vars)-1)))
				stagger(c, root, row, i).Children(func() {
					ui.Text(c, v.Key).Font("monospace").FontSize(t.Rem(0.8)).Bold().Shrink(0)
					kit.Muted(c, "=").Font("monospace").FontSize(t.Rem(0.8)).Shrink(0)
					value, ink := v.Value, t.Text
					if !*shown {
						value, ink = mask, t.TextMuted
					}
					ui.Text(c, value).Font("monospace").FontSize(t.Rem(0.8)).TextColor(ink).Grow(1).MinWidth(0).SingleLine().
						ContentTransition(ui.ContentTransition{})
					lock := ui.Icon(c, icon.Must("lock")).Size(12, 12).Shrink(0).TextColor(t.TextMuted)
					lock.Opacity(lock.AnimateTo("hidden", unit(!*shown), ui.Snappy))
				})
			}
		})
	})
	if toggle.Clicked() {
		*shown = !*shown
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("Env", "shown", "on", *shown)
	}
	return root
}
