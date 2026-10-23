package chat

import (
	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/egoist/mygo/ui"
)

// ForwardedOptions word the header, "%s" taking the source.
type ForwardedOptions struct {
	From string
}

// Forwarded is a "forwarded from" header in a bubble.
func Forwarded(c *ui.Context, from string, m Message, o ForwardedOptions) ui.Element {
	if from == "" {
		panic("chat: Forwarded needs a source")
	}
	messages("Forwarded", []Message{m})
	def(&o.From, "Forwarded from %s")
	t := c.Theme()
	root := shell(c, m.Mine)
	kit.Arrive(root, ui.Pose{Y: t.Space(2), Fade: 1}, 0)
	return root.Children(func() {
		ui.Row(c).Gap(t.Space(1)).AlignItems(ui.Center).Opacity(0.75).Children(func() {
			ui.Icon(c, icon.Must("forward")).Size(16, 16)
			ui.Text(c, sprintf(o.From, from)).Italic().FontSize(t.Rem(0.75)).SingleLine()
		})
		ui.Text(c, m.Text).Wrap()
	})
}
