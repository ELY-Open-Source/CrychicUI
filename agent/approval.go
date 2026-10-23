package agent

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ApprovalOptions name the title and buttons.
type ApprovalOptions struct {
	Title, Allow, Deny string
}

// Approval asks "allow this action?" with Allow and Deny.
func Approval(c *ui.Context, action string, o ApprovalOptions) (root, allow, deny ui.Element) {
	if action == "" {
		panic("agent: Approval needs an action")
	}
	def(&o.Title, "Allow this action?")
	def(&o.Allow, "Allow")
	def(&o.Deny, "Deny")
	t := c.Theme()
	root = kit.Card(c).AlignSelf(ui.Start).Width(t.Space(86)).Gap(t.Space(3)).Padding(t.Space(4)).Role(ui.RoleGroup).Label(o.Title)
	root.Children(func() {
		ui.Row(c).Gap(t.Space(2)).AlignItems(ui.Center).Children(func() {
			ui.Icon(c, icon.Must("triangle-alert")).Size(16, 16).TextColor(t.Warning)
			ui.Text(c, o.Title).Bold()
		})
		ui.Text(c, action).Font("monospace").FontSize(t.Rem(0.85)).Wrap().Padding(t.Space(2)).Radius(t.Radius).Background(t.Surface)
		ui.Row(c).FillWidth().Gap(t.Space(2)).Justify(ui.End).Children(func() {
			deny = kit.Press(ui.ButtonBase(c).Label(o.Deny).Padding(t.Space(1.5), t.Space(4)).Radius(t.Radius).Border(1, t.Border))
			deny.Background(deny.AnimateColor("bg", t.Surface.Mix(t.SurfaceHover, unit(deny.Hovered())), ui.Snappy)).
				Children(func() { ui.Text(c, o.Deny).SingleLine() })
			allow = kit.Press(ui.ButtonBase(c).Label(o.Allow).Padding(t.Space(1.5), t.Space(4)).Radius(t.Radius))
			allow.Background(allow.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(allow.Hovered())), ui.Snappy)).
				TextColor(t.AccentText).Children(func() { ui.Text(c, o.Allow).Bold().SingleLine() })
		})
	})
	if allow.Clicked() {
		root.SetChanged()
		sound.Play(sound.Success)
		kit.Log("Approval", "allowed")
	}
	if deny.Clicked() {
		root.SetChanged()
		sound.Play(sound.Close)
		kit.Log("Approval", "denied")
	}
	return root, allow, deny
}
