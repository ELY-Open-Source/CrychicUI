package agent

import (
	"fmt"
	"strings"

	"github.com/ELY-Open-Source/CrychicUI/button"
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ToolResultOptions word the line count and copy.
type ToolResultOptions struct {
	Lines, Copy string
}

// ToolResult is a collapsible result body with copy.
func ToolResult(c *ui.Context, name, body string, open *bool, o ToolResultOptions) (root, copy ui.Element) {
	if name == "" || body == "" || open == nil {
		panic(fmt.Sprintf("agent: ToolResult %q needs a name, a body and open", name))
	}
	def(&o.Lines, "%d lines")
	def(&o.Copy, "Copy")
	t := c.Theme()
	lines := strings.Count(strings.TrimRight(body, "\n"), "\n") + 1
	root = ui.Column(c).FillWidth().Radius(t.Radius).Background(t.Background).Border(1, t.Border).Clip()
	var head ui.Element
	root.Children(func() {
		ui.Row(c).FillWidth().AlignItems(ui.Center).PaddingEnd(t.Space(2)).Children(func() {
			head = ui.ButtonBase(c).Label(name).Row().Grow(1).Gap(t.Space(2)).AlignItems(ui.Center).
				Padding(t.Space(2), t.Space(3)).Checked(*open)
			head.Background(head.AnimateColor("bg", t.SurfaceHover.Alpha(unit(head.Hovered())), ui.Snappy))
			head.Children(func() {
				ch := ui.Icon(c, icon.Must("chevron-down")).Size(16, 16).TextColor(t.TextMuted)
				ch.Rotate(-90 + 90*ch.AnimateTo("turn", unit(*open), ui.Snappy))
				ui.Icon(c, icon.Must("file-output")).Size(16, 16).TextColor(t.Success)
				ui.Text(c, name).Bold().FontSize(t.Rem(0.9)).Grow(1).SingleLine()
				tag(c, fmt.Sprintf(o.Lines, lines))
			})
			copy = button.CopyButton(c, o.Copy, body)
		})
		layout.Collapse(c, *open, func() {
			ui.Scroll(c).FillWidth().MaxHeight(t.Space(48)).Padding(t.Space(2), t.Space(3)).Background(t.Surface).Children(func() {
				ui.Text(c, body).Font("monospace").FontSize(t.Rem(0.8)).Wrap()
			})
		})
	})
	if head.Clicked() {
		*open = !*open
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("ToolResult", "open", *open)
	}
	if copy.Clicked() {
		kit.Log("ToolResult", "copied", "lines", lines)
	}
	return root, copy
}
