package chat

import (
	"strings"
	"unicode/utf8"

	"github.com/ELY-Open-Source/CrychicUI/icon"
	"github.com/ELY-Open-Source/CrychicUI/kit"
	"github.com/ELY-Open-Source/CrychicUI/sound"
	"github.com/egoist/mygo/ui"
)

// ComposerOptions tune a Composer.
type ComposerOptions struct {
	Placeholder, Attach, Send string
}

// Composer has attach, a text area and send; Enter sends.
func Composer(c *ui.Context, draft *string, onSend func(text string), o ComposerOptions) (root, send, attach ui.Element) {
	if draft == nil || onSend == nil {
		panic("chat: Composer needs a draft and a send func")
	}
	def(&o.Placeholder, "Message")
	def(&o.Attach, "Attach")
	def(&o.Send, "Send")
	t := c.Theme()
	root = ui.Row(c).FillWidth().Gap(t.Space(2)).AlignItems(ui.End).Padding(t.Space(2)).Radius(t.Radius).
		Background(t.Background).Border(1, t.Border)
	pending := ui.Local(root, "pending", func() string { return "" })
	blank := strings.TrimSpace(*draft) == ""
	root.Children(func() {
		attach = round(c, o.Attach, "paperclip")
		in := ui.TextArea(c, draft).Grow(1).Lines(1, 6).Placeholder(o.Placeholder).Label(o.Placeholder).OnTextKey(func(k ui.TextKey) bool {
			if k.Key != ui.KeyEnter || k.Mods != 0 {
				return false
			}
			if text := strings.TrimSpace(k.Text()); text != "" {
				*pending = text
				k.SetText("")
			}
			return true
		})
		if in.Changed() {
			root.SetChanged()
		}
		send = kit.Press(ui.ButtonBase(c).Label(o.Send).Tooltip(o.Send).Size(t.Space(8), t.Space(8)).Shrink(0).Center().Radius(999)).
			Disabled(blank)
		send.Background(send.AnimateColor("bg", t.Accent.Mix(t.AccentHover, unit(send.Hovered())), ui.Snappy)).
			TextColor(t.AccentText).Children(func() { ui.Icon(c, icon.Must("send")).Size(16, 16) })
	})
	if kit.Clicked(attach, sound.Tap) {
		kit.Log("Composer", "attach")
	}
	text := *pending
	*pending = ""
	if text == "" && !blank && send.Clicked() {
		text = strings.TrimSpace(*draft)
		*draft = ""
	}
	if text != "" {
		onSend(text)
		root.SetChanged()
		sound.Play(sound.Send)
		kit.Log("Composer", "sent", "runes", utf8.RuneCountInString(text))
	}
	return root, send, attach
}
