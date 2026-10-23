package agent

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ToolResultOptions tune a ToolResult panel.
type ToolResultOptions struct {
	// Title words the header.
	Title string
}

// ToolResult is a collapsible result body.
func ToolResult(c *ui.Context, body string, open *bool, o ToolResultOptions) ui.Element {
	if body == "" || open == nil {
		panic("agent: ToolResult needs a body and open")
	}
	def(&o.Title, "Result")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
		Border(1, t.Border).Role(ui.RoleGroup)
	var head ui.Element
	root.Children(func() {
		head = kit.Press(ui.ButtonBase(c).Label(o.Title).Row().FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center))
		head.Children(func() {
			ch := ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(t.TextMuted)
			ch.Rotate(90 * ch.AnimateTo("turn", unit(*open), ui.Snappy))
			ui.Icon(c, icon.Must("wrench")).Size(16, 16).TextColor(t.TextMuted)
			ui.Text(c, o.Title).Bold().Grow(1).SingleLine()
		})
		layout.Collapse(c, *open, func() {
			ui.Text(c, body).Font("monospace").FontSize(t.Rem(0.85)).Wrap().Padding(t.Space(2)).Radius(t.Radius).Background(t.Surface)
		})
	})
	if head.Clicked() {
		*open = !*open
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("ToolResult", "open", *open)
	}
	return root
}
