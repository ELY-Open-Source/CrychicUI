package agent

import (
	"fmt"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// SystemPromptOptions tune a SystemPrompt editor.
type SystemPromptOptions struct {
	Placeholder, Tokens, Hint string
}

// SystemPrompt is an editor with a token count.
func SystemPrompt(c *ui.Context, label string, text *string, o SystemPromptOptions) ui.Element {
	if label == "" || text == nil {
		panic(fmt.Sprintf("agent: SystemPrompt %q needs a label and text", label))
	}
	def(&o.Placeholder, "You are a helpful assistant.")
	def(&o.Tokens, "≈ %d tokens")
	def(&o.Hint, "Sent before every message")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Radius(t.Radius).Background(t.Background).Clip().Role(ui.RoleGroup).Label(label)
	var in ui.Element
	root.Children(func() {
		ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.Center).Padding(t.Space(2), t.Space(3)).Background(t.Surface).
			BorderWidth(0, 0, 1, 0).BorderColor(t.Border).Children(func() {
			tile(c, "square-terminal", t.Accent, t.Space(7))
			ui.Column(c).Grow(1).MinWidth(0).Children(func() {
				ui.Text(c, label).Bold().SingleLine()
				kit.Muted(c, o.Hint).FontSize(t.Rem(0.75)).SingleLine()
			})
			tokens := (utf8.RuneCountInString(*text) + 3) / 4
			kit.Numeric(c, fmt.Sprintf(o.Tokens, tokens)).FontSize(t.Rem(0.75)).Bold().TextColor(t.Accent).
				Padding(t.Space(0.5), t.Space(2)).Radius(999).Background(t.Accent.Alpha(0.12)).ContentTransition(ui.ContentTransition{})
		})
		in = ui.TextAreaBase(c, text).FillWidth().Lines(4, 12).Placeholder(o.Placeholder).Label(label).
			Font("monospace").FontSize(t.Rem(0.85)).Padding(t.Space(3))
	})
	root.Border(1, root.AnimateColor("ring", t.Border.Mix(t.Accent, unit(root.FocusWithin())), ui.Snappy))
	if in.Changed() {
		root.SetChanged()
		sound.Play(sound.Tick)
		kit.Log("SystemPrompt", "edited", "runes", utf8.RuneCountInString(*text))
	}
	return root
}
