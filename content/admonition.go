package content

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// AdmonitionKind is the tone of a callout.
type AdmonitionKind int

// Callout tones.
const (
	Note AdmonitionKind = iota
	Tip
	Warning
	Danger
)

// AdmonitionOptions name the kinds; they title untitled blocks.
type AdmonitionOptions struct{ Kinds [4]string }

// Admonition is a note, tip, warning or danger block.
func Admonition(c *ui.Context, kind AdmonitionKind, title, text string, o AdmonitionOptions) ui.Element {
	for i, d := range [4]string{"Note", "Tip", "Warning", "Danger"} {
		def(&o.Kinds[i], d)
	}
	if kind < Note || kind > Danger || strings.TrimSpace(text) == "" {
		panic(fmt.Sprintf("content: Admonition %q needs a known kind and text", title))
	}
	def(&title, o.Kinds[kind])
	t := c.Theme()
	ink := [4]ui.Color{t.Accent, t.Success, t.Warning, t.Danger}[kind]
	mark := [4]string{"info", "lightbulb", "triangle-alert", "octagon-alert"}[kind]
	root := ui.Row(c).FillWidth().Gap(t.Space(3)).Padding(t.Space(3.5), t.Space(4)).Radius(t.Radius).
		Background(ink.Alpha(0.08)).BorderWidth(0, 0, 0, t.Space(1)).BorderColor(ink).AlignItems(ui.Start).Role(ui.RoleGroup).Label(title)
	kit.Arrive(root, ui.Pose{X: -t.Space(2), Fade: 1}, 0)
	return root.Children(func() {
		kit.Arrive(glyph(c, mark, 18).TextColor(ink), ui.Pose{Scale: 0.6, Fade: 1}, 120)
		ui.Column(c).Grow(1).MinWidth(0).Gap(t.Space(1)).Children(func() {
			ui.Text(c, title).FontSize(t.Rem(0.86)).Bold().TextColor(ink).FillWidth().Wrap()
			ui.Text(c, text).FontSize(t.Rem(0.84)).LineHeight(1.55).FillWidth().Wrap()
		})
	})
}
