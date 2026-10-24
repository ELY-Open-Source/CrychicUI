package agent

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/layout"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ToolCallOptions word the states and the input heading.
type ToolCallOptions struct {
	Pending, Running, Done, Failed, Input string
}

// ToolCall shows a tool name, arguments and status; it expands.
func ToolCall(c *ui.Context, name, args string, state StepState, open *bool, o ToolCallOptions) ui.Element {
	if name == "" || open == nil || state > StepFailed {
		panic(fmt.Sprintf("agent: ToolCall %q needs a name, open and a known state, got %d", name, state))
	}
	def(&o.Pending, "Pending")
	def(&o.Running, "Running")
	def(&o.Done, "Done")
	def(&o.Failed, "Failed")
	def(&o.Input, "Input")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Radius(t.Radius).Background(t.Background).Border(1, t.Border).Clip()
	var head ui.Element
	root.Children(func() {
		head = ui.ButtonBase(c).Label(name).Row().FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).
			Padding(t.Space(2), t.Space(3)).Checked(*open)
		head.Background(head.AnimateColor("bg", t.SurfaceHover.Alpha(unit(head.Hovered())), ui.Snappy))
		head.Children(func() {
			mark(c, state)
			ui.Icon(c, icon.Must("wrench")).Size(16, 16).TextColor(t.TextMuted)
			ui.Text(c, name).Bold().Font("monospace").FontSize(t.Rem(0.85)).Shrink(0).SingleLine()
			kit.Muted(c, args).Font("monospace").FontSize(t.Rem(0.8)).Grow(1).SingleLine()
			stateBadge(c, state, o.Pending, o.Running, o.Done, o.Failed)
			ch := ui.Icon(c, icon.Must("chevron-down")).Size(16, 16).TextColor(t.TextMuted)
			ch.Rotate(180 * ch.AnimateTo("turn", unit(*open), ui.Snappy))
		})
		layout.Collapse(c, *open, func() {
			ui.Column(c).FillWidth().Gap(t.Space(1.5)).Padding(t.Space(2), t.Space(3), t.Space(3)).Children(func() {
				kit.Muted(c, o.Input).FontSize(t.Rem(0.7)).Bold()
				text := args
				if text == "" {
					text = "{}"
				}
				ui.Text(c, text).Font("monospace").FontSize(t.Rem(0.8)).FillWidth().Wrap().Padding(t.Space(2)).
					Radius(t.Radius).Background(t.Surface)
			})
		})
	})
	if head.Clicked() {
		*open = !*open
		root.SetChanged()
		sound.Play(sound.Toggle)
		kit.Log("ToolCall", "open", *open)
	}
	return root
}
