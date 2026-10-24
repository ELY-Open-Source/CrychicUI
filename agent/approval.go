package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// Decision is an Approval's outcome.
type Decision uint8

// Decisions.
const (
	Undecided Decision = iota
	Allowed
	Denied
)

// ApprovalOptions name the title, buttons and outcomes.
type ApprovalOptions struct {
	Title, Allow, Deny, Allowed, Denied string
}

// Approval asks "allow this action?" with Allow and Deny.
func Approval(c *ui.Context, tool, action, reason string, decision *Decision, o ApprovalOptions) (root, allow, deny ui.Element) {
	if tool == "" || action == "" || decision == nil || *decision > Denied {
		panic(fmt.Sprintf("agent: Approval %q needs a tool, an action and a decision", tool))
	}
	def(&o.Title, "Allow this action?")
	def(&o.Allow, "Allow")
	def(&o.Deny, "Deny")
	def(&o.Allowed, "Allowed")
	def(&o.Denied, "Denied")
	t := c.Theme()
	root = kit.Card(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup).Label(o.Title)
	root.Border(1, root.AnimateColor("ring", t.Warning.Mix(t.Border, unit(*decision != Undecided)), ui.Snappy))
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(3)).AlignItems(ui.Center).Children(func() {
			tile(c, "shield-alert", t.Warning, t.Space(9))
			ui.Column(c).Grow(1).Gap(t.Space(0.5)).Children(func() {
				ui.Text(c, o.Title).Bold()
				kit.Muted(c, tool).Font("monospace").FontSize(t.Rem(0.8)).SingleLine()
			})
		})
		ui.Text(c, action).Font("monospace").FontSize(t.Rem(0.85)).FillWidth().Wrap().Padding(t.Space(2), t.Space(3)).
			Radius(t.Radius).Background(t.Surface)
		if reason != "" {
			kit.Muted(c, reason).FontSize(t.Rem(0.85)).Wrap()
		}
		ui.Box(c).FillWidth().Children(func() {
			if *decision != Undecided {
				text, glyph, ink := o.Allowed, "circle-check", t.Success
				if *decision == Denied {
					text, glyph, ink = o.Denied, "circle-x", t.Danger
				}
				done := ui.Row(c).Key(*decision).Gap(t.Space(1.5)).AlignItems(ui.Center).TextColor(ink)
				kit.Arrive(done, ui.Pose{Y: t.Space(1), Fade: 1}, 0).Children(func() {
					ui.Icon(c, icon.Must(glyph)).Size(16, 16)
					ui.Text(c, text).Bold()
				})
				return
			}
			ui.Row(c).FillWidth().Gap(t.Space(2)).Justify(ui.End).Children(func() {
				deny = kit.Press(ui.ButtonBase(c).Label(o.Deny).Padding(t.Space(1.5), t.Space(4)).Radius(t.Radius).Border(1, t.Border))
				deny.Background(deny.AnimateColor("bg", t.Background.Mix(t.SurfaceHover, unit(deny.Hovered())), ui.Snappy)).
					Children(func() { ui.Text(c, o.Deny).Shrink(0) })
				allow = kit.Press(ui.ButtonBase(c).Label(o.Allow).Padding(t.Space(1.5), t.Space(4)).Radius(t.Radius))
				allow.Background(allow.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(allow.Hovered())), ui.Snappy)).
					TextColor(t.AccentText).Children(func() { ui.Text(c, o.Allow).Bold().Shrink(0) })
			})
		})
	})
	switch {
	case allow.Clicked():
		*decision = Allowed
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("Approval", "allowed")
	case deny.Clicked():
		*decision = Denied
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("Approval", "denied")
	}
	return root, allow, deny
}
