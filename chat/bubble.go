package chat

import (
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// BubbleOptions tune a Bubble.
type BubbleOptions struct {
	// ShowTime shows the message time under the text.
	ShowTime bool
}

// Bubble is one message, mine on the right in accent.
func Bubble(c *ui.Context, m Message, o BubbleOptions) ui.Element {
	messages("Bubble", []Message{m})
	t := c.Theme()
	root := shell(c, m.Mine)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	return root.Children(func() {
		ui.Text(c, m.Text).Wrap()
		if o.ShowTime {
			ui.Text(c, when(c, m.At)).FontSize(t.Rem(0.75)).AlignSelf(ui.End).Opacity(0.7).
				ContentTransition(ui.ContentTransition{})
		}
	})
}
