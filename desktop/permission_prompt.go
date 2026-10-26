package desktop

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Access is a permission answer.
type Access int

// Answers.
const (
	Undecided Access = iota
	Allowed
	Denied
)

// PermissionPromptOptions word the prompt; Title formats app and resource.
type PermissionPromptOptions struct{ Icon, Title, Text, Allow, Deny, Allowed, Denied string }

// PermissionPrompt asks to allow access; cue Tap.
func PermissionPrompt(c *ui.Context, app, resource string, access *Access, o PermissionPromptOptions) ui.Element {
	def(&o.Icon, "shield")
	def(&o.Title, "%s would like to access your %s")
	def(&o.Allow, "Allow")
	def(&o.Deny, "Don't allow")
	def(&o.Allowed, "Access allowed. Change it in Settings.")
	def(&o.Denied, "Access denied. Change it in Settings.")
	if strings.TrimSpace(app) == "" || strings.TrimSpace(resource) == "" || access == nil || *access < Undecided || *access > Denied {
		panic(fmt.Sprintf("desktop: PermissionPrompt %q needs a resource and a known answer", app))
	}
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(5)).AlignItems(ui.Center).Role(ui.RoleAlertDialog).Label(fmt.Sprintf(o.Title, app, resource))
	return root.Children(func() {
		kit.Arrive(tile(c, o.Icon, t.Accent, t.Space(12)).Radius(999), ui.Pose{Scale: 0.6, Fade: 1}, 0)
		ui.Text(c, fmt.Sprintf(o.Title, app, resource)).FontSize(t.Rem(0.92)).Bold().FillWidth().TextAlign(ui.Center).Wrap()
		if o.Text != "" {
			kit.Muted(c, o.Text).FontSize(t.Rem(0.8)).FillWidth().TextAlign(ui.Center).Wrap()
		}
		c.Key("answers")
		ui.Row(c).FillWidth().Wrap().Gap(t.Space(2)).Children(func() {
			for _, a := range [2]Access{Denied, Allowed} {
				label := o.Deny
				if a == Allowed {
					label = o.Allow
				}
				on := *access == a
				b := kit.Press(ui.ButtonBase(c).Label(label).Checked(on).Grow(1).Basis(t.Space(30)).Padding(t.Space(2), t.Space(3)).Radius(999))
				bg := t.Background.Mix(t.SurfaceHover, unit(b.Hovered()))
				if a == Allowed {
					bg = t.Accent.Mix(t.Text, 0.12*unit(b.Hovered()))
				}
				b.Background(b.AnimateColor("bg", bg, ui.Snappy)).Border(2, b.AnimateColor("ring", t.Border.Mix(t.Accent, unit(on && a == Denied)).Mix(t.Text.Alpha(0.3), unit(on && a == Allowed)), ui.Snappy))
				fg := t.Text
				if a == Allowed {
					fg = t.AccentText
				}
				b.Children(func() {
					ui.Text(c, label).FontSize(t.Rem(0.8)).Bold().TextColor(fg).FillWidth().MinWidth(0).TextAlign(ui.Center).SingleLine()
				})
				if b.Clicked() && !on {
					*access = a
					root.SetChanged()
					sound.Play(sound.Tap)
					kit.Log("PermissionPrompt", "answered", "allowed", a == Allowed)
				}
			}
		})
		if *access != Undecided {
			note, mark, ink := o.Allowed, "circle-check", t.Success
			if *access == Denied {
				note, mark, ink = o.Denied, "circle-x", t.Danger
			}
			c.Key(int(*access))
			kit.Arrive(ui.Row(c).FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center), ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
				glyph(c, mark, 14).TextColor(ink)
				kit.Muted(c, note).FontSize(t.Rem(0.76)).Grow(1).MinWidth(0).Wrap()
			})
		}
	})
}
