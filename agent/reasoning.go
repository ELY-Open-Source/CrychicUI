package agent

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReasoningOptions tune a Reasoning panel.
type ReasoningOptions struct {
	// Title words the header; the body is caller copy.
	Title string
}

// Reasoning is a collapsible reasoning panel; height springs.
func Reasoning(c *ui.Context, body string, open *bool, o ReasoningOptions) ui.Element {
	if body == "" || open == nil {
		panic("agent: Reasoning needs a body and open")
	}
	def(&o.Title, "Reasoning")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Padding(t.Space(2), t.Space(3)).Radius(t.Radius).
		Background(t.Surface).Role(ui.RoleGroup)
	var head ui.Element
	root.Children(func() {
		head = kit.Press(ui.ButtonBase(c).Label(o.Title).Row().FillWidth().Gap(t.Space(1.5)).AlignItems(ui.Center))
		head.Children(func() {
			ch := ui.Icon(c, icon.Must("chevron-right")).Size(16, 16).TextColor(t.TextMuted)
			ch.Rotate(90 * ch.AnimateTo("turn", unit(*open), ui.Snappy))
			kit.Muted(c, o.Title).Italic().Grow(1).SingleLine()
		})
		layout.Collapse(c, *open, func() {
			kit.Muted(c, body).FontSize(t.Rem(0.9)).Italic().Wrap()
		})
	})
	if head.Clicked() {
		*open = !*open
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("Reasoning", "open", *open)
	}
	return root
}
