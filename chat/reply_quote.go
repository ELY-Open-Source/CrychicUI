package chat

import (
	"fmt"

	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ReplyQuoteOptions tune a ReplyQuote.
type ReplyQuoteOptions struct {
	// Label names the quote for assistive technology.
	Label string
}

// ReplyQuote quotes a message; a click jumps to it.
func ReplyQuote(c *ui.Context, sender, text string, o ReplyQuoteOptions) ui.Element {
	if sender == "" || text == "" {
		panic(fmt.Sprintf("chat: ReplyQuote %q needs a sender and text", sender))
	}
	def(&o.Label, "Quoted message")
	t := c.Theme()
	root := kit.Press(ui.ButtonBase(c).Label(o.Label).Row().FillWidth().Gap(t.Space(2)).Padding(t.Space(1.5), t.Space(2)).Radius(t.Radius))
	root.Background(root.AnimateColor("bg", t.Text.Alpha(0.05+0.05*unit(root.Hovered())), ui.Snappy))
	root.Children(func() {
		ui.Box(c).Width(t.Space(0.75)).Shrink(0).FillHeight().Radius(999).Background(t.Accent)
		ui.Column(c).Grow(1).Children(func() {
			ui.Text(c, sender).Bold().FontSize(t.Rem(0.85)).TextColor(t.Accent).SingleLine()
			kit.Muted(c, text).SingleLine().FontSize(t.Rem(0.9))
		})
	})
	if kit.Clicked(root, sound.Tap) {
		kit.Log("ReplyQuote", "jumped")
	}
	return root
}
