package home

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Rule is a when-then automation; Last is shown as given.
type Rule struct {
	Name, When, Then, WhenIcon, ThenIcon, Last string
}

// AutomationOptions word the rule; Switch takes the name.
type AutomationOptions struct{ When, Then, Switch string }

// Automation shows a when-then rule with a switch; cue Toggle.
func Automation(c *ui.Context, r Rule, on *bool, o AutomationOptions) ui.Element {
	if strings.TrimSpace(r.Name) == "" || r.When == "" || r.Then == "" || r.WhenIcon == "" || r.ThenIcon == "" || on == nil {
		panic(fmt.Sprintf("home: Automation %q needs a when, a then, their icons and an on state", r.Name))
	}
	def(&o.When, "When")
	def(&o.Then, "Then")
	def(&o.Switch, "Switch %s")
	t := c.Theme()
	root := card(c).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup)
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "workflow", t.Accent, t.Space(10), *on)
			ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
				ui.Text(c, r.Name).FontSize(t.Rem(0.9)).Bold().FillWidth().MinWidth(0).SingleLine()
				if r.Last != "" {
					kit.Muted(c, r.Last).FontSize(t.Rem(0.72)).FillWidth().MinWidth(0).SingleLine()
				}
			})
			c.Key("switch")
			sw := toggle(c, fmt.Sprintf(o.Switch, r.Name), *on)
			if sw.Clicked() {
				*on = !*on
				root.SetChanged()
				sound.Play(sound.Toggle)
				kit.Log("Automation", "switched", "on", *on)
			}
		})
		body := ui.Column(c).FillWidth().Gap(t.Space(1)).AlignItems(ui.Center)
		body.Opacity(body.AnimateTo("live", 0.45+0.55*unit(*on), ui.Smooth))
		body.Children(func() {
			step := func(head, text, mark string, ink ui.Color) {
				box := ui.Row(c).FillWidth().Gap(t.Space(2.5)).AlignItems(ui.Center).Padding(t.Space(2.5)).Radius(t.Radius).
					Border(1, t.Border).Background(ink.Alpha(0.05))
				box.Children(func() {
					ui.Box(c).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999).Background(ink.Alpha(0.14)).TextColor(ink).
						Children(func() { glyph(c, mark, 16) })
					ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(0.25)).Children(func() {
						kit.Muted(c, head).FontSize(t.Rem(0.64)).Bold().FillWidth().MinWidth(0).SingleLine()
						ui.Text(c, text).FontSize(t.Rem(0.82)).FillWidth().MinWidth(0).Wrap()
					})
				})
			}
			step(o.When, r.When, r.WhenIcon, t.Warning)
			glyph(c, "arrow-down", 14).TextColor(t.TextMuted)
			step(o.Then, r.Then, r.ThenIcon, t.Accent)
		})
	})
	return root
}
