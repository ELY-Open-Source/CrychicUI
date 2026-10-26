package desktop

import (
	"strings"
	"time"

	"github.com/egoist/mygo/ui"
)

// TerminalWindowOptions set the prompt that marks typed lines.
type TerminalWindowOptions struct{ Prompt string }

// TerminalWindow is a terminal window with a blinking caret.
func TerminalWindow(c *ui.Context, title string, lines []string, o TerminalWindowOptions) ui.Element {
	def(&o.Prompt, "$")
	if strings.TrimSpace(title) == "" {
		panic("desktop: TerminalWindow needs a title")
	}
	t := c.Theme()
	bg, fg := t.Text, t.Background
	root := window(c).Role(ui.RoleGroup).Label(title)
	return root.Children(func() {
		titleBar(c, title, true, nil)
		view := ui.Scroll(c).FillWidth().MaxHeight(t.Space(64)).Background(bg)
		view.Children(func() {
			ui.Column(c).FillWidth().Gap(t.Space(0.5)).Padding(t.Space(3), t.Space(4)).Children(func() {
				for i, line := range lines {
					c.Key(i)
					if cmd, ok := strings.CutPrefix(line, o.Prompt+" "); ok {
						ui.RichText(c, ui.Span{Text: o.Prompt + " ", Color: t.Success, Weight: 700}, ui.Span{Text: cmd, Color: fg, Weight: 600}).
							Font("monospace").FontSize(t.Rem(0.78)).LineHeight(1.5).FillWidth().Wrap()
						continue
					}
					ui.Text(c, line).Font("monospace").FontSize(t.Rem(0.78)).LineHeight(1.5).TextColor(fg.Alpha(0.72)).FillWidth().Wrap()
				}
				c.Key("caret")
				ui.Row(c).Gap(t.Space(1.5)).AlignItems(ui.Center).Children(func() {
					ui.Text(c, o.Prompt).Font("monospace").FontSize(t.Rem(0.78)).Bold().TextColor(t.Success)
					caret := ui.Box(c).Size(t.Space(2), t.Space(4))
					blink := caret.Loop("blink", 1100*time.Millisecond, ui.Linear)
					caret.Background(fg.Alpha(unit(blink < 0.5)))
				})
			})
		})
	})
}
