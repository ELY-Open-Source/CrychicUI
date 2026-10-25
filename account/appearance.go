package account

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ThemeMode is the chosen color theme.
type ThemeMode int

const (
	ThemeLight ThemeMode = iota
	ThemeDark
	ThemeSystem
)

// Density is the chosen spacing.
type Density int

const (
	DensityCompact Density = iota
	DensityComfortable
	DensitySpacious
)

// AppearanceOptions word the card.
type AppearanceOptions struct{ Title, Theme, Light, Dark, System, Density, Compact, Comfortable, Spacious string }

// Appearance picks theme and density; cue Select.
func Appearance(c *ui.Context, mode *ThemeMode, density *Density, o AppearanceOptions) ui.Element {
	if mode == nil || density == nil || *mode < ThemeLight || *mode > ThemeSystem || *density < DensityCompact || *density > DensitySpacious {
		panic(fmt.Sprintf("account: Appearance %q needs a theme and density in range", o.Title))
	}
	def(&o.Title, "Appearance")
	def(&o.Theme, "Theme")
	def(&o.Light, "Light")
	def(&o.Dark, "Dark")
	def(&o.System, "System")
	def(&o.Density, "Density")
	def(&o.Compact, "Compact")
	def(&o.Comfortable, "Comfortable")
	def(&o.Spacious, "Spacious")
	t := c.Theme()
	root := card(c).Gap(t.Space(4)).Padding(t.Space(5)).Role(ui.RoleGroup)
	root.Children(func() {
		head(c, "palette", o.Title, "")
		ui.Text(c, o.Theme).FontSize(t.Rem(0.78)).Bold().TextColor(t.TextMuted)
		modes := ui.Row(c).FillWidth().Gap(t.Space(2.5)).Role(ui.RoleRadioGroup)
		modes.Children(func() {
			for i, name := range []string{o.Light, o.Dark, o.System} {
				m := ThemeMode(i)
				on := *mode == m
				b := kit.Press(ui.ButtonBase(c).Label(name).Grow(1).Basis(0).MinWidth(0).Column().Gap(t.Space(2)).
					Padding(t.Space(2)).Radius(t.Radius).Checked(on))
				b.Background(b.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(b.Hovered())), ui.Snappy))
				b.Border(1, b.AnimateColor("line", t.Border.Mix(t.Accent, unit(on)), ui.Snappy))
				kit.Stagger(c, modes, b, i).Children(func() {
					preview(c, m)
					ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
						ui.Text(c, name).FontSize(t.Rem(0.8)).Bold().Grow(1).MinWidth(0).SingleLine()
						if on {
							c.Key("check")
							kit.Arrive(glyph(c, "circle-check", 14).TextColor(t.Accent), ui.Pose{Scale: 0.4, Fade: 1}, 0)
						}
					})
				})
				if b.Clicked() && !on {
					*mode = m
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Appearance", "theme", "mode", i)
				}
			}
		})
		ui.Text(c, o.Density).FontSize(t.Rem(0.78)).Bold().TextColor(t.TextMuted)
		names := []string{o.Compact, o.Comfortable, o.Spacious}
		track := ui.Row(c).FillWidth().Padding(t.Space(0.75)).Radius(t.Radius).Background(t.SurfaceHover).Role(ui.RoleRadioGroup)
		at := track.AnimateTo("thumb", float32(*density), ui.Snappy)
		pad := t.Space(0.75)
		track.Draw(func(p *ui.Painter, r ui.Rect) {
			w := (r.W - 2*pad) / 3
			p.Fill(ui.Rect{X: r.X + pad + at*w, Y: r.Y + pad, W: w, H: r.H - 2*pad}, t.Background, t.Radius*0.75)
		})
		track.Children(func() {
			for i, name := range names {
				d := Density(i)
				on := *density == d
				b := kit.Press(ui.ButtonBase(c).Label(name).Grow(1).Basis(0).MinWidth(0).Padding(t.Space(1.5), t.Space(1)).
					Radius(t.Radius * 0.75).Checked(on))
				b.TextColor(b.AnimateColor("ink", t.TextMuted.Mix(t.Text, max(unit(on), unit(b.Hovered()))), ui.Snappy))
				b.Children(func() { ui.Text(c, name).FontSize(t.Rem(0.8)).Bold().MinWidth(0).SingleLine() })
				if b.Clicked() && !on {
					*density = d
					root.SetChanged()
					sound.Play(sound.Select)
					kit.Log("Appearance", "density", "value", i)
				}
			}
		})
	})
	return root
}

// preview draws a tiny window in the mode's colors.
func preview(c *ui.Context, m ThemeMode) {
	t := c.Theme()
	box := ui.Box(c).FillWidth().Height(t.Space(14)).Radius(t.Radius*0.75).Clip().Border(1, t.Border)
	box.Draw(func(p *ui.Painter, r ui.Rect) {
		paint := func(x, w float32, bg, fg ui.Color) {
			p.Fill(ui.Rect{X: x, Y: r.Y, W: w, H: r.H}, bg, 0)
			p.Fill(ui.Rect{X: x + t.Space(2), Y: r.Y + t.Space(2), W: w*0.5 - t.Space(2), H: t.Space(1.5)}, t.Accent, 999)
			for k := range 3 {
				y := r.Y + t.Space(5) + float32(k)*t.Space(2.5)
				p.Fill(ui.Rect{X: x + t.Space(2), Y: y, W: w*(0.75-0.15*float32(k)) - t.Space(2), H: t.Space(1)}, fg.Alpha(0.35), 999)
			}
		}
		light, dark := t.Background, t.Text
		if t.Dark {
			light, dark = dark, light
		}
		switch m {
		case ThemeLight:
			paint(r.X, r.W, light, dark)
		case ThemeDark:
			paint(r.X, r.W, dark, light)
		default:
			paint(r.X, r.W/2, light, dark)
			paint(r.X+r.W/2, r.W/2, dark, light)
		}
	})
}
