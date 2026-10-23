// Code components for the gallery.
package code

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

type EnvVar struct {
	Key, Value string
}

// EnvOptions are the copy knobs of Env.

type EnvOptions struct {
	Show string
	Hide string
}

// Env is KEY=value rows with masked values; cue Toggle.

func Env(c *ui.Context, label string, vars []EnvVar, shown *bool, o EnvOptions) ui.Element {
	if label == "" || shown == nil {
		panic("code: Env needs a label and a bound toggle")
	}
	def(&o.Show, "Show values")
	def(&o.Hide, "Hide values")
	t := c.Theme()
	root := ui.Column(c).AlignSelf(ui.Start).Width(t.Space(80)).Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	root.Children(func() {
		ui.Row(c).FillWidth().Children(func() {
			ui.Box(c).Grow(1)
			name := o.Show
			glyph := "eye"
			if *shown {
				name = o.Hide
				glyph = "eye-off"
			}
			b := kit.Press(ui.ButtonBase(c).Label(name).Tooltip(name).Size(t.Space(8), t.Space(8)).Center().Radius(999))
			b.Background(b.AnimateColor("bg", t.SurfaceHover.Alpha(unit(b.Hovered())), ui.Snappy))
			b.Children(func() {
				ui.Icon(c, icon.Must(glyph)).Size(16, 16).ContentTransition(ui.ContentTransition{})
			})
			if b.Clicked() {
				*shown = !*shown
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("Env", "shown", "on", *shown)
			}
		})
		for _, v := range vars {
			if v.Key == "" {
				panic("code: Env needs a key in every variable")
			}
			c.Key(v.Key)
			row := ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(0.5), t.Space(1))
			row.Children(func() {
				ui.Text(c, v.Key).Font("monospace").FontSize(t.Rem(0.85)).Bold().SingleLine()
				ui.Text(c, "=").Font("monospace").FontSize(t.Rem(0.85)).TextColor(t.TextMuted)
				shownValue := v.Value
				if !*shown {
					shownValue = "••••••••"
				}
				ui.Text(c, shownValue).Font("monospace").FontSize(t.Rem(0.85)).SingleLine().
					ContentTransition(ui.ContentTransition{})
			})
		}
	})
	return root
}

// PlaygroundOptions are the copy knobs of Playground.
