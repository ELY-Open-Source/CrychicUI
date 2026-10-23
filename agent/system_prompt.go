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
	Placeholder string
	Tokens      string
}

// SystemPrompt is an editor with a token count.
func SystemPrompt(c *ui.Context, label string, text *string, o SystemPromptOptions) ui.Element {
	if label == "" || text == nil {
		panic(fmt.Sprintf("agent: SystemPrompt %q needs a label and text", label))
	}
	def(&o.Placeholder, "You are a helpful assistant.")
	def(&o.Tokens, "≈ %d tokens")
	t := c.Theme()
	root := ui.Column(c).FillWidth().Gap(t.Space(1)).Role(ui.RoleGroup).Label(label)
	return root.Children(func() {
		ui.Text(c, label).Bold().SingleLine()
		in := ui.TextArea(c, text).FillWidth().Lines(4, 10).Placeholder(o.Placeholder).Label(label).
			Padding(t.Space(2), t.Space(3)).Radius(t.Radius).Background(t.Background).Border(1, t.Border)
		if in.Changed() {
			root.SetChanged()
			sound.Play(sound.Tick)
			kit.Log("SystemPrompt", "edited", "runes", utf8.RuneCountInString(*text))
		}
		tokens := (utf8.RuneCountInString(*text) + 3) / 4
		kit.Numeric(c, fmt.Sprintf(o.Tokens, tokens)).FontSize(t.Rem(0.8)).TextColor(t.TextMuted).AlignSelf(ui.End)
	})
}
